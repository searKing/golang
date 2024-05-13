// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashset

import "errors"

// HashSet is an auto make set
type HashSet struct {
	m map[any]struct{}
}

func New() *HashSet {
	return &HashSet{}
}

// Init initializes or clears map m.
func (m *HashSet) Init() *HashSet {
	m.m = make(map[any]struct{})
	return m
}

// lazyInit lazily initializes a zero map value.
func (m *HashSet) lazyInit() {
	if m.m == nil {
		m.Init()
	}
}

func (m *HashSet) Keys() []any {
	keys := []any{}
	for key, _ := range m.m {
		keys = append(keys, key)
	}
	return keys
}

// add adds Key to the head of the linked list.
func (m *HashSet) Add(key any) error {
	m.lazyInit()
	if _, ok := m.m[key]; ok {
		return errors.New("Key was already in HashSet")
	}
	m.m[key] = struct{}{}
	return nil
}

func (m *HashSet) AddOrUpdate(key any, value any) {
	m.Remove(key)
	m.Add(key)
}

// Remove removes Key from cl.
func (m *HashSet) Remove(key any) any {
	if v, ok := m.m[key]; ok {
		delete(m.m, key)
		return v
	}
	return nil
}

func (m *HashSet) Clear() {
	m.m = nil
}
func (m *HashSet) Find(key any) bool {
	return m.Contains(key)
}

func (m *HashSet) Contains(key any) bool {
	_, ok := m.m[key]
	return ok
}

func (m *HashSet) Peek(key any) bool {
	if m.Contains(key) {
		m.Remove(key)
		return true
	}
	return false
}

// Len returns the number of items in the cache.
func (m *HashSet) Len() int {
	return len(m.m)
}

func (m *HashSet) Count() int {
	return m.Len()
}
func (m *HashSet) Size() int {
	return m.Len()
}
