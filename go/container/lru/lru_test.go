// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lru_test

import (
	"cmp"
	"slices"
	"testing"

	"github.com/searKing/golang/go/container/lru"
)

func TestLRUPairsValues(t *testing.T) {
	var l lru.LRU
	l.Add("a", 1)
	l.Add("b", 2)

	pairs := l.Pairs()
	slices.SortFunc(pairs, func(a, b lru.Pair) int { return cmp.Compare(a.Key.(string), b.Key.(string)) })
	if want := []lru.Pair{{Key: "a", Value: 1}, {Key: "b", Value: 2}}; !slices.Equal(pairs, want) {
		t.Errorf("Pairs() = %v; want %v", pairs, want)
	}

	values := l.Values()
	slices.SortFunc(values, func(a, b any) int { return cmp.Compare(a.(int), b.(int)) })
	if want := []any{1, 2}; !slices.Equal(values, want) {
		t.Errorf("Values() = %v; want %v", values, want)
	}
}

func TestLRUAddDuplicated(t *testing.T) {
	var l lru.LRU
	if err := l.Add("a", 1); err != nil {
		t.Fatalf("Add(a, 1) = %v; want nil", err)
	}
	if err := l.Add("a", 2); err == nil {
		t.Errorf("Add(a, 2) = nil; want error")
	}
	if v, ok := l.Find("a"); !ok || v != 1 {
		t.Errorf("Find(a) = %v, %t; want 1, true", v, ok)
	}
	if v, want := l.RemoveOldest(), (lru.Pair{Key: "a", Value: 1}); v != want {
		t.Errorf("RemoveOldest() = %v; want %v", v, want)
	}
	if v := l.RemoveOldest(); v != nil {
		t.Errorf("RemoveOldest() = %v; want nil", v)
	}
}

func TestLRUEmpty(t *testing.T) {
	var l lru.LRU
	if v, ok := l.Peek("a"); v != nil || ok {
		t.Errorf("Peek(a) = %v, %t; want nil, false", v, ok)
	}
	l.Add("a", 1)
	l.Remove("a")
	if v := l.RemoveOldest(); v != nil {
		t.Errorf("RemoveOldest() = %v; want nil", v)
	}
	if v, ok := l.Peek("a"); v != nil || ok {
		t.Errorf("Peek(a) = %v, %t; want nil, false", v, ok)
	}
}

func TestKeyLRUAddDuplicated(t *testing.T) {
	var l lru.KeyLRU
	if err := l.Add("a"); err != nil {
		t.Fatalf("Add(a) = %v; want nil", err)
	}
	if err := l.Add("a"); err == nil {
		t.Errorf("Add(a) = nil; want error")
	}
	if v := l.RemoveOldest(); v != "a" {
		t.Errorf("RemoveOldest() = %v; want a", v)
	}
	if v := l.RemoveOldest(); v != nil {
		t.Errorf("RemoveOldest() = %v; want nil", v)
	}
}
