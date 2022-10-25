// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps_test

import (
	"testing"

	maps_ "github.com/searKing/golang/go/exp/maps"
)

var splitTests = []struct {
	m    map[int]string
	sep  int
	want []map[int]string
}{
	{nil, -1, nil},
	{nil, 0, nil},
	{nil, 1, nil},
	{map[int]string{}, -1, nil},
	{map[int]string{}, 0, nil},
	{map[int]string{}, 1, nil},
	{map[int]string{1: "One"}, -1, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 0, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 1, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 2, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, -1, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 0, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 1, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 2, []map[int]string{{1: "One", 2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 3, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 4, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, -1, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 0, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 1, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 2, []map[int]string{{1: "One", 2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 3, []map[int]string{{1: "One", 2: "Two", 3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 4, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 5, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
}

func TestSplit(t *testing.T) {
	for i, test := range splitTests {
		got := maps_.Split(test.m, test.sep)
		if len(test.want) != len(got) {
			t.Errorf("#%d: Split(%v, %v) = %v, want %v", i, test.m, test.sep, got, test.want)
			continue
		}
		for i := range got {
			if len(got[i]) != len(test.want[i]) {
				t.Errorf("#%d: Split(%v, %v) = %v, want %v", i, test.m, test.sep, got, test.want)
				break
			}
		}
	}
}

var splitNTests = []struct {
	m    map[int]string
	sep  int
	n    int
	want []map[int]string
}{
	{nil, -1, -1, nil},
	{nil, -1, 0, nil},
	{nil, -1, 1, nil},
	{nil, 0, -1, nil},
	{nil, 0, 0, nil},
	{nil, 0, 1, nil},
	{nil, 1, -1, nil},
	{nil, 1, 0, nil},
	{nil, 1, 1, nil},
	{map[int]string{}, -1, -1, nil},
	{map[int]string{}, -1, 0, nil},
	{map[int]string{}, -1, 1, []map[int]string{{}}},
	{map[int]string{}, 0, -1, nil},
	{map[int]string{}, 0, 0, nil},
	{map[int]string{}, 0, 1, nil},
	{map[int]string{}, 1, -1, nil},
	{map[int]string{}, 1, 0, nil},
	{map[int]string{}, 1, 1, []map[int]string{{}}},
	{map[int]string{1: "One"}, -1, -1, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, -1, 0, nil},
	{map[int]string{1: "One"}, -1, 1, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, -1, 2, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 0, -1, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 0, 0, nil},
	{map[int]string{1: "One"}, 0, 1, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 0, 2, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 1, -1, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 1, 0, nil},
	{map[int]string{1: "One"}, 1, 1, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 1, 2, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 2, -1, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 2, 0, nil},
	{map[int]string{1: "One"}, 2, 1, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One"}, 2, 2, []map[int]string{{1: "One"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, -1, -1, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, -1, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, -1, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, -1, 2, []map[int]string{{1: "One"}, {2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, -1, 3, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, -1, 4, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 0, -1, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 0, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 0, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 0, 2, []map[int]string{{1: "One"}, {2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 0, 3, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 0, 4, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 1, -1, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 1, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 1, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 1, 2, []map[int]string{{1: "One"}, {2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 1, 3, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 1, 4, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 2, -1, []map[int]string{{1: "One", 2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 2, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 2, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 2, 2, []map[int]string{{1: "One", 2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 2, 3, []map[int]string{{1: "One", 2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 2, 4, []map[int]string{{1: "One", 2: "Two"}, {3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 3, -1, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 3, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 3, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 3, 2, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 3, 3, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 3, 4, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 4, -1, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 4, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 4, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 4, 2, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 4, 3, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three"}, 4, 4, []map[int]string{{1: "One", 2: "Two", 3: "Three"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, -1, -1, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, -1, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, -1, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, -1, 2, []map[int]string{{1: "One", 2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, -1, 3, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, -1, 4, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, -1, 5, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 0, -1, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 0, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 0, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 0, 2, []map[int]string{{1: "One", 2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 0, 3, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 0, 4, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 0, 5, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 1, -1, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 1, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 1, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 1, 2, []map[int]string{{1: "One"}, {2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 1, 3, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 1, 4, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 1, 5, []map[int]string{{1: "One"}, {2: "Two"}, {3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 2, -1, []map[int]string{{1: "One", 2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 2, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 2, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 2, 2, []map[int]string{{1: "One", 2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 2, 3, []map[int]string{{1: "One", 2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 2, 4, []map[int]string{{1: "One", 2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 2, 5, []map[int]string{{1: "One", 2: "Two"}, {3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 3, -1, []map[int]string{{1: "One", 2: "Two", 3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 3, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 3, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 3, 2, []map[int]string{{1: "One", 2: "Two", 3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 3, 3, []map[int]string{{1: "One", 2: "Two", 3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 3, 4, []map[int]string{{1: "One", 2: "Two", 3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 3, 5, []map[int]string{{1: "One", 2: "Two", 3: "Three"}, {4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 4, -1, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 4, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 4, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 4, 2, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 4, 3, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 4, 4, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 4, 5, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 5, -1, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 5, 0, nil},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 5, 1, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 5, 2, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 5, 3, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 5, 4, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
	{map[int]string{1: "One", 2: "Two", 3: "Three", 4: "Four"}, 5, 5, []map[int]string{{1: "One", 2: "Two", 3: "Three", 4: "Four"}}},
}

func TestSplitN(t *testing.T) {
	for i, test := range splitNTests {
		got := maps_.SplitN(test.m, test.sep, test.n)
		if len(test.want) != len(got) {
			t.Errorf("#%d: SplitN(%v, %v, %v) = %v, want %v", i, test.m, test.sep, test.n, got, test.want)
			continue
		}
		for i := range got {
			if len(got[i]) != len(test.want[i]) {
				t.Errorf("#%d: SplitN(%v, %v, %v) = %v, want %v", i, test.m, test.sep, test.n, got, test.want)
				break
			}
		}
	}
}
