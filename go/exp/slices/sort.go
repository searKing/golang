// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"golang.org/x/exp/constraints"
)

// LinearSearch searches for target in a sorted slice and returns the position
// where target is found, or the position where target would appear in the
// sort order; it also returns a bool saying whether the target is really found
// in the slice. The slice must be sorted in increasing order.
// Note: Binary-search was compared using the benchmarks. The following
// code is equivalent to the linear search above:
//
//     pos := sort.Search(len(x), func(i int) bool {
//         return target < x[i]
//     })
//
// The binary search wins for very large boundary sets, but
// the linear search performs better up through arrays between
// 256 and 512 elements, so we continue to prefer linear search.
func LinearSearch[E constraints.Ordered](x []E, target E) (int, bool) {
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
func LinearSearchFunc[E any](x []E, target E, cmp func(E, E) int) (int, bool) {
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
