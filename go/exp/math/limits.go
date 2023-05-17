// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"unsafe"

	"golang.org/x/exp/constraints"
)

// MinInt returns the smallest value representable by the type.
func MinInt[T constraints.Integer]() T {
	var zero T
	minusOne := ^zero
	if minusOne > 0 {
		return zero // Unsigned
	}
	bits := unsafe.Sizeof(zero) << 3
	return minusOne << (bits - 1) // Signed
}

// MaxInt returns the largest value representable by the type.
func MaxInt[T constraints.Integer]() T {
	return ^MinInt[T]()
}
