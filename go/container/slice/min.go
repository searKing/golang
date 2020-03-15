// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"github.com/searKing/golang/go/util/object"
)

// MinFunc returns the minimum element of this stream according to the provided.
func MinFunc(s interface{}, f func(interface{}, interface{}) int) interface{} {
	return normalizeElem(minFunc(Of(s), f), s)

}

// minFunc is the same as MinFunc
func minFunc(s []interface{}, f func(interface{}, interface{}) int) interface{} {
	object.RequireNonNil(s, "minFunc called on nil slice")
	object.RequireNonNil(f, "minFunc called on nil callfn")

	return ReduceFunc(s, func(left, right interface{}) interface{} {
		if f(left, right) < 0 {
			return left
		}
		return right
	})
}
