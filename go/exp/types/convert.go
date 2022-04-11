// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// Pointer returns a pointer to the value passed in.
func Pointer[T any](v T) *T {
	return &v
}

// Value returns the value of the pointer passed in or
// "" if the pointer is nil.
func Value[T any](v *T) T {
	if v != nil {
		return *v
	}
	var zero T
	return zero
}

// PointerSlice converts a slice of values into a slice of pointers
func PointerSlice[S ~[]E, V []*E, E any](src S) V {
	dst := make(V, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}
	return dst
}

// ValueSlice converts a slice of pointers into a slice of values
func ValueSlice[S ~[]*E, V []E, E any](src S) V {
	dst := make(V, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

// PointerMap converts a map of values into a map of pointers
func PointerMap[M ~map[K]V, N map[K]*V, K comparable, V any](src M) N {
	dst := make(N)
	for k, val := range src {
		v := val
		dst[k] = &v
	}
	return dst
}

// ValueMap converts a map of string pointers into a string map of values
func ValueMap[M ~map[K]*V, N map[K]V, K comparable, V any](src M) N {
	dst := make(N)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}
	return dst
}
