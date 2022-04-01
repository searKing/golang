// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

// SplitN slices s into subslices and returns a slice of the subslices .
//
// The count determines the number of subslices to return:
//   n > 0: at most n subslices; the last subslices will be the unsplit remainder.
//   n == 0: the result is nil (zero subslices)
//   n < 0: all subslices as n == len(s)
//
// Edge cases for s and sep (for example, empty strings) are handled
// as described in the documentation for Split.
func SplitN[S ~[]E, E any](s S, n int) []S {
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
	var chunks []S
	for len(s) > 0 {
		if len(chunks) == n-1 || chunkSize > len(s) {
			chunkSize = len(s)
		}
		chunk := append(S([]E{}), s[:chunkSize]...)
		s = s[chunkSize:]
		chunks = append(chunks, S(chunk))
	}
	return chunks
}
