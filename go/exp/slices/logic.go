// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// OrFunc tests the slice satisfying f(c), or false if none do.
// return true if len(s) == 0
func OrFunc[S ~[]E, E any](s S, f func(E) bool) bool {
	if len(s) == 0 {
		return true
	}
	for _, e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

// AndFunc tests the slice satisfying f(c), or true if none do.
// return true if len(s) == 0
func AndFunc[S ~[]E, E any](s S, f func(E) bool) bool {
	for _, e := range s {
		if !f(e) {
			return false
		}
	}
	return true
}

// Or tests whether the slice satisfying c != zero within any c in the slice.
// return true if len(s) == 0
func Or[E comparable](s ...E) bool {
	if len(s) == 0 {
		return true
	}
	var zeroE E
	for _, e := range s {
		if e != zeroE {
			return true
		}
	}
	return false
}

// And tests whether the slice satisfying c != zero within all c in the slice.
// return true if len(s) == 0
func And[E comparable](s ...E) bool {
	var zeroE E
	for _, e := range s {
		if e == zeroE {
			return false
		}
	}
	return true
}
