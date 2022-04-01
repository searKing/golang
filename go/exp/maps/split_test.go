// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps_test

import (
	"testing"

	maps_ "github.com/searKing/golang/go/exp/maps"
)

var equalIntTests = []struct {
	s    map[int]string
	n    int
	want []map[int]string
}{
	{
		nil,
		0,
		nil,
	},
	{
		map[int]string{},
		0,
		nil,
	},
	{
		map[int]string{1: "One"},
		0,
		nil,
	},
	{
		map[int]string{1: "One"},
		1,
		[]map[int]string{{1: "One"}},
	},
	{
		map[int]string{1: "One"},
		2,
		[]map[int]string{{1: "One"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three"},
		-1,
		[]map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three"},
		0,
		nil,
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three"},
		1,
		[]map[int]string{{1: "One", 2: "Two", 3: "Three"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three"},
		2,
		[]map[int]string{{1: "One"}, {2: "Two", 3: "Three"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three"},
		3,
		[]map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three"},
		4,
		[]map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"},
		-1,
		[]map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three"},
		0,
		nil,
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"},
		1,
		[]map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"},
		2,
		[]map[int]string{{1: "One", 4: "Four"}, {2: "Two", 3: "Three"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"},
		3,
		[]map[int]string{{1: "One"}, {2: "Two"}, {3: "Three", 4: "Four"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"},
		4,
		[]map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}},
	},
	{
		map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"},
		5,
		[]map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}},
	},
}

func TestSplitN(t *testing.T) {
	for _, test := range equalIntTests {
		got := maps_.SplitN(test.s, test.n)
		if len(test.want) != len(got) {
			t.Errorf("SplitN(%v, %v) = %v, want %v", test.s, test.n, got, test.want)
			continue
		}
		for i := range got {
			if len(got[i]) != len(test.want[i]) {
				t.Errorf("SplitN(%v, %v) = %v, want %v", test.s, test.n, got, test.want)
				break
			}
		}
	}
}
