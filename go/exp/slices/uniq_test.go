// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"slices"
	"strings"
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
)

var uniqTests = []struct {
	s    []int
	want []int
}{
	{
		nil,
		nil,
	},
	{
		[]int{1},
		[]int{1},
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
	},
	{
		[]int{1, 1, 2},
		[]int{1, 2},
	},
	{
		[]int{1, 2, 1},
		[]int{1, 2},
	},
	{
		[]int{1, 2, 2, 3, 3, 4},
		[]int{1, 2, 3, 4},
	},
}

func TestUniq(t *testing.T) {
	for _, test := range uniqTests {
		copy := slices.Clone(test.s)
		got := slices_.Uniq(copy)
		slices.Sort(got)
		slices.Sort(test.want)
		if !slices.Equal(got, test.want) {
			t.Errorf("Uniq(%v) = %v, want %v", test.s, got, test.want)
		}
	}
}

// equal is simply ==.
func equal[T comparable](v1, v2 T) bool {
	return v1 == v2
}

func TestUniqFunc(t *testing.T) {
	for _, test := range uniqTests {
		copy := slices.Clone(test.s)
		if got := slices_.UniqFunc(copy, equal[int]); !slices.Equal(got, test.want) {
			t.Errorf("UniqFunc(%v, equal[int]) = %v, want %v", test.s, got, test.want)
		}
	}

	s1 := []string{"a", "a", "A", "B", "b"}
	want := []string{"a", "B"}
	if got := slices_.UniqFunc(s1, strings.EqualFold); !slices.Equal(got, want) {
		t.Errorf("UniqFunc(%v, strings.EqualFold) = %v, want %v", s1, got, want)
	}
}
