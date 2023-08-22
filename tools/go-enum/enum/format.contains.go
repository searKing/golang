// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enum

// Arguments to format are:
//
//	[1]: type name
const containsTemplate = `
// %[1]sSliceContains reports whether sunEnums is within enums.
func %[1]sSliceContains(enums []%[1]s, sunEnums ...%[1]s) bool {
	var seenEnums = map[%[1]s]bool{}
	for _, e := range sunEnums {
		seenEnums[e] = false
	}

	for _, v := range enums {
		if _, has := seenEnums[v]; has {
			seenEnums[v] = true
		}
	}

	for _, seen := range seenEnums {
		if !seen {
			return false
		}
	}

	return true
}

// %[1]sSliceContainsAny reports whether any sunEnum is within enums.
func %[1]sSliceContainsAny(enums []%[1]s, sunEnums ...%[1]s) bool {
	var seenEnums = map[%[1]s]struct{}{}
	for _, e := range sunEnums {
		seenEnums[e] = struct{}{}
	}

	for _, v := range enums {
		if _, has := seenEnums[v]; has {
			return true
		}
	}

	return false
}
`
