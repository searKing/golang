// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"sync"

	"github.com/searKing/golang/go/exp/container/lru"
)

// EvictCallback is used to get a callback when a cache entry is evicted
// type EvictCallback[K comparable, V any] func(key K, value V)
type EvictCallback[K comparable, V any] lru.EvictCallback[K, V]

// LRU is like a Go map[K]V but implements a thread safe fixed size LRU cache.
// Loads, stores, and deletes run in amortized constant time.
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
	c.c.SetEvictCallback(lru.EvictCallback[K, V](onEvict))
	return c
}

// SetEvictCallbackFunc sets a callback when a cache entry is evicted
//
// Deprecated, use SetEvictCallback instead.
func (c *LRU[K, V]) SetEvictCallbackFunc(onEvict func(key K, value V)) *LRU[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c.SetEvictCallback(onEvict)
	return c
}

// Init initializes or clears LRU l.
func (c *LRU[K, V]) Init() *LRU[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c.Init()
	return c
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

// Purge is used to completely clear the cache.
func (c *LRU[K, V]) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c.Purge()
}

// Load returns the value stored in the cache for a key, or zero if no
// value is present.
// The ok result indicates whether value was found in the cache.
func (c *LRU[K, V]) Load(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Load(key)
}

// Get looks up a key's value from the cache,
// with updating the "recently used"-ness of the key.
func (c *LRU[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Get(key)
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *LRU[K, V]) Peek(key K) (value V, ok bool) {
	return c.c.Peek(key)
}

// Contains checks if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *LRU[K, V]) Contains(key K) (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Contains(key)
}

// Store sets the value for a key.
func (c *LRU[K, V]) Store(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c.Store(key, value)
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *LRU[K, V]) Add(key K, value V) (evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Add(key, value)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (c *LRU[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.LoadOrStore(key, value)
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (c *LRU[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.LoadAndDelete(key)
}

// Delete deletes the value for a key.
func (c *LRU[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c.Delete(key)
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *LRU[K, V]) Remove(key K) (present bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Remove(key)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (c *LRU[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Swap(key, value)
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (c *LRU[K, V]) CompareAndSwap(key K, old, new V) (swapped bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.CompareAndSwap(key, old, new)
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false (even if the old value is the nil interface value).
func (c *LRU[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.CompareAndDelete(key, old)
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *LRU[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.Keys()
}

// Range calls f sequentially for each key and value present in the lru from oldest to newest.
// If f returns false, range stops the iteration.
// Without updating the "recently used"-ness of the key.
func (c *LRU[K, V]) Range(f func(key K, value V) bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.c.Range(f)
}

// PeekOldest returns the value stored in the cache for the oldest entry, or zero if no
// value is present.
// The ok result indicates whether value was found in the cache.
// Without updating the "recently used"-ness of the key.
func (c *LRU[K, V]) PeekOldest() (key K, value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.PeekOldest()
}

// PeekAndDeleteOldest deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (c *LRU[K, V]) PeekAndDeleteOldest() (key K, value V, loaded bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.c.PeekAndDeleteOldest()
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
