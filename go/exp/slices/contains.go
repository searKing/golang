// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import "slices"

// Contains reports whether v is present in s.
// Deprecated: Use slices.Contains instead since go1.21.
func Contains[E comparable](s []E, v E) bool {
	return slices.Contains(s, v)
}

// ContainsFunc reports whether v satisfying f(s[i]) is present in s.
// Deprecated: Use slices.ContainsFunc instead since go1.21.
func ContainsFunc[E any](s []E, f func(E) bool) bool {
	return slices.IndexFunc(s, f) >= 0
}
