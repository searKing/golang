// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps_test

import (
	"fmt"
	"maps"
	"reflect"
	"strconv"
	"testing"

	maps_ "github.com/searKing/golang/go/exp/maps"
)

func TestFilter(t *testing.T) {
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
			{
				got := maps_.Filter(tt.data)
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("maps_.Filter(%v) = %v, want %v", tt.data, got, tt.want)
					return
				}

				if !maps.Equal(got, tt.want) {
					t.Errorf("maps_.Filter(%v) = %v, want %v", tt.data, got, tt.want)
				}
			}
		})
	}
}

func TestFilterFunc(t *testing.T) {
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
			{
				copy := maps.Clone(tt.data)
				got := maps_.FilterFunc(copy, func(k, v int) bool {
					return v != 0
				})
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("maps_.FilterFunc(%v, func(e int) bool {return e != 0}) = %v, want %v", tt.data, got, tt.want)
					return
				}

				if !maps.Equal(got, tt.want) {
					t.Errorf("maps_.FilterFunc(%v, func(e int) bool {return e != 0}) = %v, want %v", tt.data, got, tt.want)
				}
			}
		})
	}
}

func TestTypeAssertFilter(t *testing.T) {
	tests := []struct {
		data map[int]int
		want map[int8]int8
	}{
		{nil, nil},
		{map[int]int{}, map[int8]int8{}},
		{map[int]int{0: 0}, map[int8]int8{0: 0}},
		{map[int]int{3: 0}, map[int8]int8{3: 0}},
		{map[int]int{0: 3}, map[int8]int8{0: 3}},
		{map[int]int{2: 3}, map[int8]int8{2: 3}},
		{map[int]int{1: 1, 0: 0}, map[int8]int8{1: 1, 0: 0}},
		{map[int]int{1: 1, 2: 2}, map[int8]int8{1: 1, 2: 2}},
		{map[int]int{0: 0, 1: 1, 2: 2}, map[int8]int8{0: 0, 1: 1, 2: 2}},
		{map[int]int{0: 0, 1: 1, 3: 0, 2: 2}, map[int8]int8{0: 0, 1: 1, 3: 0, 2: 2}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				copy := maps.Clone(tt.data)
				got := maps_.TypeAssertFilter[map[int]int, map[int8]int8](copy)
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("maps_.TypeAssertFilter(%v) = %v, want %v", tt.data, got, tt.want)
					return
				}

				if !maps.Equal(got, tt.want) {
					t.Errorf("maps_.TypeAssertFilter(%v) = %v, want %v", tt.data, got, tt.want)
				}
			}
		})
	}
}

func TestTypeAssertFilterFunc(t *testing.T) {
	tests := []struct {
		data map[int]int
		want map[int8]int8
	}{
		{nil, nil},
		{map[int]int{}, map[int8]int8{}},
		{map[int]int{0: 0}, map[int8]int8{0: 0}},
		{map[int]int{3: 0}, map[int8]int8{3: 0}},
		{map[int]int{0: 3}, map[int8]int8{0: 3}},
		{map[int]int{2: 3}, map[int8]int8{2: 3}},
		{map[int]int{1: 1, 0: 0}, map[int8]int8{1: 1, 0: 0}},
		{map[int]int{1: 1, 2: 2}, map[int8]int8{1: 1, 2: 2}},
		{map[int]int{0: 0, 1: 1, 2: 2}, map[int8]int8{0: 0, 1: 1, 2: 2}},
		{map[int]int{0: 0, 1: 1, 3: 0, 2: 2}, map[int8]int8{0: 0, 1: 1, 3: 0, 2: 2}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				copy := maps.Clone(tt.data)
				got := maps_.TypeAssertFilterFunc[map[int]int, map[int8]int8](copy, func(k int, v int) (int8, int8, bool) {
					return int8(k), int8(v), true
				})
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("maps_.TypeAssertFilterFunc(%v, func(k int, v int) (int8, int8, bool) {return int8(k), int8(v), true}) = %v, want %v", tt.data, got, tt.want)
					return
				}

				if !maps.Equal(got, tt.want) {
					t.Errorf("maps_.TypeAssertFilterFunc(%v, func(k int, v int) (int8, int8, bool) {return int8(k), int8(v), true}) = %v, want %v", tt.data, got, tt.want)
				}
			}
		})
	}
}

func TestTypeAssertFilterFuncItoa(t *testing.T) {
	tests := []struct {
		data map[int]int
		want map[string]string
	}{
		{nil, nil},
		{map[int]int{}, map[string]string{}},
		{map[int]int{0: 0}, map[string]string{"0": "0"}},
		{map[int]int{3: 0}, map[string]string{"3": "0"}},
		{map[int]int{0: 3}, map[string]string{"0": "3"}},
		{map[int]int{2: 3}, map[string]string{"2": "3"}},
		{map[int]int{1: 1, 0: 0}, map[string]string{"1": "1", "0": "0"}},
		{map[int]int{1: 1, 2: 2}, map[string]string{"1": "1", "2": "2"}},
		{map[int]int{0: 0, 1: 1, 2: 2}, map[string]string{"0": "0", "1": "1", "2": "2"}},
		{map[int]int{0: 0, 1: 1, 3: 0, 2: 2}, map[string]string{"0": "0", "1": "1", "3": "0", "2": "2"}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				copy := maps.Clone(tt.data)
				got := maps_.TypeAssertFilterFunc[map[int]int, map[string]string](copy, func(k int, v int) (string, string, bool) {
					return strconv.Itoa(k), strconv.Itoa(v), true
				})
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("maps_.TypeAssertFilterFunc(%v, func(k int, v int) (int8, int8, bool) {return strconv.Itoa(k), strconv.Itoa(v), true}) = %v, want %v", tt.data, got, tt.want)
					return
				}

				if !maps.Equal(got, tt.want) {
					t.Errorf("maps_.TypeAssertFilterFunc(%v, func(k int, v int) (int8, int8, bool) {return strconv.Itoa(k), strconv.Itoa(v), true}) = %v, want %v", tt.data, got, tt.want)
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
		data map[any]any
		want map[Human]Human
	}{
		{nil, nil},
		{map[any]any{}, map[Human]Human{}},
		{map[any]any{nil: nil}, map[Human]Human{nil: nil}},
		{map[any]any{0: 0, nil: nil}, map[Human]Human{nil: nil}},
		{map[any]any{0: 0, nil: nil, 1: 0, nil: nil}, map[Human]Human{nil: nil, nil: nil}},
		{map[any]any{1: 1, 0: 0, &Dog{}: &Dog{}, &Male{}: &Male{}, Dog{}: Dog{}, Male{}: Male{}}, map[Human]Human{&Male{}: &Male{}, Male{}: Male{}}},
		{map[any]any{1: 1, 2: 2, zeroDog: &zeroDog, zeroMale: &zeroMale}, map[Human]Human{zeroMale: &zeroMale}},
		{map[any]any{1: 1, 2: 2, zeroDog: &zeroMale, zeroDog: &zeroMale}, map[Human]Human{}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				copy := maps.Clone(tt.data)
				got := maps_.TypeAssertFilter[map[any]any, map[Human]Human](copy)
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("maps_.TypeAssertFilterFunc(%v, func(e any) (Human, bool) {h, ok := e.(Human); return h, ok}) = %v, want %v", tt.data, got, tt.want)
					return
				}
				if len(got) != len(tt.want) {
					t.Errorf("maps_.TypeAssertFilterFunc(%v, func(e any) (Human, bool) {h, ok := e.(Human); return h, ok}) = %v, want %v", tt.data, got, tt.want)
					return
				}

				for k, v := range got {
					if reflect.TypeOf(v) != reflect.TypeOf(tt.want[k]) {
						t.Errorf("maps_.TypeAssertFilterFunc(%v, func(e any) (Human, bool) {h, ok := e.(Human); return h, ok}) =[%d]: %v, want %v", tt.data, k, reflect.TypeOf(v).String(), reflect.TypeOf(tt.want[k]).String())
					}
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
		data map[any]any
		want map[Human]Human
	}{
		{nil, nil},
		{map[any]any{}, map[Human]Human{}},
		{map[any]any{nil: nil}, map[Human]Human{nil: nil}},
		{map[any]any{0: 0, nil: nil}, map[Human]Human{nil: nil}},
		{map[any]any{0: 0, nil: nil, 1: 0, nil: nil}, map[Human]Human{nil: nil, nil: nil}},
		{map[any]any{1: 1, 0: 0, &Dog{}: &Dog{}, &Male{}: &Male{}, Dog{}: Dog{}, Male{}: Male{}}, map[Human]Human{&Male{}: &Male{}, Male{}: Male{}}},
		{map[any]any{1: 1, 2: 2, zeroDog: &zeroDog, zeroMale: &zeroMale}, map[Human]Human{zeroMale: &zeroMale}},
		{map[any]any{1: 1, 2: 2, zeroDog: &zeroMale, zeroDog: &zeroMale}, map[Human]Human{}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			{
				copy := maps.Clone(tt.data)
				got := maps_.TypeAssertFilterFunc[map[any]any, map[Human]Human](copy, func(k any, v any) (Human, Human, bool) {
					k2, ok := k.(Human)
					v2, ov := v.(Human)
					return k2, v2, (ok || k == nil) && (ov || v == nil)
				})
				if (got == nil && tt.data != nil) || (got != nil && tt.data == nil) {
					t.Errorf("maps_.TypeAssertFilterFunc(%v, func(e any) (Human, bool) {h, ok := e.(Human); return h, ok}) = %v, want %v", tt.data, got, tt.want)
					return
				}
				if len(got) != len(tt.want) {
					t.Errorf("maps_.TypeAssertFilterFunc(%v, func(e any) (Human, bool) {h, ok := e.(Human); return h, ok}) = %v, want %v", tt.data, got, tt.want)
					return
				}

				for k, v := range got {
					if reflect.TypeOf(v) != reflect.TypeOf(tt.want[k]) {
						t.Errorf("maps_.TypeAssertFilterFunc(%v, func(e any) (Human, bool) {h, ok := e.(Human); return h, ok}) =[%d]: %v, want %v", tt.data, k, reflect.TypeOf(v).String(), reflect.TypeOf(tt.want[k]).String())
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
