// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"container/list"
	"context"
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// ErrResourcePoolClosed is returned by the ResourcePool's Get, Do and Prefill
	// methods after a call to Shutdown.
	ErrResourcePoolClosed = errors.New("sync: ResourcePool closed")

	// ErrBadResource may be returned, possibly wrapped, by the function passed to
	// ResourcePool.Do to report that the resource is broken: Do then discards
	// the resource instead of returning it to the pool.
	ErrBadResource = errors.New("sync: bad resource")

	// ErrResourcePoolInvalidKey is returned by the ResourcePool's Get, Do and Prefill
	// methods for a key not equal to itself, such as NaN, which no lookup could find again.
	ErrResourcePoolInvalidKey = errors.New("sync: ResourcePool: key not equal to itself")

	errResourcePoolNilNew = errors.New("sync: ResourcePool.New is nil")
	errResourcePoolFull   = errors.New("sync: ResourcePool: MaxResourcesPerKey reached")
)

// ResourcePool is a bounded pool of long-lived resources of type T, such as
// model instances, device contexts or connections, sharded by a key of type K,
// such as a GPU id or a remote address.
//
// Unlike sync.Pool, a ResourcePool never drops a resource silently: every resource
// made by New is eventually passed to Close, when it is discarded, when it has
// been idle for too long, when there are too many idle resources, or on Shutdown.
//
// A resource is leased by Get and returned by Lease.Release, which suits
// stateful use across several calls; or it is leased for a single call by Do,
// which suits stateless use. Callers waiting for a resource are served in FIFO
// order as soon as any of their keys has a resource available.
//
// A resource that must be used on a single OS thread can be made on a Thread
// by New and used through it; the ResourcePool needs no knowledge of that.
//
// A key not equal to itself, such as NaN, is rejected with ErrResourcePoolInvalidKey;
// a key of an interface type holding a value that is not comparable causes a panic.
// Either way, the pool is left unchanged.
//
// A ResourcePool is safe for concurrent use by multiple goroutines.
// The zero value with New set is ready to use.
// The fields must not be changed after first use.
// A ResourcePool must not be copied after first use.
type ResourcePool[K comparable, T any] struct {
	// New makes a resource for key. It is required, and may be called concurrently.
	// ctx is the one passed to Get, Do or Prefill: as the caller waits for New,
	// New should return soon after ctx is done.
	New func(ctx context.Context, key K) (T, error)

	// Close, if not nil, releases a resource that the pool no longer keeps.
	// It is called once per resource, and may be called concurrently.
	Close func(key K, v T) error

	// MaxResourcesPerKey, if positive, limits the number of resources per key,
	// including those being made by New, leased and idle.
	// When the limit is reached, Get waits for a resource to be released.
	MaxResourcesPerKey int

	// MaxIdleResourcesPerKey, if positive, limits the number of idle resources
	// per key: when exceeded, the least recently used idle resource of the key is closed.
	MaxIdleResourcesPerKey int

	// MaxIdleResources, if positive, limits the number of idle resources of all keys:
	// when exceeded, the least recently used idle resource is closed.
	MaxIdleResources int

	// MaxConcurrentLeases is the number of leases that may share a resource at
	// once, for resources safe for concurrent use. A value less than 2 means
	// a resource is leased exclusively. An idle resource is always preferred
	// to sharing one, and sharing one to making a new one.
	MaxConcurrentLeases int

	// IdleResourceTimeout, if positive, is how long a resource may stay idle
	// before it is closed.
	IdleResourceTimeout time.Duration

	// DisableKeepAlives, if true, closes every resource once released,
	// so that each is used by a single lease.
	DisableKeepAlives bool

	mu        sync.Mutex
	buckets   map[K]*poolBucket[K, T]
	live      int       // resources being made, kept or closed
	idleLRU   list.List // idle *poolEntry of all keys, least recently used at front
	closeIdle bool      // CloseIdleResources was called, and Get has not been since
	closed    bool
	drained   chan struct{} // closed when closed and live == 0
	closeErrs []error       // errors from Close since Shutdown
}

// ResourcePoolStats reports the state of a ResourcePool for a key.
type ResourcePoolStats struct {
	Total   int // resources kept or being made
	Idle    int // resources with no lease
	Leases  int // leases outstanding
	Waiting int // callers of Get waiting for a resource, possibly also of other keys
}

type poolBucket[K comparable, T any] struct {
	pending int                // resources being made by New
	entries []*poolEntry[K, T] // resources kept
	idle    []*poolEntry[K, T] // most recently used at end
	waiters list.List          // *poolWaiter waiting for a resource of the key, possibly also of other keys; FIFO
}

func (b *poolBucket[K, T]) total() int { return b.pending + len(b.entries) }

func (b *poolBucket[K, T]) unused() bool { return b.total() == 0 && b.waiters.Len() == 0 }

type poolEntry[K comparable, T any] struct {
	key       K
	value     T
	leases    int
	discarded bool          // closed once its last lease is released
	idleGen   uint64        // incremented each time the resource is leased from idle, made idle or removed
	lruElem   *list.Element // in ResourcePool.idleLRU while idle
	timer     *time.Timer   // IdleResourceTimeout
}

// A poolWaiter is a caller of Get waiting for either a resource,
// a permission to make one, or an error.
//
// A poolWaiter is queued in the bucket of each of its keys, as it may be served by any of them.
// It waits only while none of its keys has a resource available, so it needs to be
// reconsidered only when one is made available for one of them.
type poolWaiter[K comparable, T any] struct {
	keys  []K             // without duplicates
	elems []*list.Element // in the waiters of the bucket of each of keys, while queued
	ready chan struct{}   // closed when delivered
	done  bool

	grant poolGrant[K, T]
	err   error
}

// A poolGrant is what a caller of Get is given: a leased resource,
// or else a reserved slot to make one for key.
type poolGrant[K comparable, T any] struct {
	e   *poolEntry[K, T]
	key K
}

// Lease is a resource leased from a ResourcePool.
// It must be returned by Release or Discard once no longer used: a Lease that
// is garbage collected stays leased, as its resource may still be in use.
type Lease[K comparable, T any] struct {
	p    *ResourcePool[K, T]
	e    *poolEntry[K, T]
	done atomic.Bool
}

// Value returns the leased resource.
// It must not be used after the lease is released.
func (l *Lease[K, T]) Value() T { return l.e.value }

// Key returns the key of the leased resource.
func (l *Lease[K, T]) Key() K { return l.e.key }

// Release returns the resource to the pool. Only the first call to Release or
// Discard has an effect.
func (l *Lease[K, T]) Release() { l.release(false) }

// Discard returns the resource to the pool, which closes it once no lease uses it,
// e.g. because it is broken. Only the first call to Release or Discard has an effect.
func (l *Lease[K, T]) Discard() { l.release(true) }

func (l *Lease[K, T]) release(discard bool) {
	if l.done.Swap(true) {
		return
	}
	p := l.p
	p.mu.Lock()
	toClose := p.releaseLocked(l.e, discard)
	p.mu.Unlock()
	p.closeEntries(toClose)
}

// Get leases a resource of any of keys, or of the zero key if keys is empty.
// It prefers, in order: an idle resource of the key with most idle resources;
// a resource shared with other leases, if MaxConcurrentLeases allows;
// a new resource made for the key with fewest resources.
// Otherwise it waits until a resource is available or ctx is done.
//
// Ties are broken by the order of keys. The lease must be released.
func (p *ResourcePool[K, T]) Get(ctx context.Context, keys ...K) (*Lease[K, T], error) {
	if p.New == nil {
		return nil, errResourcePoolNilNew
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		var zero K
		keys = []K{zero}
	}

	g, w, err := p.acquireOrEnqueue(keys)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return p.lease(ctx, g)
	}

	select {
	case <-w.ready:
	case <-ctx.Done():
		p.mu.Lock()
		if !w.done {
			w.done = true
			p.dequeueLocked(w)
			p.mu.Unlock()
			return nil, ctx.Err()
		}
		// Delivered meanwhile: give it back.
		var toClose []*poolEntry[K, T]
		if w.err == nil {
			if w.grant.e != nil {
				toClose = p.releaseLocked(w.grant.e, false)
			} else {
				p.unreserveLocked(w.grant.key)
			}
		}
		p.mu.Unlock()
		p.closeEntries(toClose)
		return nil, ctx.Err()
	}
	if w.err != nil {
		return nil, w.err
	}
	return p.lease(ctx, w.grant)
}

// acquireOrEnqueue returns what is available for keys, or else a waiter queued for them.
func (p *ResourcePool[K, T]) acquireOrEnqueue(keys []K) (g poolGrant[K, T], w *poolWaiter[K, T], err error) {
	if !validKeys(keys...) {
		return g, nil, ErrResourcePoolInvalidKey
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return g, nil, ErrResourcePoolClosed
	}
	p.closeIdle = false
	if g, ok := p.acquireLocked(keys); ok {
		return g, nil, nil
	}
	w = &poolWaiter[K, T]{keys: uniq(keys), ready: make(chan struct{})}
	p.enqueueLocked(w)
	return g, w, nil
}

// lease returns a Lease of what g grants.
func (p *ResourcePool[K, T]) lease(ctx context.Context, g poolGrant[K, T]) (*Lease[K, T], error) {
	if g.e == nil {
		return p.create(ctx, g.key)
	}
	return p.newLease(g.e), nil
}

// Do leases a resource as Get does, calls f with it and releases it.
// If f returns an error matching ErrBadResource, panics or calls runtime.Goexit,
// the resource is discarded instead. Do returns the error of f.
//
// Do panics if f is nil.
func (p *ResourcePool[K, T]) Do(ctx context.Context, f func(ctx context.Context, v T) error, keys ...K) (err error) {
	if f == nil {
		panic("sync: ResourcePool.Do called with nil function")
	}
	l, err := p.Get(ctx, keys...)
	if err != nil {
		return err
	}
	ok := false
	defer func() {
		if !ok || errors.Is(err, ErrBadResource) {
			l.Discard()
			return
		}
		l.Release()
	}()
	err = f(ctx, l.Value())
	ok = true
	return err
}

// Prefill makes n resources for key and keeps them idle.
// It returns the first error, keeping the resources already made;
// exceeding MaxResourcesPerKey is an error.
func (p *ResourcePool[K, T]) Prefill(ctx context.Context, key K, n int) error {
	if p.New == nil {
		return errResourcePoolNilNew
	}
	for range n {
		if err := p.reserve(key); err != nil {
			return err
		}
		l, err := p.create(ctx, key)
		if err != nil {
			return err
		}
		l.Release()
	}
	return nil
}

// reserve reserves a slot to make a resource for key, if the limit allows it.
func (p *ResourcePool[K, T]) reserve(key K) error {
	if !validKeys(key) {
		return ErrResourcePoolInvalidKey
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return ErrResourcePoolClosed
	}
	p.closeIdle = false
	if b := p.buckets[key]; p.MaxResourcesPerKey > 0 && b != nil && b.total() >= p.MaxResourcesPerKey {
		return errResourcePoolFull
	}
	p.reserveLocked(key)
	return nil
}

// Stats returns the state of the pool for each key in use.
func (p *ResourcePool[K, T]) Stats() map[K]ResourcePoolStats {
	p.mu.Lock()
	defer p.mu.Unlock()
	stats := make(map[K]ResourcePoolStats, len(p.buckets))
	for k, b := range p.buckets {
		s := ResourcePoolStats{Total: b.total(), Idle: len(b.idle), Waiting: b.waiters.Len()}
		for _, e := range b.entries {
			s.Leases += e.leases
		}
		stats[k] = s
	}
	return stats
}

// Shutdown closes the pool: Get, Do and Prefill fail with ErrResourcePoolClosed,
// including for callers waiting for a resource; idle resources are closed at once,
// and leased ones once released. Shutdown waits until all resources are closed,
// or ctx is done, and returns the errors of the calls to Close since the first
// call to Shutdown, joined, or ctx.Err(). It may be called multiple times and concurrently.
func (p *ResourcePool[K, T]) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	first := !p.closed
	var toClose []*poolEntry[K, T]
	if first {
		p.closed = true
		p.drained = make(chan struct{})
		for _, b := range p.buckets {
			for b.waiters.Len() > 0 {
				w := b.waiters.Front().Value.(*poolWaiter[K, T])
				p.dequeueLocked(w)
				w.done = true
				w.err = ErrResourcePoolClosed
				close(w.ready)
			}
		}
		for _, b := range p.buckets {
			for len(b.idle) > 0 {
				toClose = append(toClose, p.removeLocked(b.idle[0]))
			}
		}
		p.checkDrainedLocked()
	}
	drained := p.drained
	p.mu.Unlock()
	p.closeEntries(toClose)

	select {
	case <-drained:
	case <-ctx.Done():
		return ctx.Err()
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return errors.Join(p.closeErrs...)
}

// CloseIdleResources closes the idle resources of all keys, and those released
// until the next call to Get or Prefill, unless a caller is waiting for them.
// It does not interrupt resources in use. Unlike after Shutdown, the pool remains usable.
func (p *ResourcePool[K, T]) CloseIdleResources() {
	p.mu.Lock()
	p.closeIdle = true
	var toClose []*poolEntry[K, T]
	for _, b := range p.buckets {
		for len(b.idle) > 0 {
			toClose = append(toClose, p.removeLocked(b.idle[0]))
		}
	}
	p.mu.Unlock()
	p.closeEntries(toClose)
}

// removeIdle closes e if it is idle, reporting whether it was.
func (p *ResourcePool[K, T]) removeIdle(e *poolEntry[K, T]) bool {
	p.mu.Lock()
	if e.lruElem == nil {
		p.mu.Unlock()
		return false
	}
	toClose := []*poolEntry[K, T]{p.removeLocked(e)}
	p.mu.Unlock()
	p.closeEntries(toClose)
	return true
}

func (p *ResourcePool[K, T]) newLease(e *poolEntry[K, T]) *Lease[K, T] {
	return &Lease[K, T]{p: p, e: e}
}

// create makes a resource for key, whose slot has been reserved, and leases it.
// If New fails or panics, the slot is freed.
func (p *ResourcePool[K, T]) create(ctx context.Context, key K) (*Lease[K, T], error) {
	made := false
	defer func() {
		if !made {
			p.mu.Lock()
			p.unreserveLocked(key)
			p.mu.Unlock()
		}
	}()
	v, err := p.New(ctx, key)
	if err != nil {
		return nil, err
	}
	made = true

	p.mu.Lock()
	b := p.buckets[key]
	b.pending--
	e := &poolEntry[K, T]{key: key, value: v, leases: 1}
	b.entries = append(b.entries, e)
	if p.closed {
		toClose := p.releaseLocked(e, true)
		p.mu.Unlock()
		p.closeEntries(toClose)
		return nil, ErrResourcePoolClosed
	}
	p.serveWaitersLocked(key) // e may be shared
	p.mu.Unlock()
	return p.newLease(e), nil
}

// acquireLocked leases an available resource of any of keys, or reserves a slot
// to make one, reporting whether it could.
func (p *ResourcePool[K, T]) acquireLocked(keys []K) (g poolGrant[K, T], ok bool) {
	// An idle resource of the key with most idle resources.
	var best *poolBucket[K, T]
	for _, k := range keys {
		if b := p.buckets[k]; b != nil && len(b.idle) > 0 && (best == nil || len(b.idle) > len(best.idle)) {
			best = b
		}
	}
	if best != nil {
		e := best.idle[len(best.idle)-1]
		best.idle = best.idle[:len(best.idle)-1]
		p.idleLRU.Remove(e.lruElem)
		e.lruElem = nil
		e.idleGen++
		if e.timer != nil {
			e.timer.Stop()
			e.timer = nil
		}
		e.leases = 1
		return poolGrant[K, T]{e: e, key: e.key}, true
	}

	// A resource shared with the fewest other leases.
	if p.MaxConcurrentLeases > 1 && !p.DisableKeepAlives {
		var e *poolEntry[K, T]
		for _, k := range keys {
			b := p.buckets[k]
			if b == nil {
				continue
			}
			for _, c := range b.entries {
				if !c.discarded && c.leases > 0 && c.leases < p.MaxConcurrentLeases && (e == nil || c.leases < e.leases) {
					e = c
				}
			}
		}
		if e != nil {
			e.leases++
			return poolGrant[K, T]{e: e, key: e.key}, true
		}
	}

	// A new resource for the key with the fewest resources.
	found := false
	minTotal := 0
	for _, k := range keys {
		n := 0
		if b := p.buckets[k]; b != nil {
			n = b.total()
		}
		if p.MaxResourcesPerKey > 0 && n >= p.MaxResourcesPerKey {
			continue
		}
		if !found || n < minTotal {
			found, g.key, minTotal = true, k, n
		}
	}
	if !found {
		return g, false
	}
	p.reserveLocked(g.key)
	return g, true
}

// reserveLocked reserves a slot to make a resource for key.
func (p *ResourcePool[K, T]) reserveLocked(key K) {
	p.bucketLocked(key).pending++
	p.live++
}

// bucketLocked returns the bucket of key, adding it if needed.
func (p *ResourcePool[K, T]) bucketLocked(key K) *poolBucket[K, T] {
	if p.buckets == nil {
		p.buckets = make(map[K]*poolBucket[K, T])
	}
	b := p.buckets[key]
	if b == nil {
		b = &poolBucket[K, T]{}
		p.buckets[key] = b
	}
	return b
}

// enqueueLocked queues w in the bucket of each of its keys.
func (p *ResourcePool[K, T]) enqueueLocked(w *poolWaiter[K, T]) {
	for _, k := range w.keys {
		w.elems = append(w.elems, p.bucketLocked(k).waiters.PushBack(w))
	}
}

// dequeueLocked removes w from the bucket of each of its keys.
func (p *ResourcePool[K, T]) dequeueLocked(w *poolWaiter[K, T]) {
	for i, k := range w.keys {
		b := p.buckets[k]
		b.waiters.Remove(w.elems[i])
		p.gcBucketLocked(k, b)
	}
	w.elems = nil
}

// unreserveLocked frees a slot reserved by reserveLocked whose resource was not made.
func (p *ResourcePool[K, T]) unreserveLocked(key K) {
	b := p.buckets[key]
	b.pending--
	p.live--
	p.gcBucketLocked(key, b)
	p.serveWaitersLocked(key)
	p.checkDrainedLocked()
}

// releaseLocked returns a lease of e, and returns the entries to close.
func (p *ResourcePool[K, T]) releaseLocked(e *poolEntry[K, T], discard bool) (toClose []*poolEntry[K, T]) {
	e.leases--
	if discard {
		e.discarded = true
	}
	if e.leases > 0 {
		if !e.discarded {
			p.serveWaitersLocked(e.key) // a shared slot is free
		}
		return nil
	}
	if e.discarded || p.closed || p.DisableKeepAlives {
		toClose = append(toClose, p.removeLocked(e))
		p.serveWaitersLocked(e.key)
		return toClose
	}

	b := p.buckets[e.key]
	b.idle = append(b.idle, e)
	e.lruElem = p.idleLRU.PushBack(e)
	e.idleGen++
	if p.IdleResourceTimeout > 0 {
		gen := e.idleGen
		e.timer = time.AfterFunc(p.IdleResourceTimeout, func() { p.idleTimeout(e, gen) })
	}
	p.serveWaitersLocked(e.key)
	if p.closeIdle && e.lruElem != nil {
		return append(toClose, p.removeLocked(e))
	}
	if p.MaxIdleResourcesPerKey > 0 {
		for len(b.idle) > p.MaxIdleResourcesPerKey {
			toClose = append(toClose, p.removeLocked(b.idle[0]))
		}
	}
	if p.MaxIdleResources > 0 {
		for p.idleLRU.Len() > p.MaxIdleResources {
			toClose = append(toClose, p.removeLocked(p.idleLRU.Front().Value.(*poolEntry[K, T])))
		}
	}
	return toClose
}

func (p *ResourcePool[K, T]) idleTimeout(e *poolEntry[K, T], gen uint64) {
	p.mu.Lock()
	// idleGen changes whenever e is leased, made idle again or removed.
	if e.idleGen != gen {
		p.mu.Unlock()
		return
	}
	toClose := []*poolEntry[K, T]{p.removeLocked(e)}
	p.mu.Unlock()
	p.closeEntries(toClose)
}

// removeLocked removes e from the pool, freeing its slot, and returns it to be closed.
func (p *ResourcePool[K, T]) removeLocked(e *poolEntry[K, T]) *poolEntry[K, T] {
	e.idleGen++
	b := p.buckets[e.key]
	b.entries = slices.DeleteFunc(b.entries, func(o *poolEntry[K, T]) bool { return o == e })
	if e.lruElem != nil {
		p.idleLRU.Remove(e.lruElem)
		e.lruElem = nil
		b.idle = slices.DeleteFunc(b.idle, func(o *poolEntry[K, T]) bool { return o == e })
	}
	if e.timer != nil {
		e.timer.Stop()
		e.timer = nil
	}
	p.gcBucketLocked(e.key, b)
	return e
}

func (p *ResourcePool[K, T]) gcBucketLocked(key K, b *poolBucket[K, T]) {
	if b.unused() {
		delete(p.buckets, key)
	}
}

// serveWaitersLocked delivers what is available for key to the waiters of key, in FIFO order.
// It must be called whenever a resource of key is made available, by its release or removal;
// closing an idle resource needs not, as no caller waits for a key with idle resources.
func (p *ResourcePool[K, T]) serveWaitersLocked(key K) {
	if p.closed {
		return
	}
	b := p.buckets[key]
	if b == nil || b.waiters.Len() == 0 {
		return
	}
	keys := []K{key} // the waiters' other keys have nothing available
	for b.waiters.Len() > 0 {
		g, ok := p.acquireLocked(keys)
		if !ok {
			return
		}
		w := b.waiters.Front().Value.(*poolWaiter[K, T])
		p.dequeueLocked(w)
		w.done = true
		w.grant = g
		close(w.ready)
	}
}

// closeEntries closes entries removed from the pool, all of them even if Close panics.
// It must be called with no lock held.
func (p *ResourcePool[K, T]) closeEntries(entries []*poolEntry[K, T]) {
	if len(entries) == 0 {
		return
	}
	defer p.closeEntries(entries[1:])
	p.closeEntry(entries[0])
}

func (p *ResourcePool[K, T]) closeEntry(e *poolEntry[K, T]) {
	var err error
	defer func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.live--
		if err != nil && p.closed {
			p.closeErrs = append(p.closeErrs, err)
		}
		p.checkDrainedLocked()
	}()
	if p.Close != nil {
		err = p.Close(e.key, e.value)
	}
}

func (p *ResourcePool[K, T]) checkDrainedLocked() {
	if p.closed && p.live == 0 {
		select {
		case <-p.drained:
		default:
			close(p.drained)
		}
	}
}

// uniq returns a copy of keys without duplicates, in order.
func uniq[K comparable](keys []K) []K {
	out := make([]K, 0, len(keys))
	for _, k := range keys {
		if !slices.Contains(out, k) {
			out = append(out, k)
		}
	}
	return out
}

// validKeys reports whether every key is equal to itself, so that it can be found
// again in a map. Comparing a key of an interface type holding a value that is not
// comparable panics, before the pool is locked.
func validKeys[K comparable](keys ...K) bool {
	for _, k := range keys {
		if k != k {
			return false
		}
	}
	return true
}
