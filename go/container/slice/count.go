// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

import (
	"github.com/searKing/golang/go/util/object"
)

// CountFunc returns the maximum element of this stream according to the provided.
func CountFunc(s interface{}) int {
	return countFunc(Of(s))

}

// countFunc is the same as CountFunc
func countFunc(s []interface{}) int {
	object.RequireNonNil(s, "countFunc called on nil slice")
	return len(s)
}
