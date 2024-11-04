// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter

import (
	"fmt"
	"iter"
)

// Map2 returns an iterator that yields the pairs of values mapped by format "%v" within all [k,v] in the sequences.
// TODO: accept [K, V any, KR, KV ~string] if go support template type deduction
func Map2[K, V any, KR, KV string](seq iter.Seq2[K, V]) iter.Seq2[KR, KV] {
	return Map2Func(seq, func(k K, v V) (KR, KV) {
		return KR(fmt.Sprintf("%v", k)), KV(fmt.Sprintf("%v", v))
	})
}

// Map2Func returns an iterator that yields the pairs of values mapped by f(k,v) within all [k,v] in the sequences.
func Map2Func[K, V, KR, KV any](seq iter.Seq2[K, V], f func(K, V) (KR, KV)) iter.Seq2[KR, KV] {
	return func(yield func(KR, KV) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
				break
			}
		}
	}
}
