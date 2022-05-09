// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"container/list"
	"sync"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback[K comparable, V any] func(key K, value V)

// LRU implements a non-thread safe fixed size LRU cache
type LRU[K comparable, V any] struct {
	size      int
	evictList *list.List
	items     map[K]*list.Element
	onEvict   EvictCallback[K, V]
	mu        sync.Mutex
}

// entry is used to hold a value in the evictList
type entry[K comparable, V any] struct {
	key   K
	value V
}

// NewLRU constructs an LRU of the given size
func NewLRU[K comparable, V any](size int) *LRU[K, V] {
	c := &LRU[K, V]{
		size: size,
	}
	return c.Init()
}

// SetEvictCallback sets a callback when a cache entry is evicted
func (c *LRU[K, V]) SetEvictCallback(onEvict EvictCallback[K, V]) *LRU[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onEvict = onEvict
	return c
}

// Init initializes or clears LRU l.
func (c *LRU[K, V]) Init() *LRU[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictList = list.New()
	c.items = make(map[K]*list.Element)
	return c
}

// Purge is used to completely clear the cache.
func (c *LRU[K, V]) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entry[K, V]).value)
		}
		delete(c.items, k)
	}
	c.evictList.Init()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *LRU[K, V]) Add(key K, value V) (evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Check for existing item
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry[K, V]).value = value
		return false
	}

	// Add new item
	ent := &entry[K, V]{key, value}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry

	evict := c.evictList.Len() > c.size
	// Verify size not exceeded
	if evict {
		c.removeOldest()
	}
	return evict
}

// Get looks up a key's value from the cache.
func (c *LRU[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		if ent.Value.(*entry[K, V]) == nil {
			var zero V
			return zero, false
		}
		return ent.Value.(*entry[K, V]).value, true
	}
	return
}

// Contains checks if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *LRU[K, V]) Contains(key K) (ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok = c.items[key]
	return ok
}

// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *LRU[K, V]) Peek(key K) (value V, ok bool) {
	var ent *list.Element
	if ent, ok = c.items[key]; ok {
		return ent.Value.(*entry[K, V]).value, true
	}
	var zero V
	return zero, ok
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *LRU[K, V]) Remove(key K) (present bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU[K, V]) RemoveOldest() (key K, value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
		kv := ent.Value.(*entry[K, V])
		return kv.key, kv.value, true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// GetOldest returns the oldest entry
func (c *LRU[K, V]) GetOldest() (key K, value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ent := c.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*entry[K, V])
		return kv.key, kv.value, true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *LRU[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]K, len(c.items))
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry[K, V]).key
		i++
	}
	return keys
}

// Len returns the number of items in the cache.
func (c *LRU[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.evictList.Len()
}

// Cap returns the capacity of the cache.
func (c *LRU[K, V]) Cap() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.size
}

// Resize changes the cache size.
func (c *LRU[K, V]) Resize(size int) (evicted int) {
	c.mu.Lock()
	defer c.mu.Unlock()
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

// removeOldest removes the oldest item from the cache.
func (c *LRU[K, V]) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
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
