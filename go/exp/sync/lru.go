// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"sync"

	"github.com/searKing/golang/go/exp/container/lru"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback[K comparable, V any] func(key K, value V)

// LRU implements a thread safe fixed size LRU cache.
// A LRU is safe for use by multiple goroutines simultaneously.
// A LRU must not be copied after first use.
type LRU[K comparable, V any] struct {
	c  *lru.LRU[K, V]
	mu sync.Mutex
}

// NewLRU constructs an LRU of the given size
func NewLRU[K comparable, V any](size int) *LRU[K, V] {
	return &LRU[K, V]{
		c: lru.New[K, V](size),
	}
}

// SetEvictCallback sets a callback when a cache entry is evicted
func (c *LRU[K, V]) SetEvictCallback(onEvict EvictCallback[K, V]) *LRU[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	if onEvict == nil {
		c.c.SetEvictCallback(nil)
	} else {
		c.c.SetEvictCallback(lru.EvictCallback[K, V](onEvict))
	}
	return c
}

// SetEvictCallbackFunc sets a callback when a cache entry is evicted
//
// Deprecated, use SetEvictCallback instead.
func (c *LRU[K, V]) SetEvictCallbackFunc(onEvict func(key K, value V)) *LRU[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c.SetEvictCallbackFunc(onEvict)
	return c
}

// Init initializes or clears LRU l.
func (c *LRU[K, V]) Init() *LRU[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c.Init()
	return c
}

// Purge is used to completely clear the cache.
func (c *LRU[K, V]) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c.Purge()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *LRU[K, V]) Add(key K, value V) (evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Add(key, value)
}

// Get looks up a key's value from the cache.
func (c *LRU[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Get(key)
}

// Contains checks if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *LRU[K, V]) Contains(key K) (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Contains(key)
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *LRU[K, V]) Peek(key K) (value V, ok bool) {
	return c.c.Peek(key)
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *LRU[K, V]) Remove(key K) (present bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Remove(key)
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.RemoveOldest()
}

// GetOldest returns the oldest entry
func (c *LRU[K, V]) GetOldest() (key K, value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.GetOldest()
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *LRU[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Keys()
}

// Len returns the number of items in the cache.
func (c *LRU[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Len()
}

// Cap returns the capacity of the cache.
func (c *LRU[K, V]) Cap() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Cap()
}

// Resize changes the cache size.
func (c *LRU[K, V]) Resize(size int) (evicted int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Resize(size)
}
