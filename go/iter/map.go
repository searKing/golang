// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter

import (
	"fmt"
	"iter"
)

// Map returns an iterator that yields the individual values mapped by format "%v" within all v in the sequences.
// TODO: accept [V any, M ~string] if go support template type deduction
func Map[V any, M string](seq iter.Seq[V]) iter.Seq[M] {
	return MapFunc(seq, func(v V) M { return M(fmt.Sprintf("%v", v)) })
}

// MapFunc returns an iterator that yields the individual values mapped by f(c) within all v in the sequences.
func MapFunc[V, M any](seq iter.Seq[V], f func(V) M) iter.Seq[M] {
	return func(yield func(M) bool) {
		for v := range seq {
			if !yield(f(v)) {
				break
			}
		}
	}
}
