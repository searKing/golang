// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps_test

import (
	"fmt"
	"maps"
	"strconv"
	"testing"

	maps_ "github.com/searKing/golang/go/exp/maps"
)

func TestMap(t *testing.T) {
	tests := []struct {
		data map[int]int
		want map[string]string
	}{
		{nil, nil},
		{map[int]int{}, map[string]string{}},
		{map[int]int{0: 0}, map[string]string{"0": "0"}},
		{map[int]int{0: 1, 1: 0}, map[string]string{"0": "1", "1": "0"}},
		{map[int]int{0: 1, 1: 2}, map[string]string{"0": "1", "1": "2"}},
		{map[int]int{0: 0, 1: 1, 2: 2}, map[string]string{"0": "0", "1": "1", "2": "2"}},
		{map[int]int{0: 0, 1: 1, 2: 0, 3: 2}, map[string]string{"0": "0", "1": "1", "2": "0", "3": "2"}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				got := maps_.Map(tt.data)

				if !maps.Equal(got, tt.want) {
					t.Errorf("maps_.Map(%v) = %v, want %v", tt.data, got, tt.want)
				}
			}
		})
	}
}

func TestMapFunc(t *testing.T) {
	tests := []struct {
		data map[int]int
		want map[string]string
	}{
		{nil, nil},
		{map[int]int{}, map[string]string{}},
		{map[int]int{0: 0}, map[string]string{"0": "0"}},
		{map[int]int{0: 1, 1: 0}, map[string]string{"0": "1", "1": "0"}},
		{map[int]int{0: 1, 1: 2}, map[string]string{"0": "1", "1": "2"}},
		{map[int]int{0: 0, 1: 1, 2: 2}, map[string]string{"0": "0", "1": "1", "2": "2"}},
		{map[int]int{0: 0, 1: 1, 2: 0, 3: 2}, map[string]string{"0": "0", "1": "1", "2": "0", "3": "2"}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				got := maps_.MapFunc(tt.data, func(k, v int) (string, string) {
					return strconv.Itoa(k), strconv.Itoa(v)
				})

				if !maps.Equal(got, tt.want) {
					t.Errorf("maps_.MapFunc(%v, func(e int) string {return strconv.Itoa(e)}) = %v, want %v", tt.data, got, tt.want)
				}
			}
		})
	}
}
