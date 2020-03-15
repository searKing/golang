// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package preconditions

import "errors"

type ErrorOutOfBound error

var (
	errorOutOfBound = errors.New("out of bound")
)

// CheckIndex checks if the {@code index} is within the bounds of the range from
// {@code 0} (inclusive) to {@code length} (exclusive).
func CheckIndex(index, length int) int {
	if index < 0 || index >= length {
		panic(ErrorOutOfBound(errors.New("CheckIndex")))
	}
	return index
}

// CheckFromToIndex checks if the sub-range from {@code fromIndex} (inclusive) to
// {@code toIndex} (exclusive) is within the bounds of range from {@code 0}
// (inclusive) to {@code length} (exclusive).
func CheckFromToIndex(fromIndex, toIndex, length int) int {
	if fromIndex < 0 || fromIndex > toIndex || toIndex > length {
		panic(ErrorOutOfBound(errors.New("CheckFromToIndex")))
	}
	return fromIndex
}

// Checks if the sub-range from {@code fromIndex} (inclusive) to
// {@code fromIndex + size} (exclusive) is within the bounds of range from
// {@code 0} (inclusive) to {@code length} (exclusive).
func CheckFromIndexSize(fromIndex, size, length int) int {
	if length < 0 || fromIndex < 0 || size < 0 || size > length-fromIndex {
		panic(ErrorOutOfBound(errors.New("CheckFromIndexSize")))
	}
	return fromIndex
}
