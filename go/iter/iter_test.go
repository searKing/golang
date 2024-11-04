// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter_test

import (
	"fmt"
	"slices"
	"testing"

	iter_ "github.com/searKing/golang/go/iter"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		data []int
		want []int
	}{
		{nil, nil},
		{[]int{}, []int{}},
		{[]int{0}, []int{}},
		{[]int{1, 0}, []int{1}},
		{[]int{1, 2}, []int{1, 2}},
		{[]int{0, 1, 2}, []int{1, 2}},
		{[]int{0, 1, 0, 2}, []int{1, 2}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d: %v", i, tt.data), func(t *testing.T) {
			got := slices.Collect(iter_.Filter(slices.Values(tt.data)))
			if !slices.Equal(got, tt.want) {
				t.Errorf("iter_.Filter(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestFilterFunc(t *testing.T) {
	tests := []struct {
		data []int
		want []int
	}{
		{nil, nil},
		{[]int{}, []int{}},
		{[]int{0}, []int{}},
		{[]int{1, 0}, []int{1}},
		{[]int{1, 2}, []int{1, 2}},
		{[]int{0, 1, 2}, []int{1, 2}},
		{[]int{0, 1, 0, 2}, []int{1, 2}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d: %v", i, tt.data), func(t *testing.T) {
			got := slices.Collect(iter_.FilterFunc(slices.Values(tt.data), func(e int) bool {
				return e != 0
			}))
			if !slices.Equal(got, tt.want) {
				t.Errorf("iter_.FilterFunc(%v, func(e int) bool {return e != 0}) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestFilterN(t *testing.T) {
	tests := []struct {
		data []int
		n    int
		want []int
	}{
		{nil, 0, nil},
		{[]int{}, -1, []int{}},
		{[]int{}, 0, []int{}},
		{[]int{}, 1, []int{}},
		{[]int{0}, -1, []int{0}},
		{[]int{0}, 0, []int{}},
		{[]int{0}, 1, []int{0}},
		{[]int{1, 2}, -1, []int{1, 2}},
		{[]int{1, 2}, 0, []int{}},
		{[]int{1, 2}, 1, []int{1}},
		{[]int{1, 2}, 2, []int{1, 2}},
		{[]int{1, 2}, 3, []int{1, 2}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%d: %v, %d", i, tt.data, tt.n), func(t *testing.T) {
			got := slices.Collect(iter_.FilterN(slices.Values(tt.data), tt.n))
			if !slices.Equal(got, tt.want) {
				t.Errorf("iter_.FilterN(%v, %v) = %v, want %v", tt.data, tt.n, got, tt.want)
			}
		})
	}
}
