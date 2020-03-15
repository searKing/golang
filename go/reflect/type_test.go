// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"encoding/json"
	"reflect"
	"testing"
)

type inputType struct {
	a      reflect.Type
	expect string
}

func TestTypeDumpTypeInfoDFS(t *testing.T) {
	var nilError *json.SyntaxError
	ins := []inputType{
		{
			a:      reflect.TypeOf(nil),
			expect: `<nil>`,
		},
		{
			a:      reflect.TypeOf(true),
			expect: `bool`,
		},
		{
			a:      reflect.TypeOf(0),
			expect: `int`,
		},
		{
			a:      reflect.TypeOf(""),
			expect: `string`,
		},
		{
			a: reflect.TypeOf(json.SyntaxError{}),
			expect: `json.SyntaxError
	string
	int64`,
		},
		{
			a: reflect.TypeOf(nilError),
			expect: `*json.SyntaxError
	string
	int64`,
		},
	}
	for idx, in := range ins {
		info := DumpTypeInfoDFS(in.a)
		if info != in.expect {
			t.Errorf("#%d expect\n[\n%s\n]\nactual[\n%s\n]", idx, in.expect, info)
		}
	}
}

func TestTypeDumpTypeInfoBFS(t *testing.T) {
	var nilError *json.SyntaxError
	ins := []inputType{
		{
			a:      reflect.TypeOf(nil),
			expect: `<nil>`,
		},
		{
			a:      reflect.TypeOf(true),
			expect: `bool`,
		},
		{
			a:      reflect.TypeOf(0),
			expect: `int`,
		},
		{
			a:      reflect.TypeOf(""),
			expect: `string`,
		},
		{
			a: reflect.TypeOf(json.SyntaxError{}),
			expect: `json.SyntaxError
	string
	int64`,
		},
		{
			a: reflect.TypeOf(nilError),
			expect: `*json.SyntaxError
	string
	int64`,
		},
	}
	for idx, in := range ins {
		info := DumpTypeInfoBFS(in.a)
		if info != in.expect {
			t.Errorf("#%d expect\n[\n%s\n]\nactual[\n%s\n]", idx, in.expect, info)
		}
	}
}
