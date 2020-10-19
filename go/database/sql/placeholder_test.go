// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql_test

import (
	"testing"

	"github.com/searKing/golang/go/database/sql"
	"github.com/searKing/golang/go/strings"
)

func TestExpandAsColumns(t *testing.T) {
	table := []struct {
		Q, R []string
	}{
		{
			Q: []string{`a`},
			R: []string{`a AS a`},
		},
		{
			Q: []string{`a.b`},
			R: []string{`a.b AS a_b`},
		},
		{
			Q: []string{`a.b1asdas`},
			R: []string{`a.b1asdas AS a_b1asdas`},
		},
	}
	for i, test := range table {
		qr := sql.ExpandAsColumns(test.Q...)
		if !strings.SliceEqual(qr, test.R) {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}
