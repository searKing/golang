// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// Diff returns two slices:
//
//	one containing elements that are present in the new slice but not in the old slice (add),
//	and the other containing elements that are present in the old slice but not in the new slice (del).
//
// Diff does not modify the original slices, but creates and returns two new slices
//
//	one for added elements
//	one for deleted elements.
func Diff[T comparable](new, old []T) (add, del []T) {
	oldSet := make(map[T]struct{}, len(old))
	for _, v := range old {
		oldSet[v] = struct{}{}
	}
	newSet := make(map[T]struct{}, len(new))
	for _, v := range new {
		newSet[v] = struct{}{}
	}

	// Pre-allocate slices to avoid multiple allocations
	add = make([]T, 0, len(new))
	del = make([]T, 0, len(old))
	// Find elements that are in old but not in new (to be deleted)
	for _, v := range old {
		if _, ok := newSet[v]; !ok {
			del = append(del, v)
		}
	}
	// Find elements that are in new but not in old (to be added)
	for _, v := range new {
		if _, ok := oldSet[v]; !ok {
			add = append(add, v)
		}
	}
	return add, del
}
