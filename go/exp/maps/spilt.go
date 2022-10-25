// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps

// Split slices s into all submaps separated by sep and returns a slice of
// the submaps between those separators.
//
// If s is less than sep and sep is more than zero, Split returns a
// slice of length 1 whose only element is s.
//
// If s is nil, Split returns nil (zero submaps).
//
// If both s and sep are empty or zero, Split returns an empty slice.
//
// If sep is <= zero, Split splits after each element, as chunk size is 1.
//
// It is equivalent to SplitN with a count of -1.
func Split[M ~map[K]V, K comparable, V any](m M, sep int) []M {
	return SplitN(m, sep, -1)
}

// SplitN slices s into submaps and returns a slice of the submaps.
//
// The count determines the number of submaps to return:
//
//	  n > 0: at most n submaps; the last submaps will be the unsplit remainder.
//			The count determines the number of submaps to return:
//	  		sep > 0: Split splits every sep as chunk size; the last submaps will be the unsplit remainder.
//	  		sep <= 0: take len(S)/n as chunk size
//	  n == 0: the result is nil (zero submaps)
//	  n < 0: all submaps as n == len(s)
//
// Edge cases for s and sep (for example, zero) are handled
// as described in the documentation for Split.
func SplitN[M ~map[K]V, K comparable, V any](m M, sep int, n int) []M {
	// n < 0: all submaps as n == len(s)
	if n < 0 {
		n = len(m)
	}
	// Below: n >= 0

	// n == 0: the result is nil (zero submaps)
	// If s is nil, Split returns nil (zero submaps).
	if n == 0 || m == nil {
		return nil
	}

	// Below: s != nil && n > 0

	// If both s and sep are empty or zero, Split returns an empty slice.
	if len(m) == 0 && sep == 0 {
		return []M{}
	}

	// If s is less or equal than sep and sep is more than zero, Split returns a
	// slice of length 1 whose only element is s.
	if len(m) <= sep && sep > 0 {
		return []M{m}
	}

	// Below: len(s) > 0 && len(s) > sep && n > 0

	// n > 0: at most n submaps; the last submaps will be the unsplit remainder.
	//      The count determines the number of submaps to return:
	//      sep > 0: Split splits every sep as chunk size; the last submaps will be the unsplit remainder.
	//	  	sep <= 0: take len(S)/n as chunk size

	if n == 1 || len(m) == 1 {
		return []M{m}
	}

	// Below: len(s) > 1 && len(s) > sep && n > 1

	// If sep is <= zero, Split splits after each element, as chunk size is 1.
	chunkSize := len(m) / n
	if chunkSize == 0 {
		chunkSize = 1
	}
	if sep > 0 {
		chunkSize = sep
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
