// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

// DropUntilFunc returns, if this slice is ordered, a slice consisting of the remaining
// elements of this slice after dropping the longest prefix of elements
// that match the given predicate.  Otherwise returns, if this slice is
// unordered, a slice consisting of the remaining elements of this slice
// after dropping a subset of elements that match the given predicate.
func DropUntilFunc(s interface{}, f func(interface{}) bool) interface{} {
	return normalizeSlice(dropUntilFunc(Of(s), f, true), s)
}

// dropUntilFunc is the same as DropUntilFunc.
func dropUntilFunc(s []interface{}, f func(interface{}) bool, truth bool) []interface{} {
	return dropWhileFunc(s, f, !truth)
}
