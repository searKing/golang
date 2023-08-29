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

func TestReverse(t *testing.T) {
	tests := []struct {
		data []int
		want []int
	}{
		{nil, nil},
		{[]int{}, []int{}},
		{[]int{0}, []int{0}},
		{[]int{1, 0}, []int{0, 1}},
		{[]int{2, 1, 0}, []int{0, 1, 2}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				slices_.Reverse(tt.data)
				if !slices.Equal(tt.data, tt.want) {
					t.Errorf("reversed %v", tt.want)
					t.Errorf("   got %v", tt.data)
				}
			}
		})
	}
}
