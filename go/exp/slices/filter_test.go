// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"testing"

	slices_ "github.com/searKing/golang/go/exp/slices"
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
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				got := slices_.Filter(tt.data)
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("slices_.Filter(%v) = %v, want %v", tt.data, got, tt.want)
					return
				}

				if slices.Compare(got, tt.want) != 0 {
					t.Errorf("slices_.Filter(%v) = %v, want %v", tt.data, got, tt.want)
				}
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
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				copy := slices.Clone(tt.data)
				got := slices_.FilterFunc(copy, func(e int) bool {
					return e != 0
				})
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("slices_.FilterFunc(%v, func(e int) bool {return e != 0}) = %v, want %v", tt.data, got, tt.want)
					return
				}

				if slices.Compare(got, tt.want) != 0 {
					t.Errorf("slices_.FilterFunc(%v, func(e int) bool {return e != 0}) = %v, want %v", tt.data, got, tt.want)
				}
			}
		})
	}
}

func TestTypeAssertFilter(t *testing.T) {
	tests := []struct {
		data []int
		want []int8
	}{
		{nil, nil},
		{[]int{}, []int8{}},
		{[]int{0}, []int8{0}},
		{[]int{1, 0}, []int8{1, 0}},
		{[]int{1, 2}, []int8{1, 2}},
		{[]int{0, 1, 2}, []int8{0, 1, 2}},
		{[]int{0, 1, 0, 2}, []int8{0, 1, 0, 2}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				copy := slices.Clone(tt.data)
				got := slices_.TypeAssertFilter[[]int, []int8](copy)
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("slices_.TypeAssertFilter(%v) = %v, want %v", tt.data, got, tt.want)
					return
				}

				if slices.Compare(got, tt.want) != 0 {
					t.Errorf("slices_.TypeAssertFilter(%v) = %v, want %v", tt.data, got, tt.want)
				}
			}
		})
	}
}

func TestTypeAssertFilterConvert(t *testing.T) {
	var zeroDog Dog
	var zeroMale Male
	_ = zeroDog
	_ = zeroMale
	tests := []struct {
		data []any
		want []Human
	}{
		{nil, nil},
		{[]any{}, []Human{}},
		{[]any{nil}, []Human{nil}},
		{[]any{0, nil}, []Human{nil}},
		{[]any{0, nil, 0, nil}, []Human{nil, nil}},
		{[]any{1, 0, &Dog{}, &Male{}, Dog{}, Male{}}, []Human{&Male{}, Male{}}},
		{[]any{1, 2, zeroDog, zeroMale}, []Human{Male{}}},
		{[]any{zeroMale, zeroDog, 2, &zeroDog, 0, &zeroMale}, []Human{Male{}, &Male{}}},
		{[]any{&zeroMale, &zeroDog, 2, zeroDog, 0, zeroMale}, []Human{&Male{}, Male{}}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				copy := slices.Clone(tt.data)
				got := slices_.TypeAssertFilter[[]any, []Human](copy)
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("slices_.TypeAssertFilterFunc(%v, func(e any) (Human, bool) {h, ok := e.(Human); return h, ok}) = %v, want %v", tt.data, got, tt.want)
					return
				}

				for i, v := range got {
					if i >= len(tt.want) {
						t.Errorf("slices_.TypeAssertFilterFunc(%v, func(e any) (Human, bool) {h, ok := e.(Human); return h, ok}) = %v, want %v", tt.data, got, tt.want)
						return
					}

					if reflect.TypeOf(v) != reflect.TypeOf(tt.want[i]) {
						t.Errorf("slices_.TypeAssertFilterFunc(%v, func(e any) (Human, bool) {h, ok := e.(Human); return h, ok}) =[%d]: %v, want %v", tt.data, i, reflect.TypeOf(v).String(), reflect.TypeOf(tt.want[i]).String())
					}
				}
			}
		})
	}
}

func TestTypeAssertFilterFunc(t *testing.T) {
	tests := []struct {
		data []int
		want []int8
	}{
		{nil, nil},
		{[]int{}, []int8{}},
		{[]int{0}, []int8{0}},
		{[]int{1, 0}, []int8{1, 0}},
		{[]int{1, 2}, []int8{1, 2}},
		{[]int{0, 1, 2}, []int8{0, 1, 2}},
		{[]int{0, 1, 0, 2}, []int8{0, 1, 0, 2}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				copy := slices.Clone(tt.data)
				got := slices_.TypeAssertFilterFunc[[]int, []int8](copy, func(e int) (int8, bool) {
					return int8(e), true
				})
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("slices_.TypeAssertFilterFunc(%v, func(e int) (int8, bool) {return int8(e), true}) = %v, want %v", tt.data, got, tt.want)
					return
				}

				if slices.Compare(got, tt.want) != 0 {
					t.Errorf("slices_.TypeAssertFilterFunc(%v, func(e int) (int8, bool) {return int8(e), true}) = %v, want %v", tt.data, got, tt.want)
					return
				}
			}
		})
	}
}

func TestTypeAssertFilterFuncItoa(t *testing.T) {
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
				copy := slices.Clone(tt.data)
				got := slices_.TypeAssertFilterFunc[[]int, []string](copy, func(e int) (string, bool) {
					return strconv.Itoa(e), true
				})
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("slices_.TypeAssertFilterFunc(%v, func(e int) (string, bool) {return strconv.Itoa(e), true}) = %v, want %v", tt.data, got, tt.want)
					return
				}

				if slices.Compare(got, tt.want) != 0 {
					t.Errorf("slices_.TypeAssertFilterFunc(%v, func(e int) (string, bool) {return strconv.Itoa(e), true}) = %v, want %v", tt.data, got, tt.want)
					return
				}
			}
		})
	}
}

func TestTypeAssertFilterFuncConvert(t *testing.T) {
	var zeroDog Dog
	var zeroMale Male
	_ = zeroDog
	_ = zeroMale
	tests := []struct {
		data []any
		want []Human
	}{
		{nil, nil},
		{[]any{}, []Human{}},
		{[]any{0}, []Human{}},
		{[]any{1, 0, &Dog{}, &Male{}}, []Human{&Male{}}},
		{[]any{1, 2, zeroDog, zeroMale}, []Human{Male{}}},
		{[]any{1, 2, zeroDog, 0, zeroMale, &zeroDog, 0, &zeroMale}, []Human{Male{}, &Male{}}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				copy := slices.Clone(tt.data)
				got := slices_.TypeAssertFilterFunc[[]any, []Human](copy, func(e any) (Human, bool) {
					h, ok := e.(Human)
					return h, ok
				})
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("slices_.TypeAssertFilterFunc(%v, func(e int) (int8, bool) {return int8(e), true}) = %v, want %v", tt.data, got, tt.want)
					return
				}

				for _, v := range got {
					if v.Kind() != "male" {
						t.Errorf("slices_.TypeAssertFilterFunc(%v, func(e any) (Human, bool) {h, ok := e.(Human); return h, ok}) = %v, want %v", tt.data, got, tt.want)
						return
					}
				}
			}
		})
	}
}

type Animal interface {
	Kind() string
}
type Human interface {
	Name() string
	Animal
}

var _ Human = (*Male)(nil)

type Male struct{}

func (m Male) Kind() string {
	return "male"
}

func (Male) Name() string {
	return "bob"
}

var _ Animal = (*Dog)(nil)

type Dog struct{}

func (d Dog) Kind() string {
	return "dog"
}
