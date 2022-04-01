// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
	"golang.org/x/exp/slices"
)

var equalIntTests = []struct {
	s    []int
	n    int
	want [][]int
}{
	{
		nil,
		0,
		nil,
	},
	{
		[]int{},
		0,
		nil,
	},
	{
		[]int{1},
		0,
		nil,
	},
	{
		[]int{1},
		1,
		[][]int{{1}},
	},
	{
		[]int{1},
		2,
		[][]int{{1}},
	},
	{
		[]int{1, 2, 3},
		-1,
		[][]int{{1}, {2}, {3}},
	},
	{
		[]int{1, 2, 3},
		0,
		nil,
	},
	{
		[]int{1, 2, 3},
		1,
		[][]int{{1, 2, 3}},
	},
	{
		[]int{1, 2, 3},
		2,
		[][]int{{1}, {2, 3}},
	},
	{
		[]int{1, 2, 3},
		3,
		[][]int{{1}, {2}, {3}},
	},
	{
		[]int{1, 2, 3},
		4,
		[][]int{{1}, {2}, {3}},
	},
}

func TestSplitN(t *testing.T) {
	for _, test := range equalIntTests {
		got := slices_.SplitN(test.s, test.n)
		if len(test.want) != len(got) {
			t.Errorf("SplitN(%v, %v) = %v, want %v", test.s, test.n, got, test.want)
			continue
		}
		for i := range got {
			if !slices.Equal(got[i], test.want[i]) {
				t.Errorf("SplitN(%v, %v) = %v, want %v", test.s, test.n, got, test.want)
				break
			}
		}
	}
}
