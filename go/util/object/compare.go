// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package object

import "reflect"

type Comparator interface {
	Compare(a, b interface{}) int
}

// Equals Returns {@code true} if the arguments are equal to each other
// and {@code false} otherwise.
func Equals(a, b interface{}) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return a == b
	}
	v1 := reflect.ValueOf(a)
	v2 := reflect.ValueOf(b)
	if v1.Type() != v2.Type() {
		return false
	}
	return a == b
}

// DeepEquals Returns {@code true} if the arguments are deeply equal to each other
// and {@code false} otherwise.
func DeepEquals(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// Returns 0 if the arguments are identical and {@code
// c.compare(a, b)} otherwise.
// Consequently, if both arguments are {@code null} 0
// is returned.
func Compare(a, b interface{}, compare Comparator) int {
	if a == b {
		return 0
	}
	return compare.Compare(a, b)
}
