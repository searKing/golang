// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"fmt"
	"slices"
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
)

var intersectTests = []struct {
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
		[]int{},
	},
	{
		[]int{1, 2, 3},
		[]int{1, -2, -3},
		[]int{1},
	},
	{
		[]int{1, 2, 3},
		[]int{1, -2, 3},
		[]int{1, 3},
	},
	{
		[]int{3, 2, 1},
		[]int{-3, -2, -1},
		[]int{},
	},
}

func TestIntersectFunc(t *testing.T) {
	for _, tt := range intersectTests {
		t.Run(fmt.Sprintf("slices_.Intersect(%v, %v, equal[int])", tt.s1, tt.s2), func(t *testing.T) {
			{
				got := slices_.IntersectFunc(tt.s1, tt.s2, equal[int])
				if !slices.Equal(got, tt.want) {
					t.Errorf("slices_.Intersect(%v, %v, equal[int]) = %v, want %v", tt.s1, tt.s2, got, tt.want)
				}
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	for _, tt := range intersectTests {
		t.Run(fmt.Sprintf("slices_.Intersect(%v, %v)", tt.s1, tt.s2), func(t *testing.T) {
			{
				got := slices_.Intersect(tt.s1, tt.s2)
				if !slices.Equal(got, tt.want) {
					t.Errorf("slices_.Intersect(%v, %v) = %v, want %v", tt.s1, tt.s2, got, tt.want)
				}
			}
		})
	}
}
