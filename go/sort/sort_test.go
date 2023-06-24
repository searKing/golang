// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sort_test

import (
	"sort"
	"testing"

	sort_ "github.com/searKing/golang/go/sort"
)

var ints = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}

// Convenience types for common cases

// IntSlice attaches the methods of Interface to []int, sorting in increasing order.
type IntSlice []int

func (x IntSlice) Len() int           { return len(x) }
func (x IntSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x IntSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x IntSlice) Sort() { sort.Sort(x) }

func TestPartialSortIntSlice(t *testing.T) {
	data := ints
	data1 := ints
	k := 3
	a := IntSlice(data[:])
	sort.Sort(a[:k])
	if sort_.IsPartialSorted(a, k) {
		t.Errorf("partial sort did sort")
	}
	r := IntSlice(data1[:])
	sort.Sort(r)
	if !sort_.IsPartialSorted(r, k) {
		t.Errorf("partial sort didn't sort")
	}
}
