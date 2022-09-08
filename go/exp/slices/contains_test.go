// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
)

var containsTests = []struct {
	s    []int
	v    int
	want int
}{
	{
		nil,
		0,
		-1,
	},
	{
		[]int{},
		0,
		-1,
	},
	{
		[]int{1, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 3, 2},
		2,
		1,
	},
}

func equalToContains[T any](f func(T, T) bool, v1 T) func(T) bool {
	return func(v2 T) bool {
		return f(v1, v2)
	}
}

func TestContainsFunc(t *testing.T) {
	for _, test := range containsTests {
		if got := slices_.ContainsFunc(test.s, equalToContains(equal[int], test.v)); got != (test.want != -1) {
			t.Errorf("ContainsFunc(%v, equalToContains(equal[int], %v)) = %t, want %d", test.s, test.v, got, test.want)
		}
	}
}

func TestContains(t *testing.T) {
	for _, test := range containsTests {
		if got := slices_.Contains(test.s, test.v); got != (test.want != -1) {
			t.Errorf("Contains(%v, %v) = %t, want %t", test.s, test.v, got, test.want != -1)
		}
	}
}
