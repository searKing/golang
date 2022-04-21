// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import (
	"golang.org/x/exp/slices"
)

// FirstFunc returns the first item from the slice satisfying f(c), or zero if none do.
func FirstFunc[S ~[]E, E any](s S, f func(E) bool) (e E, ok bool) {
	var zeroE E
	i := slices.IndexFunc(s, f)
	if i == -1 {
		return zeroE, false
	}
	return s[i], true
}

// FirstOrZero returns the first non-zero item from the slice, or zero if none do.
func FirstOrZero[E comparable](s ...E) E {
	var zeroE E
	return FirstOrZeroFunc(s, func(e E) bool { return e != zeroE })
}

// FirstOrZeroFunc returns the first item from the slice satisfying f(c), or zero if none do.
func FirstOrZeroFunc[S ~[]E, E any](s S, f func(E) bool) E {
	e, _ := FirstFunc(s, f)
	return e
}
