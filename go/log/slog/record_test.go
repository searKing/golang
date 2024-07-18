// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slog

import (
	"testing"
)

func TestShortFunction(t *testing.T) {
	for i, tt := range []struct {
		function string
		want     string
	}{
		{
			"a",
			"a",
		},
		{
			"a.b",
			"b",
		},
		{
			"a.b.",
			"",
		},
		{
			"a.b.c[...]",
			"c[...]",
		},
		{
			"a.b.c[...].[...]....",
			"c[...].[...]....",
		},
	} {

		if got := shortFunction(tt.function); got != tt.want {
			t.Errorf("#%d: got %s\nwant %s", i, got, tt.want)
		}
	}
}

func TestShortFile(t *testing.T) {
	for i, tt := range []struct {
		file string
		want string
	}{
		{
			"a.go",
			"a.go",
		},
		{
			"a/b/c.go",
			"c.go",
		},
		{
			"",
			"???",
		},
		{
			"/",
			"???",
		},
		{
			"a/b/",
			"???",
		},
	} {

		if got := shortFile(tt.file); got != tt.want {
			t.Errorf("#%d: got %s\nwant %s", i, got, tt.want)
		}
	}
}
