// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"fmt"
	"slices"
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
)

var diffTests = []struct {
	sOld    []int
	sNew    []int
	wantAdd []int
	wantDel []int
}{
	{
		[]int{},
		[]int{},
		[]int{},
		[]int{},
	},
	{
		[]int{1, 2, 3},
		[]int{-1, -2, -3},
		[]int{-1, -2, -3},
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		[]int{3, 2, 1},
		[]int{},
		[]int{},
	},
	{
		[]int{1, 2, 3},
		[]int{1, -2, -3},
		[]int{-2, -3},
		[]int{2, 3},
	},
}

func TestDiff(t *testing.T) {
	for _, tt := range diffTests {
		t.Run(fmt.Sprintf("slices_.Diff(%v, %v)", tt.sNew, tt.sOld), func(t *testing.T) {
			{
				gotAdd, gotDel := slices_.Diff(tt.sNew, tt.sOld)
				if !slices.Equal(gotAdd, tt.wantAdd) {
					t.Errorf("slices_.Diff(%v, %v) = (%v, %v), want (%v, %v)", tt.sNew, tt.sOld, gotAdd, gotDel, tt.wantAdd, tt.wantDel)
				}
			}
		})
	}
}
