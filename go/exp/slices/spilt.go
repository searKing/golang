// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// Split slices s into all subslices separated by sep and returns a slice of
// the subslices between those separators.
//
// If s is less than sep and sep is more than zero, Split returns a
// slice of length 1 whose only element is s.
//
// If s is nil, Split returns nil (zero subslices).
//
// If both s and sep are empty or zero, Split returns an empty slice.
//
// If sep is <= zero, Split splits after each element, as chunk size is 1.
//
// It is equivalent to SplitN with a count of -1.
func Split[S ~[]E, E any](s S, sep int) []S {
	return SplitN(s, sep, -1)
}

// SplitN slices s into subslices and returns a slice of the subslices.
//
// The count determines the number of subslices to return:
//
//	  n > 0: at most n subslices; the last subslices will be the unsplit remainder.
//			The count determines the number of subslices to return:
//	  		sep > 0: Split splits every sep as chunk size; the last subslices will be the unsplit remainder.
//	  		sep <= 0: take len(S)/n as chunk size
//	  n == 0: the result is nil (zero subslices)
//	  n < 0: all subslices as n == len(s)
//
// Edge cases for s and sep (for example, zero) are handled
// as described in the documentation for Split.
func SplitN[S ~[]E, E any](s S, sep int, n int) []S {
	// n < 0: all subslices as n == len(s)
	if n < 0 {
		n = len(s)
	}
	// Below: n >= 0

	// n == 0: the result is nil (zero subslices)
	// If s is nil, Split returns nil (zero subslices).
	if n == 0 || s == nil {
		return nil
	}

	// Below: s != nil && n > 0

	// If both s and sep are empty or zero, Split returns an empty slice.
	if len(s) == 0 && sep == 0 {
		return []S{}
	}

	// If s is less or equal than sep and sep is more than zero, Split returns a
	// slice of length 1 whose only element is s.
	if len(s) <= sep && sep > 0 {
		return []S{s}
	}

	// Below: len(s) > 0 && len(s) > sep && n > 0

	// n > 0: at most n subslices; the last subslices will be the unsplit remainder.
	//      The count determines the number of subslices to return:
	//      sep > 0: Split splits every sep as chunk size; the last subslices will be the unsplit remainder.
	//	  	sep <= 0: take len(S)/n as chunk size

	if n == 1 || len(s) == 1 {
		return []S{s}
	}

	// Below: len(s) > 1 && len(s) > sep && n > 1

	// If sep is <= zero, Split splits after each element, as chunk size is 1.
	chunkSize := len(s) / n
	if chunkSize == 0 {
		chunkSize = 1
	}
	if sep > 0 {
		chunkSize = sep
	}
	var chunks []S
	for len(s) > 0 {
		if len(chunks) == n-1 || chunkSize > len(s) {
			chunkSize = len(s)
		}
		chunk := append([]E{}, s[:chunkSize]...)
		s = s[chunkSize:]
		chunks = append(chunks, S(chunk))
	}
	return chunks
}

// SplitMap slices s into all key-value pairs and returns a map of the key-value pairs.
//
// If s is nil, SplitMap returns nil (zero map).
//
// If len(s) is odd, it treats args[len(args)-1] as a value with a missing value.
func SplitMap[M ~map[K]V, S ~[]E, K comparable, V any, E any](s S) M {
	if s == nil {
		return nil
	}
	kvs := Split(s, 2)
	var m = make(M, len(kvs))
	for _, kv := range kvs {
		switch len(kv) {
		case 1:
			var zeroV V
			m[any(kv[0]).(K)] = zeroV
		case 2:
			m[any(kv[0]).(K)] = any(kv[1]).(V)
		default:
		}
	}
	return m
}
