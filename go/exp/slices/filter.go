// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// Filter returns a slice satisfying c != zero within all c in the slice.
// Filter modifies the contents of the slice s; it does not create a new slice.
func Filter[S ~[]E, E comparable](s S) S {
	if len(s) == 0 {
		return s
	}
	i := 0
	for _, v := range s {
		var zeroE E
		if v != zeroE {
			s[i] = v
			i++
		}
	}
	return s[:i]
}

// FilterFunc returns a slice satisfying f(c) within all c in the slice.
// FilterFunc modifies the contents of the slice s; it does not create a new slice.
func FilterFunc[S ~[]E, E any](s S, f func(E) bool) S {
	if len(s) == 0 {
		return s
	}
	i := 0
	for _, v := range s {
		if f(v) {
			s[i] = v
			i++
		}
	}
	return s[:i]
}
