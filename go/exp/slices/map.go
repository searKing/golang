// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

import "fmt"

// Map returns a slice mapped by format "%v" within all c in the slice.
// Map does not modify the contents of the slice s; it creates a new slice.
// TODO: accept [S ~[]E, E any, R ~[]M, M ~string] if go support template type deduction
func Map[S ~[]E, E any, R []M, M string](s S) R {
	return MapFunc(s, func(e E) M {
		return M(fmt.Sprintf("%v", e))
	})
}

// MapFunc returns a slice mapped by f(c) within all c in the slice.
// MapFunc does not modify the contents of the slice s; it creates a new slice.
// TODO: accept [S ~[]E, E any, R ~[]M, M any] if go support template type deduction
func MapFunc[S ~[]E, E any, R []M, M any](s S, f func(E) M) R {
	if s == nil {
		return nil
	}

	var rr = make(R, len(s))
	for i, v := range s {
		rr[i] = f(v)
	}
	return rr
}
