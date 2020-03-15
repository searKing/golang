// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binary_test

import (
	"strings"
	"testing"

	"github.com/searKing/golang/go/util/function"
	"github.com/searKing/golang/go/util/function/binary"
)

type BiFunctionTestInput struct {
	apply  func(t interface{}, u interface{}) interface{}
	afters []func(t interface{}) interface{}
	t      interface{}
	u      interface{}
}
type BiFunctionTest struct {
	input  BiFunctionTestInput
	output interface{}
}

var biFunctionTests = []BiFunctionTest{
	{
		input: BiFunctionTestInput{
			apply: func(t interface{}, u interface{}) interface{} {
				return t.(int) + u.(int)
			},
			t: 1,
			u: 2,
		},
		output: 3,
	},
	{
		input: BiFunctionTestInput{
			apply: func(t interface{}, u interface{}) interface{} {
				return t.(string) + u.(string)
			},
			afters: []func(t interface{}) interface{}{func(t interface{}) interface{} {
				return strings.ToUpper(t.(string))
			}, func(t interface{}) interface{} {
				return t.(string) + "c"
			}},
			t: "a",
			u: "b",
		},
		output: "ABc",
	},
}

func TestBiFunction(t *testing.T) {
	for n, test := range biFunctionTests {
		var bi binary.BiFunction = binary.BiFunctionFunc(test.input.apply)
		for _, after := range test.input.afters {
			bi = bi.AndThen(function.FunctionFunc(after))
		}
		got := bi.Apply(test.input.t, test.input.u)
		if got != test.output {
			t.Errorf("#%d: %v: got %v runs; expected %v", n, test.input, got, test.output)
		}
	}
}
