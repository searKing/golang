// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps

// SliceFunc returns a slice mapped by f(k,v) within all (k,v) in the map.
// SliceFunc does not modify the contents of the map m; it creates a new slice.
// TODO: accept [M ~map[K]V, K comparable, S ~[]E, V any, E any] if go support template type deduction
func SliceFunc[M ~map[K]V, K comparable, S []E, V any, E any](m M, f func(k K, v V) E) S {
	if m == nil {
		return nil
	}
	var s = make(S, len(m))
	var i int
	for k, v := range m {
		s[i] = f(k, v)
		i++
	}
	return s
}
