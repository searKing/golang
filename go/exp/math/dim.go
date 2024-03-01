// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"cmp"
	"math"

	constraints_ "github.com/searKing/golang/go/exp/constraints"
	"golang.org/x/exp/constraints"
)

// Dim returns the maximum of x-y or 0.
func Dim[T constraints_.Number](x, y T) T {
	v := x - y
	if v <= 0 {
		// v is negative or 0
		return 0
	}
	// v is positive or NaN
	return v
}

// Max returns the largest of s.
// Deprecated: Use [max] built-in instead since go1.21.
func Max[T cmp.Ordered](s ...T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s[1:] {
		if m < v {
			m = v
		}
	}
	return m
}

// Min returns the smallest of s.
// Deprecated: Use [min] built-in instead since go1.21.
func Min[T cmp.Ordered](s ...T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	m := s[0]
	for _, v := range s[1:] {
		if m > v {
			m = v
		}
	}
	return m
}

// Clamp returns the value between boundary [lo,hi], as v < lo ? v : hi > v : hi : v.
// Reference to lo if v is less than lo, reference to hi if hi is less than v, otherwise reference to v.
// If v compares equivalent to either bound, returns a reference to v, not the bound.
func Clamp[T cmp.Ordered](v, lo, hi T) T {
	if lo > hi {
		lo, hi = hi, lo
	}
	if v < lo {
		return lo
	}
	if hi < v {
		return hi
	}
	return v
}

// Sum returns the sum of s.
func Sum[T, R constraints_.Number](s ...T) R {
	if len(s) == 0 {
		var zero R
		return zero
	}
	m := R(s[0])
	for _, v := range s[1:] {
		m += R(v)
	}
	return m
}

// Mean returns the mean of s.
func Mean[T constraints_.Number, R constraints.Float](s ...T) R {
	if len(s) == 0 {
		var zero R
		return zero
	}
	return Sum[T, R](s...) / R(len(s))
}

// Variance returns the variance of s.
func Variance[T constraints_.Number, R constraints.Float](s ...T) R {
	if len(s) == 0 || len(s) == 1 {
		var zero R
		return zero
	}
	m := Mean[T, R](s...)

	var res R
	for _, v := range s {
		d := R(v) - m
		res += d * d
	}

	return res / R(len(s)-1)
}

// StandardDeviation returns the standard deviation  of s.
func StandardDeviation[T constraints_.Number, R constraints.Float](s ...T) R {
	if len(s) == 0 || len(s) == 1 {
		var zero R
		return zero
	}
	return R(math.Sqrt(Variance[T, float64](s...)))
}
