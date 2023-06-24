// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sort

import (
	"sort"

	"github.com/searKing/golang/go/exp/math"
)

type partial struct {
	// This embedded Interface permits Partial to use the methods of
	// another Interface implementation.
	sort.Interface

	Size int
}

// Len returns the opposite of the embedded implementation's Len method.
func (r partial) Len() int {
	return math.Min(r.Size, r.Interface.Len())
}

// IsPartialSorted reports whether data[:k] is partial sorted, as top k of data[:].
func IsPartialSorted(data sort.Interface, k int) bool {
	if k >= data.Len() {
		return sort.IsSorted(data)
	}

	// data[:k] sorted
	if !sort.IsSorted(&partial{data, k}) {
		return false
	}
	// data[k-1] < any element in data[k:]
	n := data.Len()
	for i := k; i < n; i++ {
		if data.Less(i, k-1) {
			return false
		}
	}
	return true
}
