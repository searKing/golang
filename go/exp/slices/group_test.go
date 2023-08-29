// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"slices"
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
)

var groupTests = []struct {
	s    []int
	want map[int]int
}{
	{nil, nil},
	{[]int{}, map[int]int{}},
	{[]int{1}, map[int]int{1: 1}},
	{[]int{1, 1}, map[int]int{1: 2}},
	{[]int{1, 2, 3}, map[int]int{1: 1, 2: 1, 3: 1}},
	{[]int{1, 1, 2, 3, 3}, map[int]int{1: 2, 2: 1, 3: 2}},
	{[]int{1, 2, 3, 4}, map[int]int{1: 1, 2: 1, 3: 1, 4: 1}},
	{[]int{1, 2, 3, 4, 2, 4, 2}, map[int]int{1: 1, 2: 3, 3: 1, 4: 2}},
}

func TestGroup(t *testing.T) {
	for i, test := range groupTests {
		got := slices_.Group(test.s)
		if len(test.want) != len(got) {
			t.Errorf("#%d: Group(%v) = %v, want %v", i, test.s, got, test.want)
			continue
		}
		for k := range got {
			if got[k] != test.want[k] {
				t.Errorf("#%d: Group(%v) = %v, want %v", i, test.s, got, test.want)
				break
			}
		}
	}
}

var groupFuncTests = []struct {
	s    []int
	want map[int][]int
}{
	{nil, nil},
	{[]int{}, map[int][]int{}},
	{[]int{1}, map[int][]int{1: {1}}},
	{[]int{1, 1}, map[int][]int{1: {1, 1}}},
	{[]int{1, 2, 3}, map[int][]int{1: {1}, 2: {2}, 3: {3}}},
	{[]int{1, 1, 2, 3, 3}, map[int][]int{1: {1, 1}, 2: {2}, 3: {3, 3}}},
	{[]int{1, 2, 3, 4}, map[int][]int{1: {1}, 2: {2}, 3: {3}, 4: {4}}},
	{[]int{1, 2, 3, 4, 2, 4, 2}, map[int][]int{1: {1}, 2: {2, 2, 2}, 3: {3}, 4: {4, 4}}},
}

func TestGroupFunc(t *testing.T) {
	for i, test := range groupFuncTests {
		got := slices_.GroupFunc(test.s, func(i int) int { return i })
		if len(test.want) != len(got) {
			t.Errorf("#%d: Group(%v) = %v, want %v", i, test.s, got, test.want)
			continue
		}
		for k := range got {
			if !slices.Equal(got[k], test.want[k]) {
				t.Errorf("#%d: Group(%v) = %v, want %v", i, test.s, got, test.want)
				break
			}
		}
	}
}
