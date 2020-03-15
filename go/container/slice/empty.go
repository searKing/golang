// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slice

// EmptyFunc returns an empty sequential {@code slice}.
func EmptyFunc(s interface{}) interface{} {
	return normalizeSlice(emptyFunc(), s)
}

// emptyFunc is the same as EmptyFunc
func emptyFunc() []interface{} {
	return []interface{}{}
}
