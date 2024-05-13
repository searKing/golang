// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lru

import (
	"container/list"
	"errors"
	"sync"
)

// KeyLRU takes advantage of list's sequence and map's efficient locate
type KeyLRU struct {
	ll   *list.List // list.Element.Value type is of interface{}
	m    map[any]*list.Element
	once sync.Once
}

// lazyInit lazily initializes a zero List value.
func (lru *KeyLRU) lazyInit() {
	lru.once.Do(func() {
		lru.ll = &list.List{}
		lru.m = make(map[any]*list.Element)
	})
}
func (lru *KeyLRU) Keys() []any {
	var keys []any
	for key := range lru.m {
		keys = append(keys, key)
	}
	return keys
}

// add adds Key to the head of the linked list.
func (lru *KeyLRU) Add(key any) error {
	lru.lazyInit()
	ele := lru.ll.PushFront(key)
	if _, ok := lru.m[key]; ok {
		return errors.New("key was already in LRU")
	}
	lru.m[key] = ele
	return nil
}
func (lru *KeyLRU) AddOrUpdate(key any) error {
	lru.Remove(key)
	return lru.Add(key)
}

func (lru *KeyLRU) RemoveOldest() any {
	if lru.ll == nil {
		return nil
	}
	ele := lru.ll.Back()
	key := ele.Value.(any)
	lru.ll.Remove(ele)
	delete(lru.m, key)
	return key
}

// Remove removes Key from cl.
func (lru *KeyLRU) Remove(key any) any {
	if ele, ok := lru.m[key]; ok {
		v := lru.ll.Remove(ele)
		delete(lru.m, key)
		return v
	}
	return nil
}

func (lru *KeyLRU) Find(key any) (any, bool) {
	e, ok := lru.m[key]
	return e, ok
}

func (lru *KeyLRU) Peek(key any) (any, bool) {
	e, ok := lru.m[key]
	if ok {
		lru.Remove(key)
	}
	return e, ok
}

// Len returns the number of items in the cache.
func (lru *KeyLRU) Len() int {
	return len(lru.m)
}
