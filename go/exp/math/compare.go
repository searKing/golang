// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"cmp"
)

// Compare returns an integer comparing two elements.
// The result will be 0 if a == b, -1 if a < b, and +1 if a > b.
// This implementation places NaN values before any others, by using:
//
//	a < b || (math.IsNaN(a) && !math.IsNaN(b))
//
// Deprecated: Use [cmp.Compare] instead since go1.21.
func Compare[E cmp.Ordered](a E, b E) int {
	if a < b || IsNaN(a) && !IsNaN(b) {
		return -1
	}

	if a == b || IsNaN(a) && IsNaN(b) {
		return 0
	}
	return 1
}

// Reverse returns the reverse comparison for cmp, as cmp(b, a).
func Reverse[E cmp.Ordered](cmp func(a E, b E) int) func(a E, b E) int {
	return func(a E, b E) int {
		return cmp(b, a)
	}
}

// IsNaN is a copy of math.IsNaN to avoid a dependency on the math package.
func IsNaN[E cmp.Ordered](f E) bool {
	return f != f
}
