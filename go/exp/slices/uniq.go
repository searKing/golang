// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// Uniq returns a slice satisfying c != zero within all c in the slice.
// Uniq modifies the contents of the slice s; it does not create a new slice.
func Uniq[S ~[]E, E comparable](s S) S {
	if len(s) == 0 {
		return s
	}
	m := make(map[E]bool)
	for _, v := range s {
		m[v] = true
	}

	i := 1
	m[s[0]] = false
	for _, v := range s[1:] {
		save, has := m[v]
		if has && save {
			s[i] = v
			m[v] = false
			i++
		}
	}
	return s[:i]
}

// UniqFunc returns a slice satisfying f(c) within all c in the slice.
func UniqFunc[S ~[]E, E any](s S, f func(v1, v2 E) bool) S {
	if len(s) == 0 {
		return s
	}

	i := 1
	for _, v := range s[1:] {
		if ContainsFunc(s[:i], func(e E) bool { return f(v, e) }) {
			continue
		}
		s[i] = v
		i++
	}
	return s[:i]
}
