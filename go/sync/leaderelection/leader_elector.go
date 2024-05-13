// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package leaderelection implements leader election of a set of endpoints.
// It uses an annotation in the endpoints object to store the record of the
// election state. This implementation does not guarantee that only one
// client is acting as a leader (a.k.a. fencing).
//
// A client only acts on timestamps captured locally to infer the state of the
// leader election. The client does not consider timestamps in the leader
// election record to be accurate because these timestamps may not have been
// produced by a local clock. The implementation does not depend on their
// accuracy and only uses their change to indicate that another client has
// renewed the leader lease. Thus, the implementation is tolerant to arbitrary
// clock skew, but is not tolerant to arbitrary clock skew rate.
//
// However, the level of tolerance to skew rate can be configured by setting
// RenewDeadline and LeaseDuration appropriately. The tolerance expressed as a
// maximum tolerated ratio of time passed on the fastest node to time passed on
// the slowest node can be approximately achieved with a configuration that sets
// the same ratio of LeaseDuration to RenewDeadline. For example if a user wanted
// to tolerate some nodes progressing forward in time twice as fast as other nodes,
// the user could set LeaseDuration to 60 seconds and RenewDeadline to 30 seconds.
//
// While not required, some method of clock synchronization between nodes in the
// cluster is highly recommended. It's important to keep in mind when configuring
// this client that the tolerance to skew rate varies inversely to master
// availability.
//
// Larger clusters often have a more lenient SLA for API latency. This should be
// taken into account when configuring the client. The rate of leader transitions
// should be monitored and RetryPeriod and LeaseDuration should be increased
// until the rate is stable and acceptably low. It's important to keep in mind
// when configuring this client that the tolerance to API latency varies inversely
// to master availability.
package leaderelection

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/searKing/golang/go/runtime"
	time_ "github.com/searKing/golang/go/time"
)

const (
	JitterFactor        = 1.2
	EventBecameLeader   = "became leader"
	EventStoppedLeading = "stopped leading"
)

// NewLeaderElector creates a LeaderElector from a LeaderElectionConfig
func NewLeaderElector(lec Config) (*LeaderElector, error) {
	if lec.LeaseDuration <= lec.RenewTimeout {
		return nil, fmt.Errorf("leaseDuration must be greater than renewDeadline")
	}
	if lec.RenewTimeout <= time.Duration(JitterFactor*float64(lec.RetryPeriod)) {
		lec.RenewTimeout = time.Duration(JitterFactor * float64(lec.RetryPeriod))
		return nil, fmt.Errorf("renewDeadline must be greater than retryPeriod*JitterFactor")
	}
	if lec.LeaseDuration < 1 {
		return nil, fmt.Errorf("leaseDuration must be greater than zero")
	}
	if lec.RenewTimeout < 1 {
		return nil, fmt.Errorf("renewDeadline must be greater than zero")
	}
	if lec.RetryPeriod < 1 {
		return nil, fmt.Errorf("retryPeriod must be greater than zero")
	}

	if lec.Lock == nil {
		return nil, fmt.Errorf("lock must not be nil")
	}
	le := LeaderElector{
		config: lec,
	}
	le.config.Lock.RecordEvent(le.config.Name, EventStoppedLeading)

	return &le, nil
}

// LeaderElector is a leader election client.
type LeaderElector struct {
	// ErrorLog specifies an optional logger for errors accepting
	// connections, unexpected behavior from handlers, and
	// underlying FileSystem errors.
	// If nil, logging is done via the log package's standard logger.
	ErrorLog *log.Logger

	config Config
	// internal bookkeeping
	observedRecord    Record
	observedRawRecord []byte
	observedTime      time.Time // Time when setObservedRecord is called
	// used to implement OnNewLeader(), may lag slightly from the
	// value observedRecord.HolderIdentity if the transition has
	// not yet been reported.
	reportedLeader string

	// used to lock the observedRecord
	observedRecordLock sync.Mutex
}

func (le *LeaderElector) logf(format string, args ...any) {
	if le.ErrorLog != nil {
		le.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// Run starts the leader election loop. Run will not return
// before leader election loop is stopped by ctx, or it has
// stopped holding the leader lease
func (le *LeaderElector) Run(ctx context.Context) {
	defer runtime.DefaultPanic.Recover()
	if le.config.Callbacks.OnStoppedLeading != nil {
		defer func() {
			le.config.Callbacks.OnStoppedLeading()
		}()
	}

	// wait until we are a leader
	if !le.acquire(ctx) {
		return // ctx signalled done
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if le.config.Callbacks.OnStartedLeading != nil {
		go le.config.Callbacks.OnStartedLeading(ctx)
	}
	le.renew(ctx)
}

// Run starts a client with the provided config or failed if the config
// fails to validate. RunOrDie blocks until leader election loop is
// stopped by ctx or it has stopped holding the leader lease
func Run(ctx context.Context, lec Config) (*LeaderElector, error) {
	le, err := NewLeaderElector(lec)
	if err != nil {
		return nil, err
	}
	le.Run(ctx)
	return nil, err
}

// GetLeader returns the identity of the last observed leader or returns the empty string if
// no leader has yet been observed.
// This function is for informational purposes. (e.g. monitoring, logs, etc.)
func (le *LeaderElector) GetLeader() string {
	return le.getObservedRecord().HolderIdentity
}

// IsLeader returns true if the last observed leader was this client else returns false.
func (le *LeaderElector) IsLeader() bool {
	return le.getObservedRecord().HolderIdentity == le.config.Lock.Identity()
}

// acquire loops calling tryAcquireOrRenew and returns true immediately when tryAcquireOrRenew succeeds.
// Returns false if ctx signals done.
func (le *LeaderElector) acquire(ctx context.Context) bool {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	succeeded := false
	desc := le.config.Lock.Describe()
	le.logf("attempting to acquire leader lease %v...", desc)
	time_.JitterUntil(ctx, func(ctx context.Context) {
		// try lock
		succeeded = le.tryAcquireOrRenew(ctx)
		le.maybeReportTransition()
		if !succeeded {
			le.logf("failed to acquire lease %v", desc)
			return
		}
		le.config.Lock.RecordEvent(le.config.Name, EventBecameLeader)
		le.logf("successfully acquired lease %v", desc)
		cancel()
	}, true, time_.WithExponentialBackOffOptionRandomizationFactor(JitterFactor))
	return succeeded
}

// renew loops calling tryAcquireOrRenew and returns immediately when tryAcquireOrRenew fails or ctx signals done.
func (le *LeaderElector) renew(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	time_.Until(ctx, func(ctx context.Context) {
		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, le.config.RenewTimeout)
		defer timeoutCancel()
		var leader bool
		// block until leader elector finished
		time_.Until(timeoutCtx, func(ctx context.Context) {
			if le.tryAcquireOrRenew(timeoutCtx) {
				leader = true
				timeoutCancel()
			}
		}, le.config.RenewTimeout)
		// leader elector finished

		// maybe report leader changed
		le.maybeReportTransition()
		desc := le.config.Lock.Describe()
		if leader {
			// I'm a leader now
			le.logf("successfully renewed lease %v", desc)
			return
		}
		// I'm a follower now
		le.config.Lock.RecordEvent(le.config.Name, EventStoppedLeading)
		le.logf("failed to renew lease %v: %v", desc, timeoutCtx.Err())
		cancel()

	}, le.config.RetryPeriod)

	// if we hold the lease, give it up
	// unlock if I'm a leader
	if le.config.ReleaseOnCancel {
		le.release()
	}
}

// release attempts to release the leader lease if we have acquired it.
// unlock is we hold this lock
func (le *LeaderElector) release() bool {
	if !le.IsLeader() {
		return true
	}
	now := time.Now()
	leaderElectionRecord := Record{
		LeaderTransitions: le.observedRecord.LeaderTransitions,
		LeaseDuration:     time.Second,
		RenewTime:         now,
		AcquireTime:       now,
	}
	if err := le.config.Lock.Update(context.TODO(), leaderElectionRecord); err != nil {
		le.logf("Failed to release lock: %v", err)
		return false
	}

	le.setObservedRecord(&leaderElectionRecord)
	return true
}

// tryAcquireOrRenew tries to acquire a leader lease if it is not already acquired,
// else it tries to renew the lease if it has already been acquired. Returns true
// on success else returns false.
func (le *LeaderElector) tryAcquireOrRenew(ctx context.Context) bool {
	now := time.Now()
	leaderElectionRecord := Record{
		HolderIdentity: le.config.Lock.Identity(),
		LeaseDuration:  le.config.LeaseDuration,
		RenewTime:      now,
		AcquireTime:    now,
	}

	// 1. obtain or create the ElectionRecord
	oldLeaderElectionRecord, oldLeaderElectionRawRecord, err := le.config.Lock.Get(ctx)
	if err != nil {
		// Lock, try to lock as I'm a leader
		if err = le.config.Lock.Create(ctx, leaderElectionRecord); err != nil {
			le.logf("error initially creating leader election record: %v", err)
			return false
		}

		le.setObservedRecord(&leaderElectionRecord)
		return true
	}
	// renew

	// 2. Record obtained, check the Identity & Time
	if !bytes.Equal(le.observedRawRecord, oldLeaderElectionRawRecord) {
		le.setObservedRecord(oldLeaderElectionRecord)
		le.observedRawRecord = oldLeaderElectionRawRecord
	}

	if len(oldLeaderElectionRecord.HolderIdentity) > 0 &&
		!le.observedRecordExpired(now) && !le.IsLeader() {
		le.logf("lock is held by %v and has not yet expired", oldLeaderElectionRecord.HolderIdentity)
		// return as a follower
		return false
	}

	// 3. We're going to try to update. The leaderElectionRecord is set to it's default
	// here. Let's correct it before updating.
	if le.IsLeader() {
		// refresh the locker by leader self
		// relock, the lock is inherited, so AcquireTime is kept
		// refresh the lock by RenewTime refreshed
		leaderElectionRecord.AcquireTime = oldLeaderElectionRecord.AcquireTime
		leaderElectionRecord.LeaderTransitions = oldLeaderElectionRecord.LeaderTransitions
	} else {
		// try to lock as a leader
		leaderElectionRecord.LeaderTransitions = oldLeaderElectionRecord.LeaderTransitions + 1
	}

	// update the lock as a leader
	if err = le.config.Lock.Update(ctx, leaderElectionRecord); err != nil {
		le.logf("Failed to update lock: %v", err)
		return false
	}

	le.setObservedRecord(&leaderElectionRecord)
	return true
}

// maybeReportTransition call OnNewLeader when the client observes a leader that is
//
//	not the previously observed leader.
func (le *LeaderElector) maybeReportTransition() {
	if le.observedRecord.HolderIdentity == le.reportedLeader {
		return
	}
	le.reportedLeader = le.observedRecord.HolderIdentity
	if le.config.Callbacks.OnNewLeader != nil {
		go le.config.Callbacks.OnNewLeader(le.reportedLeader)
	}
}

// Check will determine if the current lease is expired by more than timeout.
func (le *LeaderElector) Check(maxTolerableExpiredLease time.Duration) error {
	if !le.IsLeader() {
		// Currently, not concerned with the case that we are hot standby
		return nil
	}
	// If we are more than timeout seconds after the lease duration that is past the timeout
	// on the lease renew. Time to start reporting ourselves as unhealthy. We should have
	// died but conditions like deadlock can prevent this. (See #70819)
	if time.Since(le.observedTime) > le.config.LeaseDuration+maxTolerableExpiredLease {
		return fmt.Errorf("failed election to renew leadership on lease %s", le.config.Name)
	}

	return nil
}

// setObservedRecord will set a new observedRecord and update observedTime to the current time.
// Protect critical sections with lock.
func (le *LeaderElector) setObservedRecord(observedRecord *Record) {
	le.observedRecordLock.Lock()
	defer le.observedRecordLock.Unlock()

	le.observedRecord = *observedRecord
	le.observedTime = time.Now()
}

// getObservedRecord returns observersRecord.
// Protect critical sections with lock.
func (le *LeaderElector) getObservedRecord() Record {
	le.observedRecordLock.Lock()
	defer le.observedRecordLock.Unlock()

	return le.observedRecord
}

// observedRecordExpired returns true if observersRecord expired.
// Protect critical sections with lock.
func (le *LeaderElector) observedRecordExpired(now time.Time) bool {
	return le.observedTime.Add(le.config.LeaseDuration).After(now)
}
