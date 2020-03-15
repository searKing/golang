// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"github.com/searKing/golang/go/util/object"
)

// ForEachOrderedFunc Performs an action for each element of this slice.
// <p>This operation processes the elements one at a time, in encounter
// order if one exists.  Performing the action for one element
// performing the action for subsequent elements, but for any given element,
// the action may be performed in whatever thread the library chooses.
func ForEachOrderedFunc(s interface{}, f func(interface{})) {
	forEachOrderedFunc(Of(s), f)
}

// forEachOrderedFunc is the same as ForEachOrderedFunc
func forEachOrderedFunc(s []interface{}, f func(interface{})) {
	object.RequireNonNil(s, "forEachOrderedFunc called on nil slice")
	object.RequireNonNil(f, "forEachOrderedFunc called on nil callfn")

	for _, r := range s {
		f(r)
	}
	return
}
