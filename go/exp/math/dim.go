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
