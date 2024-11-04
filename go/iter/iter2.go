// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter

import (
	"iter"
)

// Filter2 returns an iterator that yields the pairs of values satisfying v != zero within all [k,v] in the sequences.
func Filter2[K, V comparable](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			var zeroV V
			if v != zeroV && !yield(k, v) {
				break
			}
		}
	}
}

// Filter2Func returns an iterator that yields the pairs of values satisfying f(v) within all [k,v] in the sequences.
func Filter2Func[K, V any](seq iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if f(k, v) && !yield(k, v) {
				break
			}
		}
	}
}

// Filter2N returns an iterator that yields the pairs of values at most in the sequences.
//
// The count determines the number of pairs of values to return:
//   - n > 0: at most n pairs of values;
//   - n == 0: zero pairs of values;
//   - n < 0: all pairs of values.
func Filter2N[K, V any](seq iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	if n < 0 {
		return seq
	}
	if n == 0 {
		return func(yield func(K, V) bool) { return }
	}
	return func(yield func(K, V) bool) {
		var i int
		for k, v := range seq {
			if i >= n {
				break
			}
			i++

			if !yield(k, v) {
				break
			}
		}
	}
}
