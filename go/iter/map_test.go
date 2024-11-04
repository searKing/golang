// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter_test

import (
	"fmt"
	"slices"
	"strconv"
	"testing"

	iter_ "github.com/searKing/golang/go/iter"
)

func TestMap(t *testing.T) {
	tests := []struct {
		data []int
		want []string
	}{
		{nil, nil},
		{[]int{}, []string{}},
		{[]int{0}, []string{"0"}},
		{[]int{1, 0}, []string{"1", "0"}},
		{[]int{1, 2}, []string{"1", "2"}},
		{[]int{0, 1, 2}, []string{"0", "1", "2"}},
		{[]int{0, 1, 0, 2}, []string{"0", "1", "0", "2"}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			got := slices.Collect(iter_.Map(slices.Values(tt.data)))
			if slices.Compare(got, tt.want) != 0 {
				t.Errorf("iter_.Map(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestMapFunc(t *testing.T) {
	tests := []struct {
		data []int
		want []string
	}{
		{nil, nil},
		{[]int{}, []string{}},
		{[]int{0}, []string{"0"}},
		{[]int{1, 0}, []string{"1", "0"}},
		{[]int{1, 2}, []string{"1", "2"}},
		{[]int{0, 1, 2}, []string{"0", "1", "2"}},
		{[]int{0, 1, 0, 2}, []string{"0", "1", "0", "2"}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			got := slices.Collect(iter_.MapFunc(slices.Values(tt.data), func(e int) string { return strconv.Itoa(e) }))
			if slices.Compare(got, tt.want) != 0 {
				t.Errorf("iter_.MapFunc(%v, func(e int) string {return strconv.Itoa(e)}) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}
