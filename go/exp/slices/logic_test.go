// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"fmt"
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
)

func TestOrFunc(t *testing.T) {
	tests := []struct {
		data []int
		want bool
	}{
		{nil, true},
		{[]int{}, true},
		{[]int{0}, false},
		{[]int{1, 0}, true},
		{[]int{1, 2}, true},
		{[]int{0, 1, 2}, true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				got := slices_.OrFunc(tt.data, func(v int) bool { return v != 0 })
				if got != tt.want {
					t.Errorf("and of %v = %t", tt.data, tt.want)
					t.Errorf("   got %t", got)
				}
			}
		})
	}
}

func TestAndFunc(t *testing.T) {
	tests := []struct {
		data []int
		want bool
	}{
		{nil, true},
		{[]int{}, true},
		{[]int{0}, false},
		{[]int{1, 0}, false},
		{[]int{1, 2}, true},
		{[]int{0, 1, 2}, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				got := slices_.AndFunc(tt.data, func(v int) bool { return v != 0 })
				if got != tt.want {
					t.Errorf("and of %v = %t", tt.data, tt.want)
					t.Errorf("   got %t", got)
				}
			}
		})
	}
}

func TestOr(t *testing.T) {
	tests := []struct {
		data []int
		want bool
	}{
		{nil, true},
		{[]int{}, true},
		{[]int{0}, false},
		{[]int{1, 0}, true},
		{[]int{1, 2}, true},
		{[]int{0, 1, 2}, true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				got := slices_.Or(tt.data...)
				if got != tt.want {
					t.Errorf("and of %v = %t", tt.data, tt.want)
					t.Errorf("   got %t", got)
				}
			}
		})
	}
}

func TestAnd(t *testing.T) {
	tests := []struct {
		data []int
		want bool
	}{
		{nil, true},
		{[]int{}, true},
		{[]int{0}, false},
		{[]int{1, 0}, false},
		{[]int{1, 2}, true},
		{[]int{0, 1, 2}, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				got := slices_.And(tt.data...)
				if got != tt.want {
					t.Errorf("and of %v = %t", tt.data, tt.want)
					t.Errorf("   got %t", got)
				}
			}
		})
	}
}
