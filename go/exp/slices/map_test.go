// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"fmt"
	"slices"
	"strconv"
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
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
			{
				got := slices_.Map(tt.data)

				if slices.Compare(got, tt.want) != 0 {
					t.Errorf("slices_.Map(%v) = %v, want %v", tt.data, got, tt.want)
				}
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
			{
				got := slices_.MapFunc(tt.data, func(e int) string {
					return strconv.Itoa(e)
				})

				if slices.Compare(got, tt.want) != 0 {
					t.Errorf("slices_.MapFunc(%v, func(e int) string {return strconv.Itoa(e)}) = %v, want %v", tt.data, got, tt.want)
				}
			}
		})
	}
}

func TestMapIndexFunc(t *testing.T) {
	tests := []struct {
		data []int
		want []string
	}{
		{nil, nil},
		{[]int{}, []string{}},
		{[]int{0}, []string{"0:0"}},
		{[]int{1, 0}, []string{"0:1", "1:0"}},
		{[]int{1, 2}, []string{"0:1", "1:2"}},
		{[]int{0, 1, 2}, []string{"0:0", "1:1", "2:2"}},
		{[]int{0, 1, 0, 2}, []string{"0:0", "1:1", "2:0", "3:2"}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				got := slices_.MapIndexFunc(tt.data, func(i int, e int) string {
					return fmt.Sprintf("%d:%d", i, e)
				})

				if slices.Compare(got, tt.want) != 0 {
					t.Errorf("slices_.MapIndexFunc(%v, func(i int, e int) string {return fmt.Sprintf(\"%%d:%%d\", i, e)}) = %v, want %v", tt.data, got, tt.want)
				}
			}
		})
	}
}
