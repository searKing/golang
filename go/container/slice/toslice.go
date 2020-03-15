// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

// ToSliceFunc returns an array containing the elements of this stream.
func ToSliceFunc(s interface{}) interface{} {
	return toSliceFunc(Of(s))
}

// toSliceFunc is the same as ToSliceFunc
func toSliceFunc(s []interface{}) []interface{} {
	return s
}
