// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"math/rand"
	"runtime"
	"slices"
	"sync"
	"testing"

	sync_ "github.com/searKing/golang/go/exp/sync"
)

func TestLRU(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) {
		if k != v {
			t.Fatalf("Evict values not equal (%v!=%v)", k, v)
		}
		evictCounter++
	}
	l := sync_.NewLRU[int, int](128)
	l.SetEvictCallback(onEvicted)

	for i := 0; i < 256; i++ {
		if i%2 == 0 {
			l.Store(i, i)
			continue
		}
		l.Add(i, i)
	}
	if l.Len() != 128 {
		t.Fatalf("bad len: %v", l.Len())
	}

	if evictCounter != 128 {
		t.Fatalf("bad evict count: %v", evictCounter)
	}

	for i, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || v != k || v != i+128 {
			t.Fatalf("bad key: %v", k)
		}
	}
	for i := 0; i < 128; i++ {
		_, ok := l.Get(i)
		if ok {
			t.Fatalf("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		_, ok := l.Get(i)
		if !ok {
			t.Fatalf("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		v, ok := l.LoadAndDelete(i)
		if !ok {
			t.Fatalf("should be contained")
		}
		if v != i {
			t.Fatalf("bad key: %v", i)
		}
		ok = l.Remove(i)
		if ok {
			t.Fatalf("should not be contained")
		}
		_, ok = l.Get(i)
		if ok {
			t.Fatalf("should be deleted")
		}
	}

	l.Get(192) // expect 192 to be last key in l.Keys()

	for i, k := range l.Keys() {
		if (i < 63 && k != i+193) || (i == 63 && k != 192) {
			t.Fatalf("out of order key: %v", k)
		}
	}

	l.Purge()
	if l.Len() != 0 {
		t.Fatalf("bad len: %v", l.Len())
	}
	if _, ok := l.Get(200); ok {
		t.Fatalf("should contain nothing")
	}
}

// Test that Resize can upsize and downsize
func TestLRU_Resize(t *testing.T) {
	onEvictCounter := 0
	onEvicted := func(k int, v int) { onEvictCounter++ }

	l := sync_.NewLRU[int, int](2).SetEvictCallback(onEvicted)

	// Downsize
	l.Add(1, 1)
	l.Add(2, 2)
	evicted := l.Resize(1)
	if evicted != 1 {
		t.Errorf("1 element should have been evicted: %v", evicted)
	}
	if onEvictCounter != 1 {
		t.Errorf("onEvicted should have been called 1 time: %v", onEvictCounter)
	}

	l.Add(3, 3)
	if l.Contains(1) {
		t.Errorf("Element 1 should have been evicted")
	}

	// Upsize
	evicted = l.Resize(2)
	if evicted != 0 {
		t.Errorf("0 elements should have been evicted: %v", evicted)
	}

	l.Add(4, 4)
	if !l.Contains(3) || !l.Contains(4) {
		t.Errorf("Cache should have contained 2 elements")
	}
}

// Test that Contains doesn't update recent-ness
func TestLRU_Contains(t *testing.T) {
	l := sync_.NewLRU[int, int](2)

	l.Add(1, 1)
	l.Add(2, 2)
	if !l.Contains(1) {
		t.Errorf("1 should be contained")
	}

	l.Add(3, 3)
	if l.Contains(1) {
		t.Errorf("Contains should not have updated recent-ness of 1")
	}
}

// Test that Peek doesn't update recent-ness
func TestLRU_Peek(t *testing.T) {
	l := sync_.NewLRU[int, int](2)

	l.Add(1, 1)
	l.Add(2, 2)
	if v, ok := l.Peek(1); !ok || v != 1 {
		t.Errorf("1 should be set to 1: %v, %v", v, ok)
	}

	l.Add(3, 3)
	if l.Contains(1) {
		t.Errorf("should not have updated recent-ness of 1")
	}
}

// Test that Add returns true/false if an eviction occurred
func TestLRU_Add(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](1).SetEvictCallback(onEvicted)

	if l.Add(1, 1) == true || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	if l.Add(1, -1) == true || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	if l.Add(2, 2) == false || evictCounter != 1 {
		t.Errorf("should have an eviction")
	}
}

func TestLRU_Store(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](1).SetEvictCallback(onEvicted)

	l.Store(1, 1)
	if l.Len() != 1 || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	l.Store(1, -1)
	if l.Len() != 1 || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	l.Store(2, 2)
	if l.Len() != 1 || evictCounter != 1 {
		t.Errorf("should have an eviction")
	}
}

func TestLRU_LoadOrStore(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](1).SetEvictCallback(onEvicted)

	{
		_, loaded := l.LoadOrStore(1, 1)
		if loaded {
			t.Errorf("should not loaded")
		}
		if evictCounter != 0 {
			t.Errorf("should not have an eviction")
		}
	}
	{
		old, loaded := l.LoadOrStore(1, -1)
		if !loaded {
			t.Errorf("should loaded")
		}
		if old != 1 {
			t.Errorf("should load old value 1")
		}
		if evictCounter != 0 {
			t.Errorf("should not have an eviction")
		}
	}

	{
		_, loaded := l.LoadOrStore(2, 2)
		if loaded {
			t.Errorf("should not loaded")
		}
		if evictCounter != 1 {
			t.Errorf("should not have an eviction")
		}
	}
}

func TestLRU_LoadAndDelete(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](1).SetEvictCallback(onEvicted)
	l.Add(1, 1)
	l.Add(2, 2)

	{
		_, loaded := l.LoadAndDelete(-1)
		if loaded {
			t.Errorf("should not loaded")
		}
	}
	{
		_, loaded := l.LoadAndDelete(1)
		if loaded {
			t.Errorf("should not loaded")
		}
	}

	{
		old, loaded := l.LoadAndDelete(2)
		if !loaded {
			t.Errorf("should not loaded")
		}
		if old != 2 {
			t.Errorf("should load old value 2")
		}
	}
}

func TestLRU_Delete(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](1).SetEvictCallback(onEvicted)
	l.Add(1, 1)
	l.Add(2, 2)

	{
		l.Delete(-1)
		if l.Len() != 1 {
			t.Errorf("should not deleted")
		}
	}
	{
		l.Delete(1)
		if l.Len() != 1 {
			t.Errorf("should not deleted")
		}
	}
	{
		l.Delete(2)
		if l.Len() != 0 {
			t.Errorf("should not deleted")
		}
	}
}

func TestLRU_Remove(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](1).SetEvictCallback(onEvicted)
	l.Add(1, 1)
	l.Add(2, 2)

	{
		loaded := l.Remove(-1)
		if loaded {
			t.Errorf("should not loaded")
		}
	}
	{
		loaded := l.Remove(1)
		if loaded {
			t.Errorf("should not loaded")
		}
	}

	{
		loaded := l.Remove(2)
		if !loaded {
			t.Errorf("should not loaded")
		}
	}
}

func TestLRU_Swap(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](1).SetEvictCallback(onEvicted)

	if _, loaded := l.Swap(1, 1); loaded == true || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	if pre, loaded := l.Swap(1, -1); loaded == false || pre != 1 || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	if _, loaded := l.Swap(2, 2); loaded == true || evictCounter != 1 {
		t.Errorf("should have an eviction")
	}
}

func TestLRU_CompareAndSwap(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](1).SetEvictCallback(onEvicted)
	l.Add(1, 1)
	l.Add(2, 2)

	if swapped := l.CompareAndSwap(-1, -1, -2); swapped {
		t.Errorf("should not swapped")
	}
	if swapped := l.CompareAndSwap(1, 1, -1); swapped {
		t.Errorf("should not swapped")
	}
	if swapped := l.CompareAndSwap(2, -2, 2); swapped {
		t.Errorf("should not swapped")
	}
	if swapped := l.CompareAndSwap(2, 2, -2); !swapped {
		t.Errorf("should swapped")
	} else {
		if val, ok := l.Peek(2); !ok || val != -2 {
			t.Errorf("should swapped with value -2")
		}
	}
}

func TestLRU_CompareAndDelete(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](1).SetEvictCallback(onEvicted)
	l.Add(1, 1)
	l.Add(2, 2)

	if deleted := l.CompareAndDelete(-1, -1); deleted {
		t.Errorf("should not deleted")
	}
	if deleted := l.CompareAndDelete(1, 1); deleted {
		t.Errorf("should not deleted")
	}
	if deleted := l.CompareAndDelete(2, -2); deleted {
		t.Errorf("should not deleted")
	}
	if deleted := l.CompareAndDelete(2, 2); !deleted {
		t.Errorf("should deleted")
	} else {
		if l.Len() != 0 {
			t.Errorf("should empty")
		}
	}
}

func TestLRU_Keys(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](2).SetEvictCallback(onEvicted)
	l.Add(1, 1)
	l.Add(2, 2)
	l.Add(3, 3)

	if !slices.Equal(l.Keys(), []int{2, 3}) {
		t.Fatalf("bad key order: %v", l.Keys())
	}
}

func TestLRU_Range(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k int, v int) { evictCounter++ }

	l := sync_.NewLRU[int, int](2).SetEvictCallback(onEvicted)
	l.Add(1, 1)
	l.Add(2, 2)
	l.Add(3, 3)

	var keys, vals []int
	l.Range(func(key int, value int) bool {
		keys = append(keys, key)
		vals = append(vals, value)
		return true
	})

	if !slices.Equal(l.Keys(), keys) {
		t.Fatalf("bad key order: %v", l.Keys())
	}
	if !slices.Equal(keys, vals) {
		t.Fatalf("mismatched kv pairs: %v:%v", keys, vals)
	}
}

func TestLRU_GetOldest(t *testing.T) {
	l := sync_.NewLRU[int, int](128)

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	k, _, ok := l.GetOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 128 {
		t.Fatalf("bad: %v", k)
	}
}

func TestLRU_GetOldest_RemoveOldest(t *testing.T) {
	l := sync_.NewLRU[int, int](128)

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	k, _, ok := l.GetOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 128 {
		t.Fatalf("bad: %v", k)
	}

	k, _, ok = l.RemoveOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 128 {
		t.Fatalf("bad: %v", k)
	}

	k, _, ok = l.RemoveOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 129 {
		t.Fatalf("bad: %v", k)
	}
}

func TestLRU_PeekAndDeleteOldest(t *testing.T) {
	l := sync_.NewLRU[int, int](128)

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	k, _, ok := l.PeekOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 128 {
		t.Fatalf("bad: %v", k)
	}

	k, v, ok := l.PeekAndDeleteOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 128 {
		t.Fatalf("bad key: %v", k)
	}
	if v != 128 {
		t.Fatalf("bad value: %v", k)
	}

	k, _, ok = l.PeekAndDeleteOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 129 {
		t.Fatalf("bad: %v", k)
	}
}

func TestLRU_RemoveOldest(t *testing.T) {
	l := sync_.NewLRU[int, int](128)

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	k, _, ok := l.GetOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 128 {
		t.Fatalf("bad: %v", k)
	}

	k, _, ok = l.RemoveOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 128 {
		t.Fatalf("bad: %v", k)
	}

	k, _, ok = l.RemoveOldest()
	if !ok {
		t.Fatalf("missing")
	}
	if k != 129 {
		t.Fatalf("bad: %v", k)
	}
}

func TestConcurrentRange(t *testing.T) {
	const lruSize = 1 << 10

	m := sync_.NewLRU[int64, int64](lruSize)
	for n := int64(1); n <= lruSize; n++ {
		m.Store(n, n)
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	defer func() {
		close(done)
		wg.Wait()
	}()
	for g := int64(runtime.GOMAXPROCS(0)); g > 0; g-- {
		r := rand.New(rand.NewSource(g))
		wg.Add(1)
		go func(g int64) {
			defer wg.Done()
			for i := int64(0); ; i++ {
				select {
				case <-done:
					return
				default:
				}
				for n := int64(1); n < lruSize; n++ {
					if r.Int63n(lruSize) == 0 {
						m.Store(n, n*i*g)
					} else {
						m.Load(n)
					}
				}
			}
		}(g)
	}

	iters := 1 << 10
	if testing.Short() {
		iters = 16
	}
	for n := iters; n > 0; n-- {
		seen := make(map[int64]bool, lruSize)

		m.Range(func(ki, vi int64) bool {
			k, v := ki, vi
			if v%k != 0 {
				t.Fatalf("while Storing multiples of %v, Range saw value %v", k, v)
			}
			if seen[k] {
				t.Fatalf("Range visited key %v twice", k)
			}
			seen[k] = true
			return true
		})

		if len(seen) != lruSize {
			t.Fatalf("Range visited %v elements of %v-element Map", len(seen), lruSize)
		}
	}
}
