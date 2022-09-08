// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
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
func Max[T constraints.Ordered](s ...T) T {
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
func Min[T constraints.Ordered](s ...T) T {
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
func Clamp[T constraints.Ordered](v, lo, hi T) T {
	if lo > hi {
		t := lo
		lo = hi
		hi = t
	}
	if v < lo {
		return lo
	}
	if hi < v {
		return hi
	}
	return v
}
