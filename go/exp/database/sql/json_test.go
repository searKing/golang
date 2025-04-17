// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql_test

import (
	"reflect"
	"testing"

	"github.com/searKing/golang/go/exp/database/sql"
)

func TestJsonByteSlice(t *testing.T) {
	for n, test := range []struct {
		input  any
		holder sql.Json[[]byte]
		output sql.Json[[]byte]
	}{
		{
			nil,
			sql.Json[[]byte]{},
			sql.Json[[]byte]{Data: nil},
		},
		{
			"[]",
			sql.Json[[]byte]{},
			sql.Json[[]byte]{Data: []byte{}},
		},
		{
			"[1,2,3]",
			sql.Json[[]byte]{},
			sql.Json[[]byte]{Data: []byte{1, 2, 3}},
		},
	} {
		err := test.holder.Scan(test.input)
		if err != nil {
			t.Errorf("#%d: Scan error:%v\n", n, err)
			continue
		}

		if !reflect.DeepEqual(test.holder, test.output) {
			t.Errorf("#%d: expected %+v got %+v", n, test.output, test.holder)
		}
	}
}

func TestJsonStruct(t *testing.T) {
	for n, test := range []struct {
		input  any
		holder sql.Json[people]
		output sql.Json[people]
	}{
		{
			nil,
			sql.Json[people]{},
			sql.Json[people]{Data: people{}},
		},
		{
			"",
			sql.Json[people]{},
			sql.Json[people]{Data: people{}},
		},
		{
			`{"Name":"Alice","Age":100,"Friends":[{"Name":"Bob","Age":50,"Friends":null}]}`,
			sql.Json[people]{},
			sql.Json[people]{
				Data: people{
					Name: "Alice",
					Age:  100,
					Friends: []people{{
						Name:    "Bob",
						Age:     50,
						Friends: nil,
					}},
				},
			},
		},
	} {
		err := test.holder.Scan(test.input)
		if err != nil {
			t.Errorf("#%d: Scan error:%v\n", n, err)
			continue
		}

		if !reflect.DeepEqual(test.holder, test.output) {
			t.Errorf("#%d: expected %+v got %+v", n, test.output, test.holder)
		}
	}
}
