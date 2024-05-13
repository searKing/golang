// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"sync"
	"time"
)

// PersistResource wraps a resource, usually a persistent one
// (but may be used for non-keep-alive requests as well)
type PersistResource struct {
	t        *LruPool
	cacheKey targetKey
	object   any

	// Both guarded by LruPool.idleMu:
	idleAt    time.Time   // time it last become idle
	idleTimer *time.Timer // holding an AfterFunc to close it

	mu     sync.Mutex // guards following fields
	closed error      // set non-nil when resource is closed, before closech is closed
	broken bool       // an error has happened on this resource; marked broken so it's not reused.
	reused bool       // whether resource has had successful request/response and is being reused.
} // isBroken reports whether this resource is in a known broken state.

func (pc *PersistResource) Get() any {
	if pc == nil {
		return nil
	}
	return pc.object
}

func (pc *PersistResource) Put() {
	if pc == nil {
		return
	}
	pc.t.putOrCloseIdleResource(pc)
}

func (pc *PersistResource) isBroken() bool {
	pc.mu.Lock()
	b := pc.closed != nil
	pc.mu.Unlock()
	return b
}

// markReused marks this resource as having been successfully used for a
// request and response.
func (pc *PersistResource) markReused() {
	pc.mu.Lock()
	pc.reused = true
	pc.mu.Unlock()
}

// close closes the underlying resource and closes
// the pc.closech channel.
//
// The provided err is only for testing and debugging; in normal
// circumstances it should never be seen by users.
func (pc *PersistResource) close(err error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.closeLocked(err)
}

func (pc *PersistResource) closeLocked(err error) {
	if err == nil {
		panic("nil error")
	}
	pc.broken = true
	if pc.closed == nil {
		pc.closed = err
		pc.t.decResourcesPerBucket(pc.cacheKey)
	}
}

// closeResourceIfStillIdle closes the resource if it's still sitting idle.
// This is what's called by the PersistResource's idleTimer, and is run in its
// own goroutine.
func (pc *PersistResource) closeResourceIfStillIdle() {
	t := pc.t
	t.idleMu.Lock()
	defer t.idleMu.Unlock()
	if _, ok := t.idleLRU.m[pc]; !ok {
		// Not idle.
		return
	}
	t.removeIdleResourceLocked(pc)
	pc.close(errIdleResourceTimeout)
}
