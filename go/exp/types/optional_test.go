// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/searKing/golang/go/exp/types"
)

func ExampleOptional_String() {
	fmt.Printf("%s\n", types.Optional[string]{}.String())
	fmt.Printf("%s\n", types.Optional[int]{}.String())
	fmt.Printf("%s\n", types.Optional[*int]{}.String())

	fmt.Printf("%s\n", types.Opt("").String())
	fmt.Printf("%s\n", types.Opt(0).String())
	fmt.Printf("%s\n", types.Opt[*int](nil).String())

	fmt.Printf("%s\n", types.Opt("hello world").String())
	fmt.Printf("%s\n", types.Opt(0xFF).String())
	var p int
	fmt.Printf("%s\n", types.Opt(&p).String()[:2])

	// Output:
	// null
	// null
	// null
	//
	// 0
	// <nil>
	// hello world
	// 255
	// 0x
}

func ExampleOptional_Format() {
	fmt.Printf("'%+v'\n", types.Optional[string]{})
	fmt.Printf("'%v'\n", types.Optional[string]{})
	fmt.Printf("'%s'\n", types.Optional[string]{})
	fmt.Printf("'%+v'\n", types.Optional[int]{})
	fmt.Printf("'%v'\n", types.Optional[int]{})
	fmt.Printf("'%s'\n", types.Optional[int]{})
	fmt.Printf("'%+v'\n", types.Optional[*int]{})
	fmt.Printf("'%v'\n", types.Optional[*int]{})
	fmt.Printf("'%s'\n", types.Optional[*int]{})

	fmt.Printf("'%+v'\n", types.Opt(""))
	fmt.Printf("'%v'\n", types.Opt(""))
	fmt.Printf("'%s'\n", types.Opt(""))
	fmt.Printf("'%+v'\n", types.Opt(0))
	fmt.Printf("'%v'\n", types.Opt(0))
	fmt.Printf("'%s'\n", types.Opt(0))
	fmt.Printf("'%+v'\n", types.Opt[*int](nil))
	fmt.Printf("'%v'\n", types.Opt[*int](nil))
	fmt.Printf("'%s'\n", types.Opt[*int](nil))

	fmt.Printf("'%+v'\n", types.Opt("hello world"))
	fmt.Printf("'%v'\n", types.Opt("hello world"))
	fmt.Printf("'%s'\n", types.Opt("hello world"))
	fmt.Printf("'%+v'\n", types.Opt(0xFF))
	fmt.Printf("'%v'\n", types.Opt(0xFF))
	fmt.Printf("'%s'\n", types.Opt(0xFF))
	var p int
	fmt.Printf("'%s'\n", fmt.Sprintf("%+v", types.Opt(&p))[:2])
	fmt.Printf("'%s'\n", fmt.Sprintf("%v", types.Opt(&p))[:2])
	fmt.Printf("'%s'\n", fmt.Sprintf("%s", types.Opt(&p))[:11])

	// Output:
	// 'null string: '
	// 'null'
	// 'null'
	// 'null int: 0'
	// 'null'
	// 'null'
	// 'null *int: <nil>'
	// 'null'
	// 'null'
	// ''
	// ''
	// ''
	// '0'
	// '0'
	// '%!s(int=0)'
	// '<nil>'
	// '<nil>'
	// '%!s(*int=<nil>)'
	// 'hello world'
	// 'hello world'
	// 'hello world'
	// '255'
	// '255'
	// '%!s(int=255)'
	// '0x'
	// '0x'
	// '%!s(*int=0x'
}

func TestOptional_String(t *testing.T) {
	tests := []struct {
		a    types.Optional[int]
		want string
	}{
		{types.Optional[int]{}, "null"},
		{types.Optional[int]{}, "null"},
		{types.Opt[int](1), "1"},
		{types.Opt[int](0), "0"},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			{
				got := tt.a.String()
				if got != tt.want {
					t.Errorf("%v.String() got (%v), want (%v)", tt.a, got, tt.want)
				}
			}
		})
	}

}

func TestCompareOptional(t *testing.T) {
	tests := []struct {
		a, b types.Optional[int]
		want int
	}{
		{types.Optional[int]{}, types.Optional[int]{}, 0},
		{types.Optional[int]{}, types.Opt[int](1), -1},
		{types.Opt[int](1), types.Optional[int]{}, 1},
		{types.Opt[int](1), types.Opt[int](1), 0},
		{types.Opt[int](0), types.Opt[int](1), -1},
		{types.Opt[int](1), types.Opt[int](0), 1},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			{
				got := types.CompareOptional(tt.a, tt.b)
				if got != tt.want {
					t.Errorf("CompareOptional(%v, %v) got (%v), want (%v)", tt.a, tt.b, got, tt.want)
				}
			}
		})
	}
}
