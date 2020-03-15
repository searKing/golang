// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"github.com/searKing/golang/go/util/object"
	"github.com/searKing/golang/go/util/optional"
)

// ReduceFunc calls a defined callback function on each element of an array, and returns an array that contains the results.
func ReduceFunc(s interface{}, f func(left, right interface{}) interface{}) interface{} {
	return normalizeElem(reduceFunc(Of(s), f), s)
}

// reduceFunc is the same as ReduceFunc
func reduceFunc(s []interface{}, f func(left, right interface{}) interface{}, identity ...interface{}) interface{} {
	object.RequireNonNil(s, "reduceFunc called on nil slice")
	object.RequireNonNil(f, "reduceFunc called on nil callfn")

	var foundAny bool
	var result interface{}

	if identity != nil || len(identity) != 0 {
		foundAny = true
		result = identity
	}
	for _, r := range s {
		if !foundAny {
			foundAny = true
			result = r
		} else {
			result = f(result, r)
		}
	}
	if foundAny {
		return optional.Of(result).Get()
	}
	return optional.Empty().Get()
}
