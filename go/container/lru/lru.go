// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lru

import (
	"container/list"
	"errors"
	"sync"
)

// LRU takes advantage of list's sequence and map's efficient locate
type LRU struct {
	ll   *list.List // list.Element.Value type is of interface{}
	m    map[any]*list.Element
	once sync.Once
}
type Pair struct {
	Key   any
	Value any
}

// lazyInit lazily initializes a zero List value.
func (lru *LRU) lazyInit() {
	lru.once.Do(func() {
		lru.ll = &list.List{}
		lru.m = make(map[any]*list.Element)
	})
}

func (lru *LRU) Keys() []any {
	var keys []any
	for key := range lru.m {
		keys = append(keys, key)
	}
	return keys
}
func (lru *LRU) Values() []any {
	var values []any
	for _, value := range lru.m {
		values = append(values, value)
	}
	return values
}

func (lru *LRU) Pairs() []Pair {
	var pairs []Pair
	for key, value := range lru.m {
		pairs = append(pairs, Pair{
			Key:   key,
			Value: value,
		})
	}
	return pairs
}
func (lru *LRU) AddPair(pair Pair) error {
	return lru.Add(pair.Key, pair.Value)
}

// Add adds Key to the head of the linked list.
func (lru *LRU) Add(key any, value any) error {
	lru.lazyInit()
	ele := lru.ll.PushFront(Pair{
		Key:   key,
		Value: value,
	})
	if _, ok := lru.m[key]; ok {
		return errors.New("key was already in LRU")
	}
	lru.m[key] = ele
	return nil
}

func (lru *LRU) AddOrUpdate(key any, value any) error {
	lru.Remove(key)
	return lru.Add(key, value)
}

func (lru *LRU) RemoveOldest() any {
	if lru.ll == nil {
		return nil
	}
	ele := lru.ll.Back()
	pair := ele.Value.(Pair)
	v := lru.ll.Remove(ele)
	delete(lru.m, pair.Key)
	return v
}

// Remove removes Key from cl.
func (lru *LRU) Remove(key any) any {
	if ele, ok := lru.m[key]; ok {
		v := lru.ll.Remove(ele)
		delete(lru.m, key)
		return v
	}
	return nil
}

func (lru *LRU) Find(key any) (any, bool) {
	e, ok := lru.m[key]
	if !ok {
		return nil, ok
	}
	return e.Value.(Pair).Value, true
}

func (lru *LRU) Peek(key any) (any, bool) {
	e, ok := lru.m[key]
	if ok {
		lru.Remove(key)
	}
	return e.Value.(Pair).Value, ok
}

// Len returns the number of items in the cache.
func (lru *LRU) Len() int {
	return len(lru.m)
}
