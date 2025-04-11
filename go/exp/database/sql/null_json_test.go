// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql_test

import (
	"reflect"
	"testing"

	"github.com/searKing/golang/go/exp/database/sql"
)

type people struct {
	Name    string
	Age     int
	Friends []people
}

func TestNullJsonByteSlice(t *testing.T) {
	for n, test := range []struct {
		input  any
		holder sql.NullJson[[]byte]
		output sql.NullJson[[]byte]
	}{
		{
			nil,
			sql.NullJson[[]byte]{},
			sql.NullJson[[]byte]{Data: nil, Valid: false},
		},
		{
			"",
			sql.NullJson[[]byte]{},
			sql.NullJson[[]byte]{Data: nil, Valid: true},
		},
		{
			"[]",
			sql.NullJson[[]byte]{},
			sql.NullJson[[]byte]{Data: []byte{}, Valid: true},
		},
		{
			"[1,2,3]",
			sql.NullJson[[]byte]{},
			sql.NullJson[[]byte]{Data: []byte{1, 2, 3}, Valid: true},
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

func TestNullJsonStruct(t *testing.T) {
	for n, test := range []struct {
		input  any
		holder sql.NullJson[people]
		output sql.NullJson[people]
	}{
		{
			nil,
			sql.NullJson[people]{},
			sql.NullJson[people]{Data: people{}, Valid: false},
		},
		{
			"",
			sql.NullJson[people]{},
			sql.NullJson[people]{Data: people{}, Valid: true},
		},
		{
			`{"Name":"Alice","Age":100,"Friends":[{"Name":"Bob","Age":50,"Friends":null}]}`,
			sql.NullJson[people]{},
			sql.NullJson[people]{
				Data: people{
					Name: "Alice",
					Age:  100,
					Friends: []people{{
						Name:    "Bob",
						Age:     50,
						Friends: nil,
					}},
				},
				Valid: true,
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
