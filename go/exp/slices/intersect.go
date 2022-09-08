// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// Intersect returns a slice satisfying c != zero within all c in the slice.
// Intersect does not modify the contents of the slice s1 and s2; it creates a new slice.
func Intersect[S ~[]E, E comparable](s1, s2 S) S {
	if len(s1) == 0 {
		return s2
	}
	if len(s2) == 0 {
		return s1
	}
	m := make(map[E]struct{})
	for _, v := range s1 {
		m[v] = struct{}{}
	}

	var ss S
	for _, v := range s2 {
		if len(m) == 0 {
			break
		}
		if _, ok := m[v]; ok {
			ss = append(ss, v)
			delete(m, v)
		}
	}
	return ss
}

// IntersectFunc returns a slice satisfying f(c) within all c in the slice.
// IntersectFunc does not modify the contents of the slice s1 and s2; it creates a new slice.
func IntersectFunc[S ~[]E, E any](s1, s2 S, f func(v1, v2 E) bool) S {
	if len(s1) == 0 {
		return s2
	}
	if len(s2) == 0 {
		return s1
	}

	var ss S
	for _, v := range s1 {
		if ContainsFunc(ss, func(e E) bool { return f(v, e) }) {
			continue
		}
		if ContainsFunc(s2, func(e E) bool { return f(v, e) }) {
			ss = append(ss, v)
		}
	}
	return ss
}
