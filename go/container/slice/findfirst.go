// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

// FindFirstFunc returns an {@link Optional} describing the first element of this stream,
// or an empty {@code Optional} if the stream is empty.  If the stream has
// no encounter order, then any element may be returned.
func FindFirstFunc(s interface{}, f func(interface{}) bool) interface{} {
	return normalizeElem(findFirstFunc(Of(s), f, true), s)
}

// findFirstFunc is the same as FindFirstFunc.
func findFirstFunc(s []interface{}, f func(interface{}) bool, truth bool) interface{} {
	idx := findFirstIndexFunc(s, f, truth)
	if idx == -1 {
		return nil
	}
	return s[idx]
}
