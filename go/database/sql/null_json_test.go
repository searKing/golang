// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql_test

import (
	"reflect"
	"testing"

	"github.com/searKing/golang/go/database/sql"
)

type people struct {
	Name    string
	Age     int
	Friends []people
}

type NullJsonTest struct {
	input  string
	holder sql.NullJson
	output sql.NullJson
}

var (
	nullJsonTests = []NullJsonTest{
		{
			`[1,2,3]`,
			sql.NullJson{
				Data: &[]byte{},
			},
			sql.NullJson{Data: &[]byte{1, 2, 3},
				Valid: true,
			},
		},
		{
			`{"Name":"Alice","Age":100,"Friends":[{"Name":"Bob","Age":50,"Friends":null}]}`,
			sql.NullJson{
				Data: &people{},
			},
			sql.NullJson{Data: &people{Name: "Alice",
				Age: 100,
				Friends: []people{{
					Name:    "Bob",
					Age:     50,
					Friends: nil,
				}}},
				Valid: true,
			},
		},
	}
)

func TestNullJson(t *testing.T) {

	for n, test := range nullJsonTests {
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
