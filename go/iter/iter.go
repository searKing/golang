// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter

import (
	"iter"
)

// Filter returns an iterator that yields the individual values satisfying v != zero within all v in the sequences.
func Filter[V comparable](seq iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			var zeroV V
			if v != zeroV && !yield(v) {
				break
			}
		}
	}
}

// FilterFunc returns an iterator that yields the individual values satisfying f(v) within all v in the sequences.
func FilterFunc[V any](seq iter.Seq[V], f func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if f(v) && !yield(v) {
				break
			}
		}
	}
}

// FilterN returns an iterator that yields the individual values at most in the sequences.
//
// The count determines the number of individual values to return:
//   - n > 0: at most n individual values;
//   - n == 0: zero individual values;
//   - n < 0: all individual values.
func FilterN[V any](seq iter.Seq[V], n int) iter.Seq[V] {
	if n < 0 {
		return seq
	}
	if n == 0 {
		return func(yield func(V) bool) { return }
	}
	return func(yield func(V) bool) {
		var i int
		for v := range seq {
			if i >= n {
				break
			}
			i++

			if !yield(v) {
				break
			}
		}
	}
}
