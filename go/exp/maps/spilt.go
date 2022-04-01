// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps

// SplitN slices s into submaps and returns a map of the submaps .
//
// The count determines the number of submaps to return:
//   n > 0: at most n submaps; the last submaps will be the unsplit remainder.
//   n == 0: the result is nil (zero submaps)
//   n < 0: all submaps as n == len(s)
//
// Edge cases for s and sep (for example, empty strings) are handled
// as described in the documentation for Split.
func SplitN[M ~map[K]V, K comparable, V any](m M, n int) []M {
	if m == nil || n == 0 {
		return nil
	}
	if n < 0 {
		n = len(m)
	}
	if n <= 1 || len(m) <= 1 {
		return []M{m}
	}

	chunkSize := len(m) / n
	if chunkSize == 0 {
		chunkSize = 1
	}
	var chunks []M

	var chunk = make(M)
	for k, v := range m {
		if len(chunks) == n-1 || chunkSize > len(m) {
			chunkSize = len(m)
		}
		if len(chunk) < chunkSize {
			chunk[k] = v
		}
		if len(chunk) == chunkSize {
			chunks = append(chunks, chunk)
			chunk = make(M)
		}
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}
