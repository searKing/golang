// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leaderelection

import (
	"context"
	"time"
)

type Config struct {
	// Lock is the resource that will be used for locking
	Lock ResourceLocker

	// LeaseDuration is the duration that non-leader candidates will
	// wait to force acquire leadership. This is measured against time of
	// last observed ack.
	//
	// A client needs to wait a full LeaseDuration without observing a change to
	// the record before it can attempt to take over. When all clients are
	// shutdown and a new set of clients are started with different names against
	// the same leader record, they must wait the full LeaseDuration before
	// attempting to acquire the lease. Thus LeaseDuration should be as short as
	// possible (within your tolerance for clock skew rate) to avoid a possible
	// long waits in the scenario.
	//
	// Core clients default this value to 15 seconds.
	LeaseDuration time.Duration
	// RenewTimeout is the duration that the acting master will retry
	// refreshing leadership before giving up.
	//
	// Core clients default this value to 10 seconds.
	RenewTimeout time.Duration
	// RetryPeriod is the duration the LeaderElector clients should wait
	// between tries of actions.
	//
	// Core clients default this value to 2 seconds.
	RetryPeriod time.Duration

	// Callbacks are callbacks that are triggered during certain lifecycle
	// events of the LeaderElector
	Callbacks LeaderCallbacks

	// ReleaseOnCancel should be set true if the lock should be released
	// when the run context is cancelled. If you set this to true, you must
	// ensure all code guarded by this lease has successfully completed
	// prior to cancelling the context, or you may have two processes
	// simultaneously acting on the critical path.
	ReleaseOnCancel bool

	// Name is the name of the resource lock for debugging
	Name string
}

// LeaderCallbacks are callbacks that are triggered during certain
// lifecycle events of the LeaderElector. These are invoked asynchronously.
//
// possible future callbacks:
//   - OnChallenge()
type LeaderCallbacks struct {
	// OnStartedLeading is called when a LeaderElector client starts leading
	OnStartedLeading func(context.Context)
	// OnStoppedLeading is called when a LeaderElector client stops leading
	OnStoppedLeading func()
	// OnNewLeader is called when the client observes a leader that is
	// not the previously observed leader. This includes the first observed
	// leader when the client starts.
	OnNewLeader func(identity string)
}

func (c *Config) SetDefaults() {
	if c != nil {
		c.Lock = NewDummyLock("")
		c.LeaseDuration = 15 * time.Second
		c.RenewTimeout = 10 * time.Second
		c.RetryPeriod = 2 * time.Second
	}
}

func (c *Config) Complete() {
	if c.LeaseDuration < 1 {
		c.LeaseDuration = 15 * time.Second
	}
	if c.RenewTimeout < 1 {
		c.RenewTimeout = 10 * time.Second
	}
	if c.RetryPeriod < 1 {
		c.RetryPeriod = 2 * time.Second
	}
	if c.LeaseDuration <= c.RenewTimeout {
		c.LeaseDuration = c.RenewTimeout
	}
	if c.RenewTimeout <= time.Duration(JitterFactor*float64(c.RetryPeriod)) {
		c.RenewTimeout = time.Duration(JitterFactor * float64(c.RetryPeriod))
	}
}
func (c *Config) New() (*LeaderElector, error) {
	return NewLeaderElector(*c)
}
