// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strings

// JoinGenerator supplies sep between strings step by step, with mapping if consists
// [r0,r1,r2] -> "r0'""sep""r1'""sep""r2'"
func JoinGenerator(sep string, mapping func(s string) string) func(r string) string {
	var written bool
	return func(s string) string {
		if mapping != nil {
			s = mapping(s)
		}
		if written {
			s = sep + s
		}
		written = true
		return s
	}
}
