// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"fmt"
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
	"golang.org/x/exp/slices"
)

var unionTests = []struct {
	s1   []int
	s2   []int
	want []int
}{
	{
		[]int{},
		[]int{},
		[]int{},
	},
	{
		[]int{1, 2, 3},
		[]int{-1, -2, -3},
		[]int{1, 2, 3, -1, -2, -3},
	},
	{
		[]int{1, 2, 3},
		[]int{1, -2, -3},
		[]int{1, 2, 3, -2, -3},
	},
	{
		[]int{1, 2, 3},
		[]int{1, -2, 3},
		[]int{1, 2, 3, -2},
	},
	{
		[]int{3, 2, 1},
		[]int{-3, -2, -1},
		[]int{3, 2, 1, -3, -2, -1},
	},
}

func TestUnionFunc(t *testing.T) {
	for _, tt := range unionTests {
		t.Run(fmt.Sprintf("slices_.Union(%v, %v, equal[int])", tt.s1, tt.s2), func(t *testing.T) {
			{
				got := slices_.UnionFunc(tt.s1, tt.s2, equal[int])
				if !slices.Equal(got, tt.want) {
					t.Errorf("slices_.Union(%v, %v, equal[int]) = %v, want %v", tt.s1, tt.s2, got, tt.want)
				}
			}
		})
	}
}

func TestUnion(t *testing.T) {
	for _, tt := range unionTests {
		t.Run(fmt.Sprintf("slices_.Union(%v, %v)", tt.s1, tt.s2), func(t *testing.T) {
			{
				got := slices_.Union(tt.s1, tt.s2)
				slices.Sort(got)
				slices.Sort(tt.want)
				if !slices.Equal(got, tt.want) {
					t.Errorf("slices_.Union(%v, %v) = %v, want %v", tt.s1, tt.s2, got, tt.want)
				}
			}
		})
	}
}
