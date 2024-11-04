// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iter_test

import (
	"fmt"
	"maps"
	"testing"

	iter_ "github.com/searKing/golang/go/iter"
)

func TestFilter2(t *testing.T) {
	tests := []struct {
		data map[int]int
		want map[int]int
	}{
		{nil, nil},
		{map[int]int{}, map[int]int{}},
		{map[int]int{0: 0}, map[int]int{}},
		{map[int]int{1: 1, 0: 0}, map[int]int{1: 1}},
		{map[int]int{1: 1, 2: 2}, map[int]int{1: 1, 2: 2}},
		{map[int]int{0: 0, 1: 1, 2: 2}, map[int]int{1: 1, 2: 2}},
		{map[int]int{0: 0, 1: 1, 3: 0, 2: 2}, map[int]int{1: 1, 2: 2}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			got := maps.Collect(iter_.Filter2(maps.All(tt.data)))
			if !maps.Equal(got, tt.want) {
				t.Errorf("maps_.Filter2(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestFilter2Func(t *testing.T) {
	tests := []struct {
		data map[int]int
		want map[int]int
	}{
		{nil, nil},
		{map[int]int{}, map[int]int{}},
		{map[int]int{0: 0}, map[int]int{}},
		{map[int]int{1: 1, 0: 0}, map[int]int{1: 1}},
		{map[int]int{1: 1, 2: 2}, map[int]int{1: 1, 2: 2}},
		{map[int]int{0: 0, 1: 1, 2: 2}, map[int]int{1: 1, 2: 2}},
		{map[int]int{0: 0, 1: 1, 3: 0, 2: 2}, map[int]int{1: 1, 2: 2}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			got := maps.Collect(iter_.Filter2Func(maps.All(tt.data), func(k, v int) bool { return v != 0 }))
			if !maps.Equal(got, tt.want) {
				t.Errorf("maps_.Filter2Func(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}

func TestFilter2N(t *testing.T) {
	tests := []struct {
		data  map[int]int
		n     int
		wants []map[int]int
	}{
		{nil, 0, []map[int]int{nil}},
		{map[int]int{}, -1, []map[int]int{{}}},
		{map[int]int{}, 0, []map[int]int{{}}},
		{map[int]int{}, 1, []map[int]int{{}}},
		{map[int]int{0: 0}, -1, []map[int]int{{0: 0}}},
		{map[int]int{0: 0}, 0, []map[int]int{{}}},
		{map[int]int{0: 0}, 1, []map[int]int{{0: 0}}},
		{map[int]int{1: 1, 2: 2}, -1, []map[int]int{{1: 1, 2: 2}}},
		{map[int]int{1: 1, 2: 2}, 0, []map[int]int{{}}},
		{map[int]int{1: 1, 2: 2}, 1, []map[int]int{{1: 1}, {2: 2}}},
		{map[int]int{1: 1, 2: 2}, 2, []map[int]int{{1: 1, 2: 2}}},
		{map[int]int{1: 1, 2: 2}, 3, []map[int]int{{1: 1, 2: 2}}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			got := maps.Collect(iter_.Filter2N(maps.All(tt.data), tt.n))
			var eq bool
			for _, want := range tt.wants {
				if maps.Equal(got, want) {
					eq = true
					break
				}
			}
			if !eq {
				t.Errorf("maps_.Filter2N(%v, %d) = %v, want %v", tt.data, tt.n, got, tt.wants)
			}
		})
	}
}
