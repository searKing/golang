// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps

import (
	"fmt"
)

// Map returns a map mapped by format "%v" within all kv pairs in the map.
// Map does not modify the contents of the map m; it creates a new map.
// TODO: accept [M ~map[K]V, K comparable, V any, R ~map[KR]KV, KR comparable, KV any] if go support template type deduction
func Map[M ~map[K]V, K comparable, V any, R map[KR]KV, KR string, KV string](m M) R {
	return MapFunc(m, func(k K, v V) (KR, KV) {
		return KR(fmt.Sprintf("%v", k)), KV(fmt.Sprintf("%v", v))
	})
}

// MapFunc returns a map mapped by f(k,v) within all kv pairs in the map.
// MapFunc does not modify the contents of the map m; it creates a new map.
// TODO: accept [M ~map[K]V, K comparable, V any, R ~map[KR]KV, KR comparable, KV any] if go support template type deduction
func MapFunc[M ~map[K]V, K comparable, V any, R map[KR]KV, KR comparable, KV any](m M, f func(K, V) (KR, KV)) R {
	if m == nil {
		return nil
	}

	var rr = make(R, len(m))
	for k, v := range m {
		kr, kv := f(k, v)
		rr[kr] = kv
	}
	return rr
}
