// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// it's borrowed from https://github.com/soheilhy/cmux/blob/master/patricia_test.go

package patricia_test

import (
	"github.com/searKing/golang/go/container/patricia"
	"strings"
	"testing"
)

func run(t *testing.T, strs ...string) {
	pt := patricia.NewWithString(strs...)
	for _, s := range strs {
		if !pt.Match(strings.NewReader(s)) {
			t.Errorf("%s is not matched by %s", s, s)
		}

		if !pt.MatchPrefix(strings.NewReader(s + s)) {
			t.Errorf("%s is not matched as a prefix by %s", s+s, s)
		}

		if pt.Match(strings.NewReader(s + s)) {
			t.Errorf("%s matches %s", s+s, s)
		}

		// The following tests are just to catch index out of
		// range and off-by-one errors and not the functionality.
		pt.MatchPrefix(strings.NewReader(s[:len(s)-1]))
		pt.Match(strings.NewReader(s[:len(s)-1]))
		pt.MatchPrefix(strings.NewReader(s + "$"))
		pt.Match(strings.NewReader(s + "$"))
	}
}

func TestPatriciaOnePrefix(t *testing.T) {
	run(t, "prefix")
}

func TestPatriciaNonOverlapping(t *testing.T) {
	run(t, "foo", "bar", "dummy")
}

func TestPatriciaOverlapping(t *testing.T) {
	run(t, "foo", "far", "farther", "boo", "ba", "bar")
}
