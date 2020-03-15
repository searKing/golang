// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binary_test

import (
	"testing"

	"github.com/searKing/golang/go/util"
	"github.com/searKing/golang/go/util/function"
	"github.com/searKing/golang/go/util/function/binary"
)

type MinByTestInput struct {
	compare func(t interface{}, u interface{}) int
	afters  []func(t interface{}) interface{}
	t       interface{}
	u       interface{}
}
type MinByTest struct {
	input  MinByTestInput
	output interface{}
}

var minByTests = []MinByTest{
	{
		input: MinByTestInput{
			compare: func(t interface{}, u interface{}) int {
				return t.(int) - u.(int)
			},
			t: 1,
			u: 2,
		},
		output: 1,
	},
	{
		input: MinByTestInput{
			compare: func(t interface{}, u interface{}) int {
				return t.(int) - u.(int)
			},
			afters: []func(t interface{}) interface{}{func(t interface{}) interface{} {
				return 0 - t.(int)
			}, func(t interface{}) interface{} {
				return t.(int) * 3
			}},
			t: 1,
			u: 2,
		},
		output: -3,
	},
}

func TestMinBy(t *testing.T) {
	for n, test := range minByTests {
		var bi = binary.MinBy(util.ComparatorFunc(test.input.compare))
		for _, after := range test.input.afters {
			bi = bi.AndThen(function.FunctionFunc(after))
		}
		got := bi.Apply(test.input.t, test.input.u)
		if got != test.output {
			t.Errorf("#%d: %v: got %v runs; expected %v", n, test.input, got, test.output)
		}
	}
}

type MaxByTestInput struct {
	compare func(t interface{}, u interface{}) int
	afters  []func(t interface{}) interface{}
	t       interface{}
	u       interface{}
}
type MaxByTest struct {
	input  MaxByTestInput
	output interface{}
}

var maxByTests = []MaxByTest{
	{
		input: MaxByTestInput{
			compare: func(t interface{}, u interface{}) int {
				return t.(int) - u.(int)
			},
			t: 1,
			u: 2,
		},
		output: 2,
	},
	{
		input: MaxByTestInput{
			compare: func(t interface{}, u interface{}) int {
				return t.(int) - u.(int)
			},
			afters: []func(t interface{}) interface{}{func(t interface{}) interface{} {
				return 0 - t.(int)
			}, func(t interface{}) interface{} {
				return t.(int) * 3
			}},
			t: 1,
			u: 2,
		},
		output: -6,
	},
}

func TestMaxBy(t *testing.T) {
	for n, test := range maxByTests {
		var bi = binary.MaxBy(util.ComparatorFunc(test.input.compare))
		for _, after := range test.input.afters {
			bi = bi.AndThen(function.FunctionFunc(after))
		}
		got := bi.Apply(test.input.t, test.input.u)
		if got != test.output {
			t.Errorf("#%d: %v: got %v runs; expected %v", n, test.input, got, test.output)
		}
	}
}
