// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// Group returns a map group by all elements in s that have the same values.
//
// If s is nil, Group returns nil (zero map).
//
// If s is empty, Group returns an empty map.
//
// Else, Group returns a map group by all elements in s that have the same values.
// TODO: accept [S ~[]E, M ~map[E]N, E comparable, N constraints.Number] if go support template type deduction
func Group[S ~[]E, M map[E]N, E comparable, N int](s S) M {
	// If s is nil, Split returns nil (zero map).
	if s == nil {
		return nil
	}

	// Below: s != nil

	// If s is empty, Split returns an empty map.
	if len(s) == 0 {
		var emptyM = M{}
		return emptyM
	}

	// Below: len(s) > 0 && f != nil
	var m = M{}
	for _, e := range s {
		m[e] = m[e] + 1
	}
	return m
}

// GroupFunc returns a map satisfying f(c) within
// all c in the map group by all elements in s that have the same values.
//
// If s is nil, GroupFunc returns nil (zero map).
//
// If s and f are both empty or nil, GroupFunc returns an empty map.
//
// Else, GroupFunc returns a map satisfying f(c)
// within all c in the map group by all elements in s that have the same values.
// TODO: accept [S ~[]E, M ~map[K]S, E any, K comparable] if go support template type deduction
func GroupFunc[S ~[]E, M map[K]S, E any, K comparable](s S, f func(E) K) M {
	// If s is nil, Split returns nil (zero submaps).
	if s == nil {
		return nil
	}

	// Below: s != nil

	// If both s and f are empty or nil, Split returns an empty map.
	if len(s) == 0 && f == nil {
		var emptyM = M{}
		return emptyM
	}

	// Below: len(s) > 0 && f != nil
	var m = M{}
	for _, e := range s {
		k := f(e)
		m[k] = append(m[k], e)
	}
	return m
}
