// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type inputValue struct {
	a      reflect.Value
	expect bool
}

func TestIsEmptyValue(t *testing.T) {
	ins := []inputValue{
		{
			a:      reflect.ValueOf(nil),
			expect: true,
		},
		{
			a:      reflect.ValueOf(true),
			expect: false,
		},
		{
			a:      reflect.ValueOf(0),
			expect: true,
		},
		{
			a:      reflect.ValueOf(""),
			expect: true,
		},
		{
			a:      reflect.ValueOf(time.Now()),
			expect: false,
		},
		{
			a:      reflect.ValueOf(time.Time{}),
			expect: true,
		},
		{
			a:      reflect.ValueOf([]byte{}),
			expect: true,
		},
		{
			a:      reflect.ValueOf([]byte(nil)),
			expect: true,
		},
		{
			a:      reflect.ValueOf(map[int]string{}),
			expect: true,
		},
		{
			a:      reflect.ValueOf(map[int]string(nil)),
			expect: true,
		},
		{
			a:      reflect.ValueOf(any([]byte{})),
			expect: true,
		},
		{
			a:      reflect.ValueOf(any([]byte(nil))),
			expect: true,
		},
		{
			a:      reflect.ValueOf(any(map[int]string{})),
			expect: true,
		},
		{
			a:      reflect.ValueOf(any(map[int]string(nil))),
			expect: true,
		},
		{
			a:      reflect.ValueOf(struct{}{}),
			expect: true,
		},
	}
	for idx, in := range ins {
		if IsEmptyValue(in.a) != in.expect {
			t.Errorf("#%d expect %t", idx, in.expect)
		}
	}
}

func TestIsZeroValue(t *testing.T) {
	ins := []inputValue{
		{
			a:      reflect.ValueOf(nil),
			expect: true,
		},
		{
			a:      reflect.ValueOf(true),
			expect: false,
		},
		{
			a:      reflect.ValueOf(0),
			expect: true,
		},
		{
			a:      reflect.ValueOf(1),
			expect: false,
		},
		{
			a:      reflect.ValueOf(""),
			expect: true,
		},
		{
			a:      reflect.ValueOf(time.Now()),
			expect: false,
		},
		{
			a:      reflect.ValueOf(time.Time{}),
			expect: true,
		},
		{
			a:      reflect.ValueOf([]byte{}),
			expect: false,
		},
		{
			a:      reflect.ValueOf([]byte(nil)),
			expect: true,
		},
		{
			a:      reflect.ValueOf(map[int]string{}),
			expect: false,
		},
		{
			a:      reflect.ValueOf(map[int]string(nil)),
			expect: true,
		},
		{
			a:      reflect.ValueOf(any([]byte{})),
			expect: false,
		},
		{
			a:      reflect.ValueOf(any([]byte(nil))),
			expect: true,
		},
		{
			a:      reflect.ValueOf(any(map[int]string{})),
			expect: false,
		},
		{
			a:      reflect.ValueOf(any(map[int]string(nil))),
			expect: true,
		},
		{
			a:      reflect.ValueOf(struct{}{}),
			expect: true,
		},
	}
	for idx, in := range ins {
		if IsZeroValue(in.a) != in.expect {
			t.Errorf("#%d expect %t", idx, in.expect)
		}
	}
}

func TestIsNilValue(t *testing.T) {
	var nilTime *time.Time
	ins := []inputValue{
		{
			a:      reflect.ValueOf(nil),
			expect: true,
		},
		{
			a:      reflect.ValueOf(true),
			expect: false,
		},
		{
			a:      reflect.ValueOf(0),
			expect: false,
		},
		{
			a:      reflect.ValueOf(""),
			expect: false,
		},
		{
			a:      reflect.ValueOf(time.Now()),
			expect: false,
		},
		{
			a:      reflect.ValueOf(nilTime), // typed nil
			expect: true,
		},
	}
	for idx, in := range ins {
		if IsNilValue(in.a) != in.expect {
			t.Errorf("#%d expect %t", idx, in.expect)
		}
	}
}

type inputDumpValue struct {
	a      reflect.Value
	expect string
}

func TestTypeDumpValueInfoDFS(t *testing.T) {
	var nilError *json.SyntaxError
	ins := []inputDumpValue{
		{
			a:      reflect.ValueOf(nil),
			expect: `<invalid value>`,
		},
		{
			a:      reflect.ValueOf(true),
			expect: `[bool: true]`,
		},
		{
			a:      reflect.ValueOf(0),
			expect: `[int: 0]`,
		},
		{
			a:      reflect.ValueOf("HelloWorld"),
			expect: `[string: HelloWorld]`,
		},
		{
			a: reflect.ValueOf(json.SyntaxError{}),
			expect: `[json.SyntaxError: {msg: Offset:0}]
	[string: ]
	[int64: 0]`,
		},
		{
			a:      reflect.ValueOf(nilError),
			expect: `[*json.SyntaxError: <nil>]`,
		},
	}
	for idx, in := range ins {
		info := DumpValueInfoDFS(in.a)
		if info != in.expect {
			t.Errorf("#%d expect\n[\n%s\n]\nactual\n[\n%s\n]", idx, in.expect, info)
		}
	}
}

func TestTypeDumpValueInfoBFS(t *testing.T) {
	var nilError *json.SyntaxError
	ins := []inputDumpValue{
		{
			a:      reflect.ValueOf(nil),
			expect: `<invalid value>`,
		},
		{
			a:      reflect.ValueOf(true),
			expect: `[bool: true]`,
		},
		{
			a:      reflect.ValueOf(0),
			expect: `[int: 0]`,
		},
		{
			a:      reflect.ValueOf(""),
			expect: `[string: ]`,
		},
		{
			a: reflect.ValueOf(json.SyntaxError{}),
			expect: `[json.SyntaxError: {msg: Offset:0}]
	[string: ]
	[int64: 0]`,
		},
		{
			a:      reflect.ValueOf(nilError),
			expect: `[*json.SyntaxError: <nil>]`,
		},
	}
	for idx, in := range ins {
		info := DumpValueInfoBFS(in.a)
		if info != in.expect {
			t.Errorf("#%d expect\n[\n%s\n]\nactual\n[\n%s\n]", idx, in.expect, info)
		}
	}
}

func TestDumpValueInfoBFS(t *testing.T) {
	str := "HelloWorld"
	s := &str
	ss := &s
	valueS := reflect.ValueOf(ss)
	indirectValueS := reflect.Indirect(valueS)
	fmt.Printf("valueS: %s\n", valueS.Kind().String())
	fmt.Printf("indirect valueS: %s\n", indirectValueS.Kind().String())
}
