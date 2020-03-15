// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"github.com/searKing/golang/go/util/object"
)

// TakeWhileFunc returns, if this slice is ordered, a slice consisting of the longest
// prefix of elements taken from this slice that match the given predicate.
// Otherwise returns, if this slice is unordered, a slice consisting of a
// subset of elements taken from this slice that match the given predicate.
func TakeWhileFunc(s interface{}, f func(interface{}) bool) interface{} {
	return normalizeSlice(takeWhileFunc(Of(s), f, true), s)
}

// takeWhileFunc is the same as TakeWhileFunc.
func takeWhileFunc(s []interface{}, f func(interface{}) bool, truth bool) []interface{} {
	object.RequireNonNil(s, "takeWhileFunc called on nil slice")
	object.RequireNonNil(f, "takeWhileFunc called on nil callfn")

	var sTaken = []interface{}{}
	for _, r := range s {
		if f(r) == truth {
			sTaken = append(sTaken, r)
			continue
		}
		break
	}
	return sTaken
}
