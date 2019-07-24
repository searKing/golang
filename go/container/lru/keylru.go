package lru

import "container/list"

// LRU takes advantage of list's sequence and map's efficient locate
type KeyLRU struct {
	ll *list.List // list.Element.Value type is of interface{}
	m  map[interface{}]*list.Element
}

func NewKeyLRU() *KeyLRU {
	return &KeyLRU{}
}

// Init initializes or clears list l.
func (lru *KeyLRU) Init() *KeyLRU {
	lru.ll = &list.List{}
	lru.m = make(map[interface{}]*list.Element)
	return lru
}

// lazyInit lazily initializes a zero List value.
func (lru *KeyLRU) lazyInit() {
	if lru.ll == nil {
		lru.Init()
	}
}
func (lru *KeyLRU) Keys() []interface{} {
	keys := []interface{}{}
	for key, _ := range lru.m {
		keys = append(keys, key)
	}
	return keys
}

// add adds Key to the head of the linked list.
func (lru *KeyLRU) Add(key interface{}) {
	lru.lazyInit()
	ele := lru.ll.PushFront(key)
	if _, ok := lru.m[key]; ok {
		panic("persistConn was already in LRU")
	}
	lru.m[key] = ele
}
func (lru *KeyLRU) AddOrUpdate(key interface{}, value interface{}) {
	lru.Remove(key)
	lru.Add(key)
}

func (lru *KeyLRU) RemoveOldest() interface{} {
	if lru.ll == nil {
		return nil
	}
	ele := lru.ll.Back()
	key := ele.Value.(interface{})
	lru.ll.Remove(ele)
	delete(lru.m, key)
	return key
}

// Remove removes Key from cl.
func (lru *KeyLRU) Remove(key interface{}) interface{} {
	if ele, ok := lru.m[key]; ok {
		v := lru.ll.Remove(ele)
		delete(lru.m, key)
		return v
	}
	return nil
}

func (lru *KeyLRU) Find(key interface{}) (interface{}, bool) {
	e, ok := lru.m[key]
	return e, ok
}

func (lru *KeyLRU) Peek(key interface{}) (interface{}, bool) {
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
