// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"github.com/searKing/golang/go/util/object"
)

// MapFunc returns a slice consisting of the results of applying the given
// function to the elements of this slice.
func MapFunc(s interface{}, f func(interface{}) interface{}) interface{} {
	return normalizeSlice(mapFunc(Of(s), f), s)
}

// mapFunc is the same as MapFunc
func mapFunc(s []interface{}, f func(interface{}) interface{}) []interface{} {
	object.RequireNonNil(s, "mapFunc called on nil slice")
	object.RequireNonNil(f, "mapFunc called on nil callfn")

	var sMapped = []interface{}{}
	for _, r := range s {
		sMapped = append(sMapped, f(r))
	}
	return sMapped
}
