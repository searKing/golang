// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"container/heap"
	"sort"

	sort_ "github.com/searKing/golang/go/sort"
	"golang.org/x/exp/constraints"
)

// LinearSearch searches for target in a sorted slice and returns the position
// where target is found, or the position where target would appear in the
// sort order; it also returns a bool saying whether the target is really found
// in the slice. The slice must be sorted in increasing order.
// Note: Binary-search was compared using the benchmarks. The following
// code is equivalent to the linear search above:
//
//	pos := sort.Search(len(x), func(i int) bool {
//	    return target < x[i]
//	})
//
// The binary search wins for very large boundary sets, but
// the linear search performs better up through arrays between
// 256 and 512 elements, so we continue to prefer linear search.
func LinearSearch[S ~[]E, E constraints.Ordered](x S, target E) (int, bool) {
	// search returns the leftmost position where f returns true, or len(x) if f
	// returns false for all x. This is the insertion position for target in x,
	// and could point to an element that's either == target or not.
	pos := search(len(x), func(i int) bool {
		return x[i] >= target
	})
	if pos >= len(x) || x[pos] != target {
		return pos, false
	} else {
		return pos, true
	}
}

// LinearSearchFunc works like LinearSearch, but uses a custom comparison
// function. The slice must be sorted in increasing order, where "increasing" is
// defined by cmp. cmp(a, b) is expected to return an integer comparing the two
// parameters: 0 if a == b, a negative number if a < b and a positive number if
// a > b.
func LinearSearchFunc[S ~[]E, E any](x S, target E, cmp func(E, E) int) (int, bool) {
	pos := search(len(x), func(i int) bool { return cmp(x[i], target) >= 0 })
	if pos >= len(x) || cmp(x[pos], target) != 0 {
		return pos, false
	} else {
		return pos, true
	}
}

func search(n int, f func(int) bool) int {
	pos := n
	for i := 0; i < n; i++ {
		if f(i) {
			return i // preserves f(i) == true
		}
	}
	return pos
}

// PartialSort rearranges elements such that the range [0, m)
// contains the sorted m smallest elements in the range [first, data.Len).
// The order of equal elements is not guaranteed to be preserved.
// The order of the remaining elements in the range [m, data.Len) is unspecified.
// PartialSort modifies the contents of the slice s; it does not create a new slice.
func PartialSort[S ~[]E, E constraints.Ordered](s S, k int) {
	if s == nil {
		return
	}

	if k <= 0 {
		return
	}
	if k >= len(s) {
		sort.Sort(SortSlice[E](s))
		return
	}

	var ss = make(S, k)
	copy(ss, s[:k])
	h := MaxHeap[E](ss)
	heap.Init(&h)
	for i, v := range s[k:] {
		heap.Push(&h, v)
		vv := heap.Pop(&h).(E)
		s[i+k] = vv
	}

	for h.Len() > 0 {
		s[h.Len()-1] = heap.Pop(&h).(E)
	}
	return
}

func PartialSortFunc[S ~[]E, E any](s S, k int, cmp func(E, E) int) {
	if s == nil {
		return
	}

	if k <= 0 {
		return
	}

	if k >= len(s) {
		k = len(s)
	}

	var ss = make(S, k)
	copy(ss, s[:k])

	if k <= 0 {
		return
	}
	// MaxHeap
	h := NewHeapFunc(ss, func(a E, b E) int {
		if cmp == nil {
			return 0
		}
		return -cmp(a, b)
	})
	heap.Init(h)
	for i, v := range s[k:] {
		heap.Push(h, v)
		vv := heap.Pop(h).(E)
		s[i+k] = vv
	}

	for h.Len() > 0 {
		s[h.Len()-1] = heap.Pop(h).(E)
	}
	return
}

// IsPartialSorted reports whether data[:k] is partial sorted, as top k of data[:].
func IsPartialSorted[S ~[]E, E constraints.Ordered](s S, k int) bool {
	return sort_.IsPartialSorted(SortSlice[E](s), k)
}

// Convenience types for common cases

// SortSlice attaches the methods of Interface to []E, sorting in increasing order.
type SortSlice[E constraints.Ordered] []E

func (x SortSlice[E]) Len() int { return len(x) }

// Less reports whether x[i] should be ordered before x[j], as required by the sort Interface.
// Note that floating-point comparison by itself is not a transitive relation: it does not
// report a consistent ordering for not-a-number (NaN) values.
// This implementation of Less places NaN values before any others, by using:
//
//	x[i] < x[j] || (math.IsNaN(x[i]) && !math.IsNaN(x[j]))
func (x SortSlice[E]) Less(i, j int) bool { return x[i] < x[j] || (isNaN(x[i]) && !isNaN(x[j])) }
func (x SortSlice[E]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x SortSlice[E]) Sort() { sort.Sort(x) }

// isNaN is a copy of math.IsNaN to avoid a dependency on the math package.
func isNaN[E constraints.Ordered](f E) bool {
	return f != f
}
