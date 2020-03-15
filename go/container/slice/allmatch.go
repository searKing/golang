// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

// AllMatchFunc returns whether all elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for
// determining the result.  If the stream is empty then {@code true} is
// returned and the predicate is not evaluated.
func AllMatchFunc(s interface{}, f func(interface{}) bool) bool {
	return allMatchFunc(Of(s), f, true)
}

// allMatchFunc is the same as AllMatchFunc.
func allMatchFunc(s []interface{}, f func(interface{}) bool, truth bool) bool {
	return !anyMatchFunc(s, f, !truth)
}
