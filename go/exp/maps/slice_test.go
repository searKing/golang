// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	maps_ "github.com/searKing/golang/go/exp/maps"
)

func TestSliceFunc(t *testing.T) {
	type KV struct {
		K int
		V string
	}
	sort := func(a, b KV) int {
		if a.K == b.K {
			return strings.Compare(a.V, b.V)
		}
		return a.K - b.K
	}
	tests := []struct {
		data map[int]string
		want []KV
	}{
		{nil, nil},
		{map[int]string{}, []KV{}},
		{map[int]string{0: "0"}, []KV{{0, "0"}}},
		{map[int]string{1: "1", 0: "0"}, []KV{{1, "1"}, {0, "0"}}},
		{map[int]string{1: "1", 2: "2"}, []KV{{1, "1"}, {2, "2"}}},
		{map[int]string{0: "0", 1: "1", 2: "2"}, []KV{{0, "0"}, {1, "1"}, {2, "2"}}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			got := maps_.SliceFunc(tt.data, func(k int, v string) KV {
				return KV{k, v}
			})
			slices.SortFunc(got, sort)
			slices.SortFunc(tt.want, sort)
			if len(got) != len(tt.want) {
				t.Errorf("SliceFunc(%v) = %v, want %v", tt.data, got, tt.want)
			}
			if slices.CompareFunc(got, tt.want, sort) != 0 {
				t.Errorf("SliceFunc(%v) = %v, want %v", tt.data, got, tt.want)
			}
		})
	}
}
