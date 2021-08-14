// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leaderelection

import (
	"context"
	"time"
)

// Record is the record that is stored in the leader election annotation.
// This information should be used for observational purposes only and could be replaced
// with a random string (e.g. UUID) with only slight modification of this code.
type Record struct {
	// HolderIdentity is the ID that owns the lease. If empty, no one owns this lease and
	// all callers may acquire.
	// This value is set to empty when a client voluntarily steps down.
	HolderIdentity string
	// LeaseDuration is the duration that non-leader candidates will
	// wait to force acquire leadership. This is measured against time of
	// last observed ack.
	LeaseDuration     time.Duration
	AcquireTime       time.Time // when the leader hold this record
	RenewTime         time.Time // when the locker is renewed recently
	LeaderTransitions int       // +1 if expired checked by followers, changed to trigger a new lock acquire or new
}

// ResourceLocker offers a common interface for locking on arbitrary
// resources used in leader election.  The ResourceLocker is used
// to hide the details on specific implementations in order to allow
// them to change over time.  This interface is strictly for use
// by the leaderelection code.
type ResourceLocker interface {
	// Get returns the LeaderElectionRecord
	// get locker's state
	Get(ctx context.Context) (record *Record, rawRecord []byte, err error)

	// Create attempts to create a LeaderElectionRecord
	// return err if lock failed
	// Lock
	Create(ctx context.Context, ler Record) error

	// Update will update and existing LeaderElectionRecord
	// return err if lock failed
	// Lock or UnLock
	Update(ctx context.Context, ler Record) error

	// RecordEvent is used to record events
	RecordEvent(name, event string)

	// Identity will return the locks Identity
	Identity() string

	// Describe is used to convert details on current resource lock
	// into a string
	Describe() string
}
