// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// Split slices s into all substrings separated by sep and returns a slice of
// the substrings between those separators.
//
// If sep is <= zero, Split splits after each element, as chunk size is 1. If both s
// and sep are empty or zero, Split returns an empty slice.
//
// It is equivalent to SplitN with a count of -1.
func Split[S ~[]E, E any](s S, sep int) []S {
	return SplitN(s, sep, -1)
}

// SplitN slices s into subslices and returns a slice of the subslices.
//
// The count determines the number of subslices to return:
//   n > 0: at most n subslices; the last subslices will be the unsplit remainder.
// 		The count determines the number of subslices to return:
//   		sep > 0: Split splits every sep as chunk size; the last subslices will be the unsplit remainder.
//   		sep <= 0: take len(S)/n as chunk size
//   n == 0: the result is nil (zero subslices)
//   n < 0: all subslices as n == len(s)
//
// Edge cases for s and sep (for example, zero) are handled
// as described in the documentation for Split.
func SplitN[S ~[]E, E any](s S, sep int, n int) []S {
	if s == nil || n == 0 {
		return nil
	}
	if n < 0 {
		n = len(s)
	}
	if n <= 1 || len(s) <= 1 {
		return []S{s}
	}

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
