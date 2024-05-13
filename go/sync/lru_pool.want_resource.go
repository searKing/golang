// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"sync"
)

// A wantResource records state about a wanted resource.
// The resource may be gotten by New() or by finding an idle resource,
// or a cancellation may make the resource no longer wanted.
// These three options are racing against each other and use
// wantResource to coordinate and agree about the winning outcome.
type wantResource struct {
	req   any
	key   any             // cm.key()
	ctx   context.Context // context for New
	ready chan struct{}   // closed when pr, err pair is delivered

	mu  sync.Mutex // protects pr, err, close(ready)
	pr  *PersistResource
	err error
}

// waiting reports whether w is still waiting for an answer (connection or error).
func (w *wantResource) waiting() bool {
	select {
	case <-w.ready:
		return false
	default:
		return true
	}
}

// deliveredLock returns whether resource has been delivered already.
func (w *wantResource) deliveredLock() bool {
	return w.pr != nil || w.err != nil
}

// tryDeliver attempts to deliver pr, err to w and reports whether it succeeded.
func (w *wantResource) tryDeliver(pr *PersistResource, err error) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.deliveredLock() {
		return false
	}

	// deliver and notify delivered event
	w.pr = pr
	w.err = err
	if w.pr == nil && w.err == nil {
		panic("sync: internal error: misuse of tryDeliver")
	}
	close(w.ready)
	return true
}

// cancel marks w as no longer wanting a result (for example, due to cancellation).
// If a connection has been delivered already, cancel returns it with t.putOrCloseIdleResource.
func (w *wantResource) cancel(t *LruPool, err error) {
	w.mu.Lock()
	if !w.deliveredLock() {
		// notify deliver cancelled event
		close(w.ready) // catch misbehavior in future delivery
	}

	// deliver and notify delivered event
	pc := w.pr
	w.pr = nil
	w.err = err // not nil always
	w.mu.Unlock()

	if pc != nil {
		t.putOrCloseIdleResource(pc)
	}
}
