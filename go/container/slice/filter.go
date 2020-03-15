// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"github.com/searKing/golang/go/util/object"
)

// FilterFunc returns a slice consisting of the elements of this slice that match
// the given predicate.
func FilterFunc(s interface{}, f func(interface{}) bool) interface{} {
	return normalizeSlice(filterFunc(Of(s), f, true), s)
}

// filterFunc is the same as FilterFunc except that if
// truth==false, the sense of the predicate function is
// inverted.
func filterFunc(s []interface{}, f func(interface{}) bool, truth bool) []interface{} {
	object.RequireNonNil(s, "filterFunc called on nil slice")
	object.RequireNonNil(f, "filterFunc called on nil callfn")

	var sFiltered = []interface{}{}
	for _, r := range s {
		if f(r) == truth {
			sFiltered = append(sFiltered, r)
		}
	}
	return sFiltered
}
