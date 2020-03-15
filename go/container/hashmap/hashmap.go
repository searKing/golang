// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashmap

import (
	"errors"
)

// HashMap is an auto make map
type HashMap struct {
	m map[interface{}]interface{}
}
type Pair struct {
	Key   interface{}
	Value interface{}
}

func New() *HashMap {
	return &HashMap{}
}

// Init initializes or clears map m.
func (m *HashMap) Init() *HashMap {
	m.m = make(map[interface{}]interface{})
	return m
}

// lazyInit lazily initializes a zero map value.
func (m *HashMap) lazyInit() {
	if m.m == nil {
		m.Init()
	}
}

func (m *HashMap) Keys() []interface{} {
	var keys []interface{}
	for key, _ := range m.m {
		keys = append(keys, key)
	}
	return keys
}
func (m *HashMap) Values() []interface{} {
	var values []interface{}
	for _, value := range m.m {
		values = append(values, value)
	}
	return values
}

func (m *HashMap) Pairs() []Pair {
	var pairs []Pair
	for key, value := range m.m {
		pairs = append(pairs, Pair{
			Key:   key,
			Value: value,
		})
	}
	return pairs
}
func (m *HashMap) AddPair(pair Pair) error {
	return m.Add(pair.Key, pair.Value)
}

// add adds Key to the head of the linked list.
func (m *HashMap) Add(key, value interface{}) error {
	m.lazyInit()
	if _, ok := m.m[key]; ok {
		return errors.New("Key was already in HashMap")
	}
	m.m[key] = value
	return nil
}

func (m *HashMap) AddOrUpdate(key interface{}, value interface{}) {
	m.Remove(key)
	m.Add(key, value)
}

// Remove removes Key from cl.
func (m *HashMap) Remove(key interface{}) interface{} {
	if v, ok := m.m[key]; ok {
		delete(m.m, key)
		return v
	}
	return nil
}

func (m *HashMap) Clear() {
	m.m = nil
}
func (m *HashMap) Find(key interface{}) (interface{}, bool) {
	v, ok := m.m[key]
	return v, ok
}

func (m *HashMap) Contains(key interface{}) bool {
	_, ok := m.m[key]
	return ok
}

func (m *HashMap) Peek(key interface{}) (interface{}, bool) {
	v, ok := m.m[key]
	if ok {
		m.Remove(key)
	}
	return v, ok
}

// Len returns the number of items in the cache.
func (m *HashMap) Len() int {
	return len(m.m)
}

func (m *HashMap) Count() int {
	return m.Len()
}
func (m *HashMap) Size() int {
	return m.Len()
}
