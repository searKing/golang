// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"fmt"
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
)

func TestFirstOrZero_Int(t *testing.T) {
	tests := []struct {
		data []int
		want int
	}{
		{nil, 0},
		{[]int{}, 0},
		{[]int{0}, 0},
		{[]int{1, 0}, 1},
		{[]int{1, 2}, 1},
		{[]int{0, 1, 2}, 1},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				got := slices_.FirstOrZero(tt.data...)
				if got != tt.want {
					t.Errorf("first of %v = %v", tt.data, tt.want)
					t.Errorf("   got %v", got)
				}
			}
		})
	}
}

func TestFirstOrZero_String(t *testing.T) {
	tests := []struct {
		data []string
		want string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{""}, ""},
		{[]string{"a", ""}, "a"},
		{[]string{"a", "b"}, "a"},
		{[]string{"", "a", "b"}, "a"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				got := slices_.FirstOrZero(tt.data...)
				if got != tt.want {
					t.Errorf("first of %v = %v", tt.data, tt.want)
					t.Errorf("   got %v", got)
				}
			}
		})
	}
}
