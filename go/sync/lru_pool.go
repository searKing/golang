// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"sync"
	"time"
)

// DefaultLruPool is new resources as needed and caches them for reuse by subsequent calls.
var DefaultLruPool = &LruPool{
	MaxIdleResources:    100,
	IdleResourceTimeout: 90 * time.Second,
}

// DefaultMaxIdleResourcesPerBucket is the default value of LruPool's
// MaxIdleResourcesPerBucket.
const DefaultMaxIdleResourcesPerBucket = 2

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
//
// Deprecated: Use ResourcePool, which is type safe, closes the resources it drops,
// and can lease a resource of any of several keys.
// LruPool is a thin wrapper around ResourcePool[any, any].
type LruPool struct {
	// New optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	New func(ctx context.Context, req any) (resp any, err error)

	// DisableKeepAlives, if true, disables keep-alives and
	// will only use the resource to the server for a single request.
	DisableKeepAlives bool

	// MaxIdleResources controls the maximum number of idle (keep-alive)
	// resources across all buckets. Zero means no limit.
	// If negative, no resource is kept idle.
	MaxIdleResources int

	// MaxIdleResourcesPerBucket, if non-zero, controls the maximum idle
	// (keep-alive) resources to keep per-bucket. If zero,
	// DefaultMaxIdleResourcesPerBucket is used. If negative,
	// no resource is kept idle.
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

	once sync.Once
	pool ResourcePool[any, any]
}

// lruPoolReqKey is the context key of the req passed to LruPool.New.
type lruPoolReqKey struct{}

func (t *LruPool) init() *ResourcePool[any, any] {
	t.once.Do(func() {
		maxIdle := t.MaxIdleResourcesPerBucket
		if maxIdle == 0 {
			maxIdle = DefaultMaxIdleResourcesPerBucket
		}
		t.pool = ResourcePool[any, any]{
			New: func(ctx context.Context, key any) (any, error) {
				if t.New == nil {
					return nil, nil
				}
				return t.New(ctx, ctx.Value(lruPoolReqKey{}))
			},
			MaxResourcesPerKey:     t.MaxResourcesPerBucket,
			MaxIdleResourcesPerKey: max(maxIdle, 0),
			MaxIdleResources:       t.MaxIdleResources,
			IdleResourceTimeout:    t.IdleResourceTimeout,
			DisableKeepAlives:      t.DisableKeepAlives || maxIdle < 0 || t.MaxIdleResources < 0,
		}
	})
	return &t.pool
}

// PersistResource wraps a resource leased from a LruPool.
type PersistResource struct {
	lease *Lease[any, any]
}

// Get returns the resource, or nil if pc is nil.
func (pc *PersistResource) Get() any {
	if pc == nil || pc.lease == nil {
		return nil
	}
	return pc.lease.Value()
}

// Put returns the resource to its LruPool. Only the first call has an effect.
func (pc *PersistResource) Put() {
	if pc == nil || pc.lease == nil {
		return
	}
	pc.lease.Release()
}

// GetByKeyOrError leases a resource of the bucket key, made by New with req if none is idle.
// If MaxResourcesPerBucket is reached, it waits for a resource to be put back, or ctx to be done.
// The resource must be put back by PersistResource.Put.
func (t *LruPool) GetByKeyOrError(ctx context.Context, key any, req any) (*PersistResource, error) {
	l, err := t.init().Get(context.WithValue(ctx, lruPoolReqKey{}, req), key)
	if err != nil {
		return nil, err
	}
	return &PersistResource{lease: l}, nil
}

// GetByKey is like GetByKeyOrError, but returns the resource and a function to put it back.
// On error, it returns nil and a no-op function.
func (t *LruPool) GetByKey(ctx context.Context, key any, req any) (v any, put context.CancelFunc) {
	pc, _ := t.GetByKeyOrError(ctx, key, req)
	return pc.Get(), pc.Put
}

// GetOrError is like GetByKeyOrError with req as the key, but returns the resource
// and a function to put it back. On error, it returns nil and a no-op function.
func (t *LruPool) GetOrError(ctx context.Context, req any) (v any, put context.CancelFunc, err error) {
	pc, err := t.GetByKeyOrError(ctx, req, req)
	return pc.Get(), pc.Put, err
}

// Get is like GetOrError, ignoring the error.
func (t *LruPool) Get(ctx context.Context, req any) (v any, put context.CancelFunc) {
	return t.GetByKey(ctx, req, req)
}

// Put puts presource back, as presource.Put does.
func (t *LruPool) Put(presource *PersistResource) {
	presource.Put()
}

// CloseIdleResources closes any resources which are now sitting idle.
// It does not interrupt any resources currently in use.
func (t *LruPool) CloseIdleResources() {
	t.init().CloseIdleResources()
}

// RemoveIdleResource closes presource if it is idle, i.e. put back and not leased since,
// reporting whether it was.
func (t *LruPool) RemoveIdleResource(presource *PersistResource) bool {
	if presource == nil || presource.lease == nil {
		return false
	}
	pool := t.init()
	if presource.lease.p != pool {
		return false
	}
	return pool.removeIdle(presource.lease.e)
}
