// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// borrowed from https://github.com/golang/go/blob/master/src/net/http/transport.go

package sync

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

var (
	errKeepAlivesDisabled  = errors.New("sync: putIdleResource: keep alives disabled")
	errResourceBroken      = errors.New("sync: putIdleResource: resource is in bad state")
	errCloseIdle           = errors.New("sync: putIdleResource: CloseIdleResources was called")
	errTooManyIdle         = errors.New("sync: putIdleResource: too many idle resources")
	errTooManyIdleResource = errors.New("sync: putIdleResource: too many idle resources for bucket")
	errCloseIdleResources  = errors.New("sync: CloseIdleResources called")

	errIdleResourceTimeout = errors.New("sync: idle resource timeout")
)

// DefaultLruPool is new resources as needed and caches them for reuse by subsequent calls.
var DefaultLruPool = &LruPool{
	MaxIdleResources:    100,
	IdleResourceTimeout: 90 * time.Second,
}

// DefaultMaxIdleResourcesPerBucket is the default value of LruPool's
// MaxIdleResourcesPerBucket.
const DefaultMaxIdleResourcesPerBucket = 2

type targetKey any

// LruPool is an implementation of sync.Pool with LRU.
//
// By default, LruPool caches resources for future re-use.
// This may leave many open resources when accessing many buckets.
// This behavior can be managed using LruPool's CloseIdleResources method
// and the MaxIdleResourcesPerBucket and DisableKeepAlives fields.
//
// LruPools should be reused instead of created as needed.
// LruPools are safe for concurrent use by multiple goroutines.
//
// A LruPool is a low-level primitive for making resources.
type LruPool struct {
	// New optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	New func(ctx context.Context, req any) (resp any, err error)

	idleMu           sync.Mutex
	closeIdle        bool                             // user has requested to close all idle resources
	idleResource     map[targetKey][]*PersistResource // most recently used at end
	idleResourceWait map[targetKey]wantResourceQueue  // waiting getResources
	idleLRU          resourceLRU

	resourcesPerBucketMu   sync.Mutex
	resourcesPerBucket     map[targetKey]int
	resourcesPerBucketWait map[targetKey]wantResourceQueue // waiting getResources
	// DisableKeepAlives, if true, disables keep-alives and
	// will only use the resource to the server for a single request.
	DisableKeepAlives bool

	// MaxIdleResources controls the maximum number of idle (keep-alive)
	// resources across all buckets. Zero means no limit.
	MaxIdleResources int

	// MaxIdleResourcesPerBucket, if non-zero, controls the maximum idle
	// (keep-alive) resources to keep per-bucket. If zero,
	// DefaultMaxIdleResourcesPerBucket is used.
	MaxIdleResourcesPerBucket int

	// MaxResourcesPerBucket optionally limits the total number of
	// resources per bucket, including resources in the newResource,
	// active, and idle states. On limit violation, news will block.
	//
	// Zero means no limit.
	MaxResourcesPerBucket int

	// IdleResourceTimeout is the maximum amount of time an idle
	// (keep-alive) resource will remain idle before closing
	// itself.
	// Zero means no limit.
	IdleResourceTimeout time.Duration
}

// GetByKeyOrError creates a new PersistResource to the target as specified in the key.
// If this doesn't return an error, the PersistResource is ready to write requests to.
func (t *LruPool) GetByKeyOrError(ctx context.Context, key any, req any) (pc *PersistResource, err error) {

	w := &wantResource{
		req:   req,
		key:   key,
		ctx:   ctx,
		ready: make(chan struct{}, 1),
	}
	defer func() {
		if err != nil {
			w.cancel(t, err)
		}
	}()

	// Queue for idle resource.
	if delivered := t.queueForIdleResource(w); delivered {
		pc := w.pr
		return pc, nil
	}

	// Queue for permission to new resource.
	t.queueForNewResource(w)

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
			default:
				// return below
			}
		}
		return w.pr, w.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetByKey creates a new PersistResource to the target as specified in the key.
// If this doesn't return an error, the PersistResource is ready to write requests to.
func (t *LruPool) GetByKey(ctx context.Context, key any, req any) (v any, put context.CancelFunc) {
	pc, _ := t.GetByKeyOrError(ctx, key, req)
	put = func() {
		pc.Put()
	}
	return pc.Get(), put
}

// GetOrError creates a new PersistResource to the target as specified in the key.
// If this doesn't return an error, the PersistResource is ready to write requests to.
func (t *LruPool) GetOrError(ctx context.Context, req any) (v any, put context.CancelFunc, err error) {
	pc, err := t.GetByKeyOrError(ctx, req, req)
	put = func() {
		pc.Put()
	}
	return pc.Get(), put, err
}

// Get creates a new PersistResource to the target as specified in the key.
// If this doesn't return an error, the PersistResource is ready to write requests to.
func (t *LruPool) Get(ctx context.Context, req any) (v any, put context.CancelFunc) {
	return t.GetByKey(ctx, req, req)
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

func (t *LruPool) maxIdleResourcesPerBucket() int {
	if v := t.MaxIdleResourcesPerBucket; v != 0 {
		return v
	}
	return DefaultMaxIdleResourcesPerBucket
}

// tryPutIdleResource adds presource to the list of idle persistent resources awaiting
// a new request.
// If presource is no longer needed or not in a good state, tryPutIdleResource returns
// an error explaining why it wasn't registered.
// tryPutIdleResource does not close presource. Use putOrCloseIdleResource instead for that.
func (t *LruPool) tryPutIdleResource(presource *PersistResource) error {
	if t.DisableKeepAlives || t.MaxIdleResourcesPerBucket < 0 {
		return errKeepAlivesDisabled
	}
	if presource.isBroken() {
		return errResourceBroken
	}
	presource.markReused()

	t.idleMu.Lock()
	defer t.idleMu.Unlock()

	// Deliver presource to goroutine waiting for idle resource, if any.
	// (They may be actively newResource, but this resource is ready first.
	// Chrome calls this socket late binding.
	// See syncs://www.chromium.org/developers/design-documents/network-stack#TOC-Resourceection-Management.)
	key := presource.cacheKey
	if q, ok := t.idleResourceWait[key]; ok {
		done := false
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
	if len(idles) >= t.maxIdleResourcesPerBucket() {
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
	// resource.
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
				// broken, but LruPool.RemoveIdleResource has not yet removed it
				// from the idle list, or if this PersistResource is too old (it was
				// idle too long), then ignore it and look for another. In both
				// cases it's already in the process of being closed.
				list = list[:len(list)-1]
				continue
			}
			delivered = w.tryDeliver(presource, nil)
			if delivered {
				// only one client can use presource.
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

// RemoveIdleResource marks presource as dead.
func (t *LruPool) RemoveIdleResource(presource *PersistResource) bool {
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
			// Slide down, keeping most recently-used resources at the end.
			copy(presources[i:], presources[i+1:])
			t.idleResource[key] = presources[:len(presources)-1]
			removed = true
			break
		}
	}
	return removed
}

// queueForNewResource queues w to wait for permission to begin newResource.
// Once w receives permission to dial, it will do so in a separate goroutine.
func (t *LruPool) queueForNewResource(w *wantResource) {
	if t.MaxResourcesPerBucket <= 0 {
		go t.newResourceFor(w)
		return
	}

	t.resourcesPerBucketMu.Lock()
	defer t.resourcesPerBucketMu.Unlock()

	if n := t.resourcesPerBucket[w.key]; n < t.MaxResourcesPerBucket {
		if t.resourcesPerBucket == nil {
			t.resourcesPerBucket = make(map[targetKey]int)
		}
		t.resourcesPerBucket[w.key] = n + 1
		go t.newResourceFor(w)
		return
	}

	if t.resourcesPerBucketWait == nil {
		t.resourcesPerBucketWait = make(map[targetKey]wantResourceQueue)
	}
	q := t.resourcesPerBucketWait[w.key]
	q.cleanFront()
	q.pushBack(w)
	t.resourcesPerBucketWait[w.key] = q
}

// newResourceFor news on behalf of w and delivers the result to w.
// newResourceFor has received permission to dial w.cm and is counted in t.resourceCount[w.cm.key()].
// If the dial is cancelled or unsuccessful, newResourceFor decrements t.resourceCount[w.cm.key()].
func (t *LruPool) newResourceFor(w *wantResource) {
	pc, err := t.buildResource(w.ctx, w.key, w.req)
	delivered := w.tryDeliver(pc, err)
	if err == nil && (!delivered) {
		// presource was not passed to w,
		// or it can be shared.
		// Add to the idle resource pool.
		t.putOrCloseIdleResource(pc)
	}
	if err != nil {
		t.decResourcesPerBucket(w.key)
	}
}

// decResourcesPerBucket decrements the per-bucket resource count for key,
// which may in turn give a different waiting goroutine permission to dial.
func (t *LruPool) decResourcesPerBucket(key targetKey) {
	if t.MaxResourcesPerBucket <= 0 {
		return
	}

	t.resourcesPerBucketMu.Lock()
	defer t.resourcesPerBucketMu.Unlock()
	n := t.resourcesPerBucket[key]
	if n == 0 {
		// Shouldn't happen, but if it does, the counting is buggy and could
		// easily lead to a silent deadlock, so report the problem loudly.
		panic("sync: internal error: resourceCount underflow")
	}

	// Can we hand this count to a goroutine still waiting to dial?
	// (Some goroutines on the wait list may have timed out or
	// gotten a resource another way. If they're all gone,
	// we don't want to kick off any spurious dial operations.)
	if q := t.resourcesPerBucketWait[key]; q.len() > 0 {
		done := false
		for q.len() > 0 {
			w := q.popFront()
			if w.waiting() {
				go t.newResourceFor(w)
				done = true
				break
			}
		}
		if q.len() == 0 {
			delete(t.resourcesPerBucketWait, key)
		} else {
			// q is a value (like a slice), so we have to store
			// the updated q back into the map.
			t.resourcesPerBucketWait[key] = q
		}
		if done {
			return
		}
	}

	// Otherwise, decrement the recorded count.
	if n--; n == 0 {
		delete(t.resourcesPerBucket, key)
	} else {
		t.resourcesPerBucket[key] = n
	}
}

func (t *LruPool) buildResource(ctx context.Context, key any, req any) (presource *PersistResource, err error) {
	presource = &PersistResource{
		t:        t,
		cacheKey: key,
	}

	if t.New != nil {
		presource.object, err = t.New(ctx, req)
	}
	return presource, err
}

// CloseIdleResources closes any connections which were previously
// connected from previous requests but are now sitting idle in
// a "keep-alive" state. It does not interrupt any connections currently
// in use.
func (t *LruPool) CloseIdleResources() {
	t.idleMu.Lock()
	m := t.idleResource
	t.idleResource = nil
	t.closeIdle = true // close newly idle connections
	t.idleLRU = resourceLRU{}
	t.idleMu.Unlock()
	for _, conns := range m {
		for _, pconn := range conns {
			pconn.close(errCloseIdleResources)
		}
	}
}
