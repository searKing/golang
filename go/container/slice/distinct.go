// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"github.com/searKing/golang/go/util/object"
)

// DistinctFunc returns a slice consisting of the distinct elements (according to
// {@link Object#equals(Object)}) of this slice.
func DistinctFunc(s interface{}, f func(interface{}, interface{}) int) interface{} {
	return normalizeSlice(distinctFunc(Of(s), f), s)
}

// distinctFunc is the same as DistinctFunc except that if
//// truth==false, the sense of the predicate function is
//// inverted.
func distinctFunc(s []interface{}, f func(interface{}, interface{}) int) []interface{} {
	object.RequireNonNil(s, "distinctFunc called on nil slice")
	object.RequireNonNil(f, "distinctFunc called on nil callfn")

	sDistinctMap := map[interface{}]struct{}{}
	var sDistincted = []interface{}{}
	for _, r := range s {
		if _, ok := sDistinctMap[r]; ok {
			continue
		}
		sDistinctMap[r] = struct{}{}
		sDistincted = append(sDistincted, r)
	}
	return sDistincted
}
