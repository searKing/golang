// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql_test

import (
	"testing"

	"github.com/searKing/golang/go/database/sql"
)

func TestCompliantName(t *testing.T) {
	table := []struct {
		Q, R string
	}{
		{
			Q: `a`,
			R: `a`,
		},
		{
			Q: `t.a`,
			R: `t_a`,
		},
		{
			Q: `':foo'`,
			R: `__foo_`,
		},
		{
			Q: `'a:b:c' || first_name`,
			R: `_a_b_c_____first_name`,
		},
		{
			Q: `'::ABC:_:'`,
			R: `___ABC____`,
		},
	}

	for i, test := range table {
		qr := sql.CompliantName(test.Q)
		if qr != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}
