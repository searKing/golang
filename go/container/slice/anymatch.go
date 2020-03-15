// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"github.com/searKing/golang/go/util/object"
)

// AnyMatchFunc returns whether any elements of this stream match the provided
// predicate.  May not evaluate the predicate on all elements if not
// necessary for determining the result.  If the stream is empty then
// {@code false} is returned and the predicate is not evaluated.
func AnyMatchFunc(s interface{}, f func(interface{}) bool) bool {
	return anyMatchFunc(Of(s), f, true)
}

// anyMatchFunc is the same as AnyMatchFunc.
func anyMatchFunc(s []interface{}, f func(interface{}) bool, truth bool) bool {
	object.RequireNonNil(s, "anyMatchFunc called on nil slice")
	object.RequireNonNil(f, "anyMatchFunc called on nil callfn")

	for _, r := range s {
		if f(r) == truth {
			return true
		}
	}
	return false
}
