// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lru

import (
	"container/list"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback[K comparable, V any] func(key K, value V)

type EvictCallbackFunc[K comparable, V any] interface {
	Evict(key K, value V)
}

// LRU is like a Go map[K]V but implements a non-thread safe fixed size LRU cache.
// Loads, stores, and deletes run in amortized constant time.
type LRU[K comparable, V any] struct {
	size int // LRU size limit

	evictList *list.List          // sequence order for lru: Latest, Old, Older, ..., Oldest
	items     map[K]*list.Element // index to element access accelerate
	onEvict   EvictCallback[K, V]
}

// entry is used to hold a value in the evictList
type entry[K comparable, V any] struct {
	key   K
	value V
}

// New constructs an LRU of the given size
func New[K comparable, V any](size int) *LRU[K, V] {
	c := &LRU[K, V]{
		size: size,
	}
	return c.Init()
}

// SetEvictCallback sets a callback when a cache entry is evicted
func (c *LRU[K, V]) SetEvictCallback(onEvict EvictCallback[K, V]) *LRU[K, V] {
	c.onEvict = onEvict
	return c
}

// SetEvictCallbackFunc sets a callback func when a cache entry is evicted
//
// Deprecated, use SetEvictCallback instead.
func (c *LRU[K, V]) SetEvictCallbackFunc(onEvict func(key K, value V)) *LRU[K, V] {
	c.onEvict = onEvict
	return c
}

// Init initializes or clears LRU l.
func (c *LRU[K, V]) Init() *LRU[K, V] {
	c.evictList = list.New()
	c.items = make(map[K]*list.Element)
	return c
}

// Len returns the number of items in the cache.
func (c *LRU[K, V]) Len() int {
	return c.evictList.Len()
}

// Cap returns the capacity of the cache.
func (c *LRU[K, V]) Cap() int {
	return c.size
}

// Resize changes the cache size.
func (c *LRU[K, V]) Resize(size int) (evicted int) {
	diff := c.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.removeOldest()
	}
	c.size = size
	return diff
}

// Purge is used to completely clear the cache.
func (c *LRU[K, V]) Purge() {
	for k, v := range c.items {
		delete(c.items, k)
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entry[K, V]).value)
		}
	}
	c.evictList.Init()
}

// Load returns the value stored in the cache for a key, or zero if no
// value is present.
// The ok result indicates whether value was found in the cache.
func (c *LRU[K, V]) Load(key K) (value V, ok bool) {
	return c.load(key, true)
}

// Get looks up a key's value from the cache,
// with updating the "recently used"-ness of the key.
func (c *LRU[K, V]) Get(key K) (value V, ok bool) {
	return c.Load(key)
}

// Peek returns the value stored in the cache for a key, or zero if no
// value is present.
// The ok result indicates whether value was found in the cache.
// Without updating the "recently used"-ness of the key.
func (c *LRU[K, V]) Peek(key K) (value V, ok bool) {
	return c.load(key, false)
}

// Contains reports whether key is within the cache.
// The ok result indicates whether value was found in the cache.
// Without updating the "recently used"-ness of the key.
func (c *LRU[K, V]) Contains(key K) (ok bool) {
	_, ok = c.Peek(key)
	return ok
}

func (c *LRU[K, V]) load(key K, update bool) (value V, ok bool) {
	if e, ok := c.items[key]; ok {
		if update {
			// update the "recently used"-ness.
			c.evictList.MoveToFront(e)
		}
		return e.Value.(*entry[K, V]).value, true
	}
	return
}

// Store sets the value for a key.
func (c *LRU[K, V]) Store(key K, value V) {
	_, _ = c.Swap(key, value)
}

// Add adds a value to the cache,
// with updating the "recently used"-ness of the key.
// Returns true if an eviction occurred.
func (c *LRU[K, V]) Add(key K, value V) (evicted bool) {
	full := c.evictList.Len() >= c.size
	_, loaded := c.Swap(key, value)

	return !loaded && full
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (c *LRU[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	if e, ok := c.items[key]; ok {
		// update the "recently used"-ness.
		c.evictList.MoveToFront(e)
		return e.Value.(*entry[K, V]).value, true
	}

	// Add new item and update the "recently used"-ness of the key.
	e := &entry[K, V]{key, value}
	c.items[key] = c.evictList.PushFront(e)

	evict := c.evictList.Len() > c.size
	// Verify size not exceeded
	if evict {
		c.removeOldest()
	}
	return value, false
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (c *LRU[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	if e, ok := c.items[key]; ok {
		c.removeElement(e)
		return e.Value.(*entry[K, V]).value, true
	}
	return
}

// Delete deletes the value for a key.
func (c *LRU[K, V]) Delete(key K) {
	c.LoadAndDelete(key)
}

// Remove removes the provided key from the cache, returning true if the
// key was contained.
func (c *LRU[K, V]) Remove(key K) (present bool) {
	_, present = c.LoadAndDelete(key)
	return present
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (c *LRU[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	// Check for existing item
	if e, ok := c.items[key]; ok {
		// update the "recently used"-ness.
		c.evictList.MoveToFront(e)
		previous, e.Value.(*entry[K, V]).value = e.Value.(*entry[K, V]).value, value
		loaded = true
		return previous, loaded
	}

	// Add new item and update the "recently used"-ness of the key.
	e := &entry[K, V]{key, value}
	// update the "recently used"-ness.
	entry := c.evictList.PushFront(e)
	c.items[key] = entry

	evict := c.evictList.Len() > c.size
	// Verify size not exceeded
	if evict {
		c.removeOldest()
	}
	return previous, loaded
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (c *LRU[K, V]) CompareAndSwap(key K, old, new V) (swapped bool) {
	// Check for existing item
	if e, ok := c.items[key]; ok && any(e.Value.(*entry[K, V]).value) == any(old) {
		// update the "recently used"-ness.
		c.evictList.MoveToFront(e)
		e.Value.(*entry[K, V]).value = new
		return true
	}
	return false
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false (even if the old value is the nil interface value).
func (c *LRU[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	// Check for existing item
	if e, ok := c.items[key]; ok && any(e.Value.(*entry[K, V]).value) == any(old) {
		c.removeElement(e)
		return true
	}
	return false
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
// Without updating the "recently used"-ness of the key.
func (c *LRU[K, V]) Keys() []K {
	keys := make([]K, len(c.items))
	i := 0
	for e := c.evictList.Back(); e != nil; e = e.Prev() {
		keys[i] = e.Value.(*entry[K, V]).key
		i++
	}
	return keys
}

// Range calls f sequentially for each key and value present in the lru from oldest to newest.
// If f returns false, range stops the iteration.
// Without updating the "recently used"-ness of the key.
func (c *LRU[K, V]) Range(f func(key K, value V) bool) {
	// Iterate through list and print its contents.
	for e := c.evictList.Back(); e != nil; e = e.Prev() {
		if !f(e.Value.(*entry[K, V]).key, e.Value.(*entry[K, V]).value) {
			break
		}
	}
}

// PeekOldest returns the value stored in the cache for the oldest entry, or zero if no
// value is present.
// The ok result indicates whether value was found in the cache.
// Without updating the "recently used"-ness of the key.
func (c *LRU[K, V]) PeekOldest() (key K, value V, ok bool) {
	e := c.evictList.Back()
	if e != nil {
		kv := e.Value.(*entry[K, V])
		return kv.key, kv.value, true
	}
	return
}

// PeekAndDeleteOldest deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (c *LRU[K, V]) PeekAndDeleteOldest() (key K, value V, loaded bool) {
	e := c.evictList.Back()
	if e != nil {
		c.removeElement(e)
		kv := e.Value.(*entry[K, V])
		return kv.key, kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	return c.PeekAndDeleteOldest()
}

// GetOldest returns the oldest entry, without updating the "recently used"-ness
// or deleting it for being stale.
// Without updating the "recently used"-ness of the key.
func (c *LRU[K, V]) GetOldest() (key K, value V, ok bool) {
	return c.PeekOldest()
}

// removeOldest removes the oldest item from the cache.
func (c *LRU[K, V]) removeOldest() {
	e := c.evictList.Back()
	if e != nil {
		c.removeElement(e)
	}
}

// removeElement is used to remove a given list element from the cache
func (c *LRU[K, V]) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry[K, V])
	delete(c.items, kv.key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}
