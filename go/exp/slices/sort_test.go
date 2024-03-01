// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
	"testing"

	math_ "github.com/searKing/golang/go/exp/math"
	slices_ "github.com/searKing/golang/go/exp/slices"
)

func TestLinearSearch(t *testing.T) {
	str1 := []string{"foo"}
	str2 := []string{"ab", "ca"}
	str3 := []string{"mo", "qo", "vo"}
	str4 := []string{"ab", "ad", "ca", "xy"}

	// slice with repeating elements
	strRepeats := []string{"ba", "ca", "da", "da", "da", "ka", "ma", "ma", "ta"}

	// slice with all element equal
	strSame := []string{"xx", "xx", "xx"}

	tests := []struct {
		data      []string
		target    string
		wantPos   int
		wantFound bool
	}{
		{[]string{}, "foo", 0, false},
		{[]string{}, "", 0, false},

		{str1, "foo", 0, true},
		{str1, "bar", 0, false},
		{str1, "zx", 1, false},

		{str2, "aa", 0, false},
		{str2, "ab", 0, true},
		{str2, "ad", 1, false},
		{str2, "ca", 1, true},
		{str2, "ra", 2, false},

		{str3, "bb", 0, false},
		{str3, "mo", 0, true},
		{str3, "nb", 1, false},
		{str3, "qo", 1, true},
		{str3, "tr", 2, false},
		{str3, "vo", 2, true},
		{str3, "xr", 3, false},

		{str4, "aa", 0, false},
		{str4, "ab", 0, true},
		{str4, "ac", 1, false},
		{str4, "ad", 1, true},
		{str4, "ax", 2, false},
		{str4, "ca", 2, true},
		{str4, "cc", 3, false},
		{str4, "dd", 3, false},
		{str4, "xy", 3, true},
		{str4, "zz", 4, false},

		{strRepeats, "da", 2, true},
		{strRepeats, "db", 5, false},
		{strRepeats, "ma", 6, true},
		{strRepeats, "mb", 8, false},

		{strSame, "xx", 0, true},
		{strSame, "ab", 0, false},
		{strSame, "zz", 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			{
				pos, found := slices_.LinearSearch(tt.data, tt.target)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("LinearSearch got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
				wantPos, wantFound := slices.BinarySearch(tt.data, tt.target)
				if pos != wantPos || found != wantFound {
					t.Errorf("LinearSearch got (%v, %v), BinarySearch want (%v, %v)", pos, found, wantPos, wantFound)
				}
			}

			{
				pos, found := slices_.LinearSearchFunc(tt.data, tt.target, strings.Compare)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("LinearSearchFunc got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
				wantPos, wantFound := slices.BinarySearchFunc(tt.data, tt.target, strings.Compare)
				if pos != wantPos || found != wantFound {
					t.Errorf("LinearSearch got (%v, %v), BinarySearchFunc want (%v, %v)", pos, found, wantPos, wantFound)
				}
			}
		})
	}
}

func TestLinearSearchInts(t *testing.T) {
	tests := []struct {
		data      []int
		target    int
		wantPos   int
		wantFound bool
	}{
		{nil, 20, 0, false},
		{[]int{}, 20, 0, false},
		{[]int{20, 20, 30, 30}, 20, 0, true},
		{[]int{20, 20, 30, 30}, 23, 2, false},
		{[]int{20, 20, 30, 30}, 30, 2, true},
		{[]int{20, 20, 30, 30}, 43, 4, false},
		{[]int{20, 30, 40, 50, 60, 70, 80, 90}, 20, 0, true},
		{[]int{20, 30, 40, 50, 60, 70, 80, 90}, 23, 1, false},
		{[]int{20, 30, 40, 50, 60, 70, 80, 90}, 43, 3, false},
		{[]int{20, 30, 40, 50, 60, 70, 80, 90}, 80, 6, true},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.target), func(t *testing.T) {
			{
				pos, found := slices_.LinearSearch(tt.data, tt.target)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("LinearSearch got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
				wantPos, wantFound := slices.BinarySearch(tt.data, tt.target)
				if pos != wantPos || found != wantFound {
					t.Errorf("LinearSearch got (%v, %v), BinarySearch want (%v, %v)", pos, found, wantPos, wantFound)
				}
			}

			{
				compare := cmp.Compare[int]
				pos, found := slices_.LinearSearchFunc(tt.data, tt.target, compare)
				if pos != tt.wantPos || found != tt.wantFound {
					t.Errorf("LinearSearchFunc got (%v, %v), want (%v, %v)", pos, found, tt.wantPos, tt.wantFound)
				}
				wantPos, wantFound := slices.BinarySearchFunc(tt.data, tt.target, compare)
				if pos != wantPos || found != wantFound {
					t.Errorf("LinearSearch got (%v, %v), BinarySearchFunc want (%v, %v)", pos, found, wantPos, wantFound)
				}
			}
		})
	}
}

func TestPartialSortInts(t *testing.T) {
	tests := []struct {
		data []int
		k    int
	}{
		{nil, 3},
		{[]int{}, 3},
		{[]int{20, 20, 30, 30}, 3},
		{[]int{20, 30}, 3},
		{[]int{20, 30, 40, 50, 60, 70, 80, 90}, 3},
		{[]int{90, 80, 70, 60, 50, 40, 30, 20}, 3},
		{[]int{90, 30, 70, 40, 50, 60, 80, 20}, 3},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			{
				data1 := slices.Clone(tt.data)
				slices_.PartialSort(data1, tt.k)
				if !slices_.IsPartialSorted(data1, tt.k) {
					t.Errorf("partial sort didn't sort")
				}
			}

			{
				data2 := slices.Clone(tt.data)
				slices_.PartialSortFunc(data2, tt.k, cmp.Compare[int])
				if !slices_.IsPartialSorted(data2, tt.k) {
					t.Errorf("partial sort func didn't sort")
				}
			}
		})
	}
}

func TestSearchMin(t *testing.T) {
	tests := []struct {
		data []int
		want int
	}{
		{nil, 0},
		{[]int{}, 0},
		{[]int{20, 20, 30, 30}, 0},
		{[]int{20, 30}, 0},
		{[]int{20, 30, 40, 50, 60, 70, 80, 90}, 0},
		{[]int{90, 80, 70, 60, 50, 40, 30, 20}, 7},
		{[]int{90, 30, 70, 40, 50, 60, 80, 20}, 7},
		{[]int{90, 30, 70, 40, 50, 60, 80, 20, 100, 30}, 7},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			{
				got := slices_.SearchMin(tt.data)
				if got != tt.want {
					t.Errorf("SearchMin(%v) got (%v), want (%v)", tt.data, got, tt.want)
				}
			}
			{
				got := slices_.SearchMinFunc(tt.data, cmp.Compare[int])
				if got != tt.want {
					t.Errorf("SearchMinFunc(%v, math_.Compare[int]) got (%v), want (%v)", tt.data, got, tt.want)
				}
			}
		})
	}
}

func TestSearchMax(t *testing.T) {
	tests := []struct {
		data []int
		want int
	}{
		{nil, 0},
		{[]int{}, 0},
		{[]int{20, 20, 30, 30}, 2},
		{[]int{20, 30}, 1},
		{[]int{20, 30, 40, 50, 60, 70, 80, 90}, 7},
		{[]int{90, 80, 70, 60, 50, 40, 30, 20}, 0},
		{[]int{90, 30, 70, 40, 50, 60, 80, 20}, 0},
		{[]int{90, 30, 70, 40, 50, 60, 80, 20, 100, 30}, 8},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			{
				got := slices_.SearchMax(tt.data)
				if got != tt.want {
					t.Errorf("SearchMax(%v) got (%v), want (%v)", tt.data, got, tt.want)
				}
			}
			{
				got := slices_.SearchMinFunc(tt.data, math_.Reverse(cmp.Compare[int]))
				if got != tt.want {
					t.Errorf("SearchMinFunc(%v, math_.Reverse(math_.Compare[int])) got (%v), want (%v)", tt.data, got, tt.want)
				}
			}
		})
	}
}
