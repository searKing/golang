// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

var errRequestCanceledResource = errors.New("sync: request canceled while waiting for resource")

// errRequestCanceled is a copy of sync's errRequestCanceled because it's not
// exported. At least they'll be DeepEqual for h1-vs-h2 comparisons tests.
var errRequestCanceled = errors.New("sync: request canceled")

var (
	errKeepAlivesDisabled  = errors.New("sync: putIdleResource: keep alives disabled")
	errResourceBroken      = errors.New("sync: putIdleResource: resource is in bad state")
	errCloseIdle           = errors.New("sync: putIdleResource: CloseIdleResources was called")
	errTooManyIdle         = errors.New("sync: putIdleResource: too many idle resources")
	errTooManyIdleResource = errors.New("sync: putIdleResource: too many idle resources for host")

	errIdleResourceTimeout = errors.New("sync: idle resource timeout")
)

// DefaultMaxIdleResourcesPerHost is the default value of LruPool's
// MaxIdleResourcesPerHost.
const DefaultMaxIdleResourcesPerHost = 2

// A cancelKey is the key of the reqCanceler map.
// We wrap the *Request in this type since we want to use the original request,
// not any transient one created by roundTrip.
type cancelKey struct {
	req interface{}
}

type targetKey interface{}

type LruPool struct {
	// New optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	New func(ctx context.Context) interface{}

	idleMu           sync.Mutex
	closeIdle        bool                             // user has requested to close all idle conns
	idleResource     map[targetKey][]*PersistResource // most recently used at end
	idleResourceWait map[targetKey]wantResourceQueue  // waiting getResources
	idleLRU          connLRU
	reqMu            sync.Mutex
	reqCanceler      map[cancelKey]func(error)

	connsPerHostMu   sync.Mutex
	connsPerHost     map[targetKey]int
	connsPerHostWait map[targetKey]wantResourceQueue // waiting getResources
	// DisableKeepAlives, if true, disables HTTP keep-alives and
	// will only use the resource to the server for a single
	// HTTP request.
	//
	// This is unrelated to the similarly named TCP keep-alives.
	DisableKeepAlives bool

	// MaxIdleResources controls the maximum number of idle (keep-alive)
	// resources across all hosts. Zero means no limit.
	MaxIdleResources int

	// MaxIdleResourcesPerHost, if non-zero, controls the maximum idle
	// (keep-alive) resources to keep per-host. If zero,
	// DefaultMaxIdleResourcesPerHost is used.
	MaxIdleResourcesPerHost int

	// MaxResourcesPerHost optionally limits the total number of
	// resources per host, including resources in the dialing,
	// active, and idle states. On limit violation, dials will block.
	//
	// Zero means no limit.
	MaxResourcesPerHost int

	// IdleResourceTimeout is the maximum amount of time an idle
	// (keep-alive) resource will remain idle before closing
	// itself.
	// Zero means no limit.
	IdleResourceTimeout time.Duration
}

// GetOrError creates a new PersistResource to the target as specified in the key.
// If this doesn't return an error, the PersistResource is ready to write requests to.
func (t *LruPool) GetOrError(ctx context.Context, key interface{}) (pc *PersistResource, err error) {

	w := &wantResource{
		key:   key,
		ctx:   ctx,
		ready: make(chan struct{}, 1),
	}
	defer func() {
		if err != nil {
			w.cancel(t, err)
		}
	}()
	cancelKey := cancelKey{req: key}

	// Queue for idle resource.
	if delivered := t.queueForIdleResource(w); delivered {
		pc := w.pr
		// set request canceler to some non-nil function, so we
		// can detect whether it was cleared between now and when
		// we enter resolving
		t.setReqCanceler(cancelKey, func(error) {})
		return pc, nil
	}

	cancelc := make(chan error, 1)
	t.setReqCanceler(cancelKey, func(err error) { cancelc <- err })

	// Queue for permission to dial.
	t.queueForDial(w)

	// Wait for completion or cancellation.
	select {
	case <-w.ready:
		if w.err != nil {
			// If the request has been cancelled, that's probably
			// what caused w.err; if so, prefer to return the
			// cancellation error (see golang.org/issue/16049).
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case err := <-cancelc:
				if err == errRequestCanceled {
					err = errRequestCanceledResource
				}
				return nil, err
			default:
				// return below
			}
		}
		return w.pr, w.err
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-cancelc:
		if err == errRequestCanceled {
			err = errRequestCanceledResource
		}
		return nil, err
	}
}

// Get creates a new PersistResource to the target as specified in the key.
//  If this doesn't return an error, the PersistResource is ready to write requests to.
func (t *LruPool) Get(ctx context.Context, key interface{}) (v interface{}, put context.CancelFunc) {
	pc, _ := t.GetOrError(ctx, key)
	put = func() {
		pc.Put()
	}
	return pc.Get(), put
}

func (t *LruPool) Put(presource *PersistResource) {
	t.putOrCloseIdleResource(presource)
}

func (t *LruPool) putOrCloseIdleResource(presource *PersistResource) {
	if t == nil {
		return
	}
	if err := t.tryPutIdleResource(presource); err != nil {
		presource.close(err)
	}
}

func (t *LruPool) maxIdleResourcesPerKey() int {
	if v := t.MaxIdleResourcesPerHost; v != 0 {
		return v
	}
	return DefaultMaxIdleResourcesPerHost
}

// tryPutIdleResource adds presource to the list of idle persistent resources awaiting
// a new request.
// If presource is no longer needed or not in a good state, tryPutIdleResource returns
// an error explaining why it wasn't registered.
// tryPutIdleResource does not close presource. Use putOrCloseIdleResource instead for that.
func (t *LruPool) tryPutIdleResource(presource *PersistResource) error {
	if t.DisableKeepAlives || t.MaxIdleResourcesPerHost < 0 {
		return errKeepAlivesDisabled
	}
	if presource.isBroken() {
		return errResourceBroken
	}
	presource.markReused()

	t.idleMu.Lock()
	defer t.idleMu.Unlock()

	// Deliver presource to goroutine waiting for idle resource, if any.
	// (They may be actively dialing, but this conn is ready first.
	// Chrome calls this socket late binding.
	// See syncs://www.chromium.org/developers/design-documents/network-stack#TOC-Resourceection-Management.)
	key := presource.cacheKey
	if q, ok := t.idleResourceWait[key]; ok {
		done := false
		// HTTP/1.
		// Loop over the waiting list until we find a w that isn't done already, and hand it presource.
		for q.len() > 0 {
			w := q.popFront()
			if w.tryDeliver(presource, nil) {
				done = true
				break
			}
		}
		if q.len() == 0 {
			delete(t.idleResourceWait, key)
		} else {
			t.idleResourceWait[key] = q
		}
		if done {
			return nil
		}
	}

	if t.closeIdle {
		return errCloseIdle
	}
	if t.idleResource == nil {
		t.idleResource = make(map[targetKey][]*PersistResource)
	}
	idles := t.idleResource[key]
	if len(idles) >= t.maxIdleResourcesPerKey() {
		return errTooManyIdleResource
	}
	for _, exist := range idles {
		if exist == presource {
			log.Fatalf("dup idle presource %p in freelist", presource)
		}
	}
	t.idleResource[key] = append(idles, presource)
	t.idleLRU.add(presource)
	if t.MaxIdleResources != 0 && t.idleLRU.len() > t.MaxIdleResources {
		oldest := t.idleLRU.removeOldest()
		oldest.close(errTooManyIdle)
		t.removeIdleResourceLocked(oldest)
	}

	// Set idle timer, but only for HTTP/1 (presource.alt == nil).
	// The HTTP/2 implementation manages the idle timer itself
	// (see idleResourceTimeout in h2_bundle.go).
	if t.IdleResourceTimeout > 0 {
		if presource.idleTimer != nil {
			presource.idleTimer.Reset(t.IdleResourceTimeout)
		} else {
			presource.idleTimer = time.AfterFunc(t.IdleResourceTimeout, presource.closeResourceIfStillIdle)
		}
	}
	presource.idleAt = time.Now()
	return nil
}

// queueForIdleResource queues w to receive the next idle resource for w.cm.
// As an optimization hint to the caller, queueForIdleResource reports whether
// it successfully delivered an already-idle resource.
func (t *LruPool) queueForIdleResource(w *wantResource) (delivered bool) {
	if t.DisableKeepAlives {
		return false
	}

	t.idleMu.Lock()
	defer t.idleMu.Unlock()

	// Stop closing resources that become idle - we might want one.
	// (That is, undo the effect of t.CloseIdleResources.)
	t.closeIdle = false

	if w == nil {
		// Happens in test hook.
		return false
	}

	// If IdleResourceTimeout is set, calculate the oldest
	// PersistResource.idleAt time we're willing to use a cached idle
	// conn.
	var oldTime time.Time
	if t.IdleResourceTimeout > 0 {
		oldTime = time.Now().Add(-t.IdleResourceTimeout)
	}

	// Look for most recently-used idle resource.
	if list, ok := t.idleResource[w.key]; ok {
		stop := false
		delivered := false
		for len(list) > 0 && !stop {
			presource := list[len(list)-1]

			// See whether this resource has been idle too long, considering
			// only the wall time (the Round(0)), in case this is a laptop or VM
			// coming out of suspend with previously cached idle resources.
			tooOld := !oldTime.IsZero() && presource.idleAt.Round(0).Before(oldTime)
			if tooOld {
				// Async cleanup. Launch in its own goroutine (as if a
				// time.AfterFunc called it); it acquires idleMu, which we're
				// holding, and does a synchronous net.Resource.Close.
				go presource.closeResourceIfStillIdle()
			}
			if presource.isBroken() || tooOld {
				// If either PersistResource.readLoop has marked the resource
				// broken, but LruPool.removeIdleResource has not yet removed it
				// from the idle list, or if this PersistResource is too old (it was
				// idle too long), then ignore it and look for another. In both
				// cases it's already in the process of being closed.
				list = list[:len(list)-1]
				continue
			}
			delivered = w.tryDeliver(presource, nil)
			if delivered {
				// HTTP/1: only one client can use presource.
				// Remove it from the list.
				t.idleLRU.remove(presource)
				list = list[:len(list)-1]
			}
			stop = true
		}
		if len(list) > 0 {
			t.idleResource[w.key] = list
		} else {
			delete(t.idleResource, w.key)
		}
		if stop {
			return delivered
		}
	}

	// Register to receive next resource that becomes idle.
	if t.idleResourceWait == nil {
		t.idleResourceWait = make(map[targetKey]wantResourceQueue)
	}
	q := t.idleResourceWait[w.key]
	q.cleanFront()
	q.pushBack(w)
	t.idleResourceWait[w.key] = q
	return false
}

// removeIdleResource marks presource as dead.
func (t *LruPool) removeIdleResource(presource *PersistResource) bool {
	t.idleMu.Lock()
	defer t.idleMu.Unlock()
	return t.removeIdleResourceLocked(presource)
}

// t.idleMu must be held.
func (t *LruPool) removeIdleResourceLocked(presource *PersistResource) bool {
	if presource.idleTimer != nil {
		presource.idleTimer.Stop()
	}
	t.idleLRU.remove(presource)
	key := presource.cacheKey
	presources := t.idleResource[key]
	var removed bool
	switch len(presources) {
	case 0:
		// Nothing
	case 1:
		if presources[0] == presource {
			delete(t.idleResource, key)
			removed = true
		}
	default:
		for i, v := range presources {
			if v != presource {
				continue
			}
			// Slide down, keeping most recently-used
			// conns at the end.
			copy(presources[i:], presources[i+1:])
			t.idleResource[key] = presources[:len(presources)-1]
			removed = true
			break
		}
	}
	return removed
}

// queueForDial queues w to wait for permission to begin dialing.
// Once w receives permission to dial, it will do so in a separate goroutine.
func (t *LruPool) queueForDial(w *wantResource) {
	if t.MaxResourcesPerHost <= 0 {
		go t.dialResourceFor(w)
		return
	}

	t.connsPerHostMu.Lock()
	defer t.connsPerHostMu.Unlock()

	if n := t.connsPerHost[w.key]; n < t.MaxResourcesPerHost {
		if t.connsPerHost == nil {
			t.connsPerHost = make(map[targetKey]int)
		}
		t.connsPerHost[w.key] = n + 1
		go t.dialResourceFor(w)
		return
	}

	if t.connsPerHostWait == nil {
		t.connsPerHostWait = make(map[targetKey]wantResourceQueue)
	}
	q := t.connsPerHostWait[w.key]
	q.cleanFront()
	q.pushBack(w)
	t.connsPerHostWait[w.key] = q
}

// dialResourceFor dials on behalf of w and delivers the result to w.
// dialResourceFor has received permission to dial w.cm and is counted in t.connCount[w.cm.key()].
// If the dial is cancelled or unsuccessful, dialResourceFor decrements t.connCount[w.cm.key()].
func (t *LruPool) dialResourceFor(w *wantResource) {

	pc, err := t.buildResource(w.ctx, w.key)
	delivered := w.tryDeliver(pc, err)
	if err == nil && (!delivered) {
		// presource was not passed to w,
		// or it is HTTP/2 and can be shared.
		// Add to the idle resource pool.
		t.putOrCloseIdleResource(pc)
	}
	if err != nil {
		t.decResourcesPerHost(w.key)
	}
}

// decResourcesPerHost decrements the per-host resource count for key,
// which may in turn give a different waiting goroutine permission to dial.
func (t *LruPool) decResourcesPerHost(key targetKey) {
	if t.MaxResourcesPerHost <= 0 {
		return
	}

	t.connsPerHostMu.Lock()
	defer t.connsPerHostMu.Unlock()
	n := t.connsPerHost[key]
	if n == 0 {
		// Shouldn't happen, but if it does, the counting is buggy and could
		// easily lead to a silent deadlock, so report the problem loudly.
		panic("sync: internal error: connCount underflow")
	}

	// Can we hand this count to a goroutine still waiting to dial?
	// (Some goroutines on the wait list may have timed out or
	// gotten a resource another way. If they're all gone,
	// we don't want to kick off any spurious dial operations.)
	if q := t.connsPerHostWait[key]; q.len() > 0 {
		done := false
		for q.len() > 0 {
			w := q.popFront()
			if w.waiting() {
				go t.dialResourceFor(w)
				done = true
				break
			}
		}
		if q.len() == 0 {
			delete(t.connsPerHostWait, key)
		} else {
			// q is a value (like a slice), so we have to store
			// the updated q back into the map.
			t.connsPerHostWait[key] = q
		}
		if done {
			return
		}
	}

	// Otherwise, decrement the recorded count.
	if n--; n == 0 {
		delete(t.connsPerHost, key)
	} else {
		t.connsPerHost[key] = n
	}
}

func (t *LruPool) buildResource(ctx context.Context, key interface{}) (presource *PersistResource, err error) {
	presource = &PersistResource{
		t:        t,
		cacheKey: key,
	}

	if t.New != nil {
		presource.object = t.New(ctx)
	}
	return presource, nil
}
func (t *LruPool) setReqCanceler(key cancelKey, fn func(error)) {
	t.reqMu.Lock()
	defer t.reqMu.Unlock()
	if t.reqCanceler == nil {
		t.reqCanceler = make(map[cancelKey]func(error))
	}
	if fn != nil {
		t.reqCanceler[key] = fn
	} else {
		delete(t.reqCanceler, key)
	}
}

// replaceReqCanceler replaces an existing cancel function. If there is no cancel function
// for the request, we don't set the function and return false.
// Since CancelRequest will clear the canceler, we can use the return value to detect if
// the request was canceled since the last setReqCancel call.
func (t *LruPool) replaceReqCanceler(key cancelKey, fn func(error)) bool {
	t.reqMu.Lock()
	defer t.reqMu.Unlock()
	_, ok := t.reqCanceler[key]
	if !ok {
		return false
	}
	if fn != nil {
		t.reqCanceler[key] = fn
	} else {
		delete(t.reqCanceler, key)
	}
	return true
}

// PersistResource wraps a resource, usually a persistent one
// (but may be used for non-keep-alive requests as well)
type PersistResource struct {
	t        *LruPool
	cacheKey targetKey
	object   interface{}

	// Both guarded by LruPool.idleMu:
	idleAt    time.Time   // time it last become idle
	idleTimer *time.Timer // holding an AfterFunc to close it

	mu     sync.Mutex // guards following fields
	closed error      // set non-nil when conn is closed, before closech is closed
	broken bool       // an error has happened on this resource; marked broken so it's not reused.
	reused bool       // whether conn has had successful request/response and is being reused.
} // isBroken reports whether this resource is in a known broken state.

func (pc *PersistResource) Get() interface{} {
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

// close closes the underlying TCP resource and closes
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
		pc.t.decResourcesPerHost(pc.cacheKey)
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
