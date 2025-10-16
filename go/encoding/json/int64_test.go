// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json_test

import (
	"encoding/json"
	"testing"

	json_ "github.com/searKing/golang/go/encoding/json"
)

type Int64Test struct {
	input  string
	holder json_.Int64
	output json_.Int64
	want   string
}

var (
	Int64Tests = []Int64Test{
		{
			`0`,
			0,
			0,
			`"0"`,
		},
		{
			`"0"`,
			0,
			0,
			`"0"`,
		},
		{
			`""`,
			0,
			0,
			`"0"`,
		},
		{
			`"-1"`,
			0,
			-1,
			`"-1"`,
		},
		{
			`-1`,
			0,
			-1,
			`"-1"`,
		},
	}
)

func TestInt64(t *testing.T) {
	for n, test := range Int64Tests {
		err := json.Unmarshal([]byte(test.input), &test.holder)
		if err != nil {
			t.Errorf("#%d: Unmarshal error:%v\n", n, err)
			continue
		}

		if test.holder != test.output {
			t.Errorf("#%d: Unmarshal expected %+v got %+v", n, test.output, test.holder)
		}

		got, err := json.Marshal(test.holder)
		if err != nil {
			t.Errorf("#%d: Marshal error:%v\n", n, err)
			continue
		}
		if string(got) != test.want {
			t.Errorf("#%d: Marshal expected %+v got %+v", n, test.want, string(got))
		}
	}
}
