// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/searKing/golang/go/exp/container/queue"
	"github.com/searKing/golang/go/pragma"
)

const (
	starvationThresholdNs = 1e6

	UnlimitedResident = -1
	UnlimitedCapacity = 0
)

// FixedPool is a set of resident and temporary items that may be individually saved and
// retrieved.
//
// Any item stored in the Pool may be removed automatically at any time without
// notification. If the Pool holds the only reference when this happens, the
// item might be deallocated.
//
// A Pool is safe for use by multiple goroutines simultaneously.
type FixedPool[E any] struct {
	noCopy pragma.DoNotCopy

	length   atomic.Int64 // items available
	capacity atomic.Int64 // items allocated

	// fixed-size pool for keep-alive
	// localC + localQ + localV
	// [0, MaxResidentSize) + [MaxResidentSize, MaxCapSize) + [MaxCapSize, UnlimitedCapacity)
	localC chan *FixedPoolElement[E] // fixed-size pool for keep-alive
	mu     sync.Mutex
	localQ queue.Queue[*FixedPoolElement[E]] // temporary pool for allocated, <keep-alive> excluded
	localV sync.Pool                         // A second GC should drop the victim cache, try put into local first.

	pinChan chan struct{} // bell for Put or New

	// New optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	New func() E

	// MinResidentSize controls the minimum number of keep-alive items. items will be preallocated.
	MinResidentSize int
	// MaxResidentSize controls the maximum number of keep-alive items. Negative means no limit.
	MaxResidentSize int
	// MaxCapacity controls the maximum number of allocated items. Zero means no limit.
	MaxCapacity int
}

// NewFixedPool returns an initialized fixed pool.
// resident controls the maximum number of keep-alive items. Negative means no limit.
// cap controls the maximum number of allocated items. Zero means no limit.
func NewFixedPool[E any](f func() E, size int) *FixedPool[E] {
	p := &FixedPool[E]{
		New:             f,
		MinResidentSize: size,
		MaxResidentSize: size,
		MaxCapacity:     size,
	}
	return p.Init()
}

// NewCachedPool Creates a pool that creates new items as needed, but
// will reuse previously constructed items when they are available.
// the pool will reuse previously constructed items and items will never be dropped.
func NewCachedPool[E any](f func() E) *FixedPool[E] {
	p := &FixedPool[E]{
		New:             f,
		MinResidentSize: 0,
		MaxResidentSize: UnlimitedResident,
		MaxCapacity:     UnlimitedCapacity,
	}
	return p.Init()
}

// NewTempPool Creates a pool that creates new items as needed, but
// will be dropped at second GC if only referenced by the pool self.
// the pool will reuse previously constructed items when they are available and not dropped.
func NewTempPool[E any](f func() E) *FixedPool[E] {
	p := &FixedPool[E]{
		New:             f,
		MinResidentSize: 0,
		MaxResidentSize: 0,
		MaxCapacity:     UnlimitedCapacity,
	}
	return p.Init()
}

// Init initializes fixed pool l.
func (p *FixedPool[E]) Init() *FixedPool[E] {
	if p.MaxResidentSize < 0 {
		p.MaxCapacity = 0
	} else {
		p.MaxCapacity = max(p.MaxCapacity, p.MaxResidentSize, 0)
	}
	p.pinChan = make(chan struct{})
	p.localC = make(chan *FixedPoolElement[E], max(p.MaxResidentSize, 0))

	p.preallocAllResident()
	return p
}

func (p *FixedPool[E]) preallocAllResident() {
	if p.New != nil {
		xs := make([]*FixedPoolElement[E], 0, p.MinResidentSize)
		for i := 0; i < p.MinResidentSize; i++ {
			x := p.TryGet()
			// short circuit
			if x == nil {
				break
			}
			xs = append(xs, x)
		}
		for _, x := range xs {
			p.Put(x)
		}
	}
}

func (p *FixedPool[E]) signal() {
	select {
	case p.pinChan <- struct{}{}:
	default:
	}
}

// Len returns the len of pool, that is object len idle, allocated and still in cache
// The complexity is O(1).
func (p *FixedPool[E]) Len() int {
	return int(p.length.Load())
}

// Cap returns the capacity of pool, that is object len allocated
// The complexity is O(1).
func (p *FixedPool[E]) Cap() int { return int(p.capacity.Load()) }

// Emplace adds x to the pool.
// NOTE: Emplace may break through the len and cap boundaries, as x be allocated already.
func (p *FixedPool[E]) Emplace(x E) {
	p.Put(newFixedPoolElement(x, p).markAvailable(false))
}

// Put adds x to the pool.
func (p *FixedPool[E]) Put(x *FixedPoolElement[E]) (stored bool) {
	return p.put(x, true)
}

// TryPut adds x to the pool, .
func (p *FixedPool[E]) TryPut(x *FixedPoolElement[E]) (stored bool) {
	return p.put(x, false)
}

func (p *FixedPool[E]) put(x *FixedPoolElement[E], victim bool) (stored bool) {
	if x == nil {
		return
	}
	x.markAvailable(true)
	defer func() {
		if stored {
			p.signal()
		} else {
			x.markAvailable(false)
		}
	}()
	select {
	case p.localC <- x:
		return true
	default:
		return p.putSlow(x, victim)
	}
}

// Get selects an arbitrary item from the Pool, removes it from the
// Pool, and returns it to the caller.
// Get may choose to ignore the pool and treat it as empty.
// Callers should not assume any relation between values passed to Put and
// the values returned by Get.
//
// If Get would otherwise return nil and p.New is non-nil, Get returns
// the result of calling p.New.
//
// Get uses context.Background internally; to specify the context, use
// GetContext.
func (p *FixedPool[E]) Get() *FixedPoolElement[E] {
	e, _ := p.GetContext(context.Background())
	return e
}

// GetContext selects an arbitrary item from the Pool, removes it from the
// Pool, and returns it to the caller.
// Get may choose to ignore the pool and treat it as empty.
// Callers should not assume any relation between values passed to Put and
// the values returned by Get.
//
// If GetContext would otherwise return nil and p.New is non-nil, Get returns
// the result of calling p.New.
func (p *FixedPool[E]) GetContext(ctx context.Context) (*FixedPoolElement[E], error) {
	return p.get(ctx, -1)
}

func (p *FixedPool[E]) TryGet() *FixedPoolElement[E] {
	e, _ := p.get(context.Background(), 1)
	return e
}

func (p *FixedPool[E]) get(ctx context.Context, maxIter int) (*FixedPoolElement[E], error) {
	select {
	case e := <-p.localC:
		return e.markAvailable(false), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		p.mu.Lock()
		if p.localQ.Next() {
			e := p.localQ.PopFront()
			p.mu.Unlock()
			return e.markAvailable(false), nil
		}
		p.mu.Unlock()
	}
	return p.getSlow(ctx, maxIter)
}

func (p *FixedPool[E]) getSlow(ctx context.Context, maxIter int) (*FixedPoolElement[E], error) {
	if ctx == nil {
		panic("sync.FixedPool: nil Context")
	}
	var timer *time.Timer // for canceling TLS handshake
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()
	iter := 0
	for {
		select {
		case e := <-p.localC:
			return e.markAvailable(false), nil
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			p.mu.Lock()
			if p.localQ.Next() {
				e := p.localQ.PopFront()
				p.mu.Unlock()
				return e.markAvailable(false), nil
			}
			x, allocated := p.tryAllocateLocked()
			if allocated {
				p.mu.Unlock()
				return x.markAvailable(false), nil
			}
			p.mu.Unlock()
		}

		iter++
		if maxIter > 0 && iter >= maxIter {
			return nil, nil
		}
		if timer != nil {
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(starvationThresholdNs)
		} else {
			timer = time.NewTimer(starvationThresholdNs)
		}
		select {
		case e := <-p.localC:
			return e.markAvailable(false), nil
		case <-p.pinChan:
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

func (p *FixedPool[E]) putSlow(x *FixedPoolElement[E], victim bool) (stored bool) {
	if x == nil {
		return true
	}

	x.markAvailable(true)
	defer func() {
		if !stored {
			x.markAvailable(false)
		}
	}()
	select {
	case p.localC <- x:
		return true
	default:
	}

	p.mu.Lock()
	move := p.moveToVictimLocked()
	p.mu.Unlock()

	if move { // overcapacity
		if victim { // drop this element into victim cache for reuse
			// After one GC, the victim cache should keep them alive.
			// A second GC should drop the victim cache.
			p.localV.Put(x)
			return true
		}
		// drop this element directly as overcapacity
		return false
	}

	// Try to put this element into localC or localQ if localC is full.
	select {
	case p.localC <- x:
		return true
	default:
		p.mu.Lock()
		defer p.mu.Unlock()
		p.localQ.PushBack(x)
		return true
	}
}

func (p *FixedPool[E]) isCapacityUnLimited() bool { return p.MaxCapacity <= UnlimitedCapacity }
func (p *FixedPool[E]) isResidentUnLimited() bool { return p.MaxResidentSize <= UnlimitedResident }
func (p *FixedPool[E]) moveToVictimLocked() bool {
	// resident no limit
	if p.isResidentUnLimited() {
		return false
	}
	c := p.Cap()
	// cap and resident both has limit
	return c > p.MaxResidentSize
}

func (p *FixedPool[E]) tryAllocateLocked() (x *FixedPoolElement[E], allocated bool) {
	// Try to pop the head of the victim for temporal locality of
	// reuse.
	{
		x := p.localV.Get()
		if x != nil {
			return x.(*FixedPoolElement[E]), true
		}
	}

	if p.isCapacityUnLimited() || p.Cap() < p.MaxCapacity {
		if n := p.New; n != nil {
			x := newFixedPoolElement(n(), p)
			return x, true
		}
		return nil, true
	}
	return nil, false
}

type FixedPoolElement[E any] struct {
	// The value stored with this element.
	Value E

	available bool // available as idle for the pool
	pool      *FixedPool[E]
}

func newFixedPoolElement[E any](Value E, pool *FixedPool[E]) *FixedPoolElement[E] {
	e := &FixedPoolElement[E]{
		Value: Value,
		pool:  pool,
	}
	e.pool.capacity.Add(1)
	runtime.SetFinalizer(e, (*FixedPoolElement[E]).Finalize)
	return e.markAvailable(true)
}

func (e *FixedPoolElement[E]) Finalize() {
	stored := e.pool.putSlow(e, false)
	if stored {
		runtime.SetFinalizer(e, (*FixedPoolElement[E]).Finalize)
		return
	}
	e.markAvailable(false).pool.capacity.Add(-1)
	// no need for a finalizer anymore
	runtime.SetFinalizer(e, nil)
}

func (e *FixedPoolElement[E]) Get() E {
	if e == nil {
		var zeroE E
		return zeroE
	}
	return e.Value
}

func (e *FixedPoolElement[E]) markAvailable(available bool) *FixedPoolElement[E] {
	if e != nil {
		if available == e.available {
			return e
		}
		e.available = available
		if e.available {
			e.pool.length.Add(1)
		} else {
			e.pool.length.Add(-1)
		}
	}
	return e
}
