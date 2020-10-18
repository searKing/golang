// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlx_test

import (
	"testing"

	sqlx_ "github.com/searKing/golang/third_party/github.com/jmoiron/sqlx"
)

func TestCompileQuery(t *testing.T) {
	table := []struct {
		Q, R string
	}{
		{
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`,
			R: `insert into foo(a, b, c, d) values (:a, :b, :c, :d)`,
		},
		{
			Q: `SELECT t.a, b FROM t WHERE first_name= :hehe AND middle_name=? OR last_name=?`,
			R: `select t.a as t_a, b as b from t where first_name = :first_name and middle_name = :middle_name or last_name = :last_name`,
		},
		{
			Q: `SELECT ":foo" FROM a WHERE first_name=1 AND last_name='NAME'`,
			R: `select ':foo' as __foo_ from a where first_name = 1 and last_name = 'NAME'`,
		},
		{
			Q: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name=? AND last_name=?`,
			R: `select 'a:b:c' or first_name as _a_b_c__or_first_name, '::ABC:_:' as ___ABC____ from person where first_name = :first_name and last_name = :last_name`,
		},
	}

	for i, test := range table {
		qr, err := sqlx_.CompileQuery(test.Q)
		if err != nil {
			t.Errorf("%d. got err %s, want err nil", i, err)
		}
		if qr != test.R {
			t.Errorf("%d. got %q, want %q", i, qr, test.R)
		}
	}
}

func TestWithCompileQueryOptionAliasWithSelect(t *testing.T) {
	table := []struct {
		Q, R string
		A    bool
	}{
		{
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`,
			R: `insert into foo(a, b, c, d) values (:a, :b, :c, :d)`,
			A: false,
		},
		{
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`,
			R: `insert into foo(a, b, c, d) values (:a, :b, :c, :d)`,
			A: true,
		},
		{
			Q: `SELECT t.a, b FROM t WHERE first_name= :hehe AND middle_name=? OR last_name=?`,
			R: `select t.a, b from t where first_name = :first_name and middle_name = :middle_name or last_name = :last_name`,
			A: false,
		},
		{
			Q: `SELECT t.a, b FROM t WHERE first_name= :hehe AND middle_name=? OR last_name=?`,
			R: `select t.a as t_a, b as b from t where first_name = :first_name and middle_name = :middle_name or last_name = :last_name`,
			A: true,
		},
		{
			Q: `SELECT ":foo" FROM a WHERE first_name=1 AND last_name='NAME'`,
			R: `select ':foo' from a where first_name = 1 and last_name = 'NAME'`,
			A: false,
		},
		{
			Q: `SELECT ":foo" FROM a WHERE first_name=1 AND last_name='NAME'`,
			R: `select ':foo' as __foo_ from a where first_name = 1 and last_name = 'NAME'`,
			A: true,
		},
		{
			Q: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name=? AND last_name=?`,
			R: `select 'a:b:c' or first_name, '::ABC:_:' from person where first_name = :first_name and last_name = :last_name`,
			A: false,
		},
		{
			Q: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name=? AND last_name=?`,
			R: `select 'a:b:c' or first_name as _a_b_c__or_first_name, '::ABC:_:' as ___ABC____ from person where first_name = :first_name and last_name = :last_name`,
			A: true,
		},
	}

	for i, test := range table {
		qr, err := sqlx_.CompileQuery(test.Q, sqlx_.WithCompileQueryOptionAliasWithSelect(test.A))
		if err != nil {
			t.Errorf("%d. got err %s, want err nil", i, err)
		}
		if qr != test.R {
			t.Errorf("%d. got %q, want %q", i, qr, test.R)
		}
	}
}
func TestWithCompileQueryOptionTrimByColumn(t *testing.T) {
	table := []struct {
		Q, R string
		T    map[string]interface{}
	}{
		{ // 0
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`,
			R: `insert into foo(a, b, c, d) values (:a, :b, :c, :d)`,
			T: nil,
		},
		{ // 1
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`,
			R: `insert into foo(a, b, c, d) values (:a, :b, :c, :d)`,
			T: map[string]interface{}{"b": nil},
		},
		{ // 2
			Q: `SELECT t.a, b FROM t WHERE first_name= :hehe AND middle_name=? OR last_name=?`,
			R: `select t.a as t_a, b as b from t where first_name = :first_name and middle_name = :middle_name or last_name = :last_name`,
			T: nil,
		},
		{ // 3
			Q: `SELECT t.a, b FROM t WHERE first_name= :hehe AND middle_name=? OR last_name=?`,
			R: `select t.a as t_a, b as b from t where first_name = :first_name or last_name = :last_name`,
			T: map[string]interface{}{"middle_name": nil},
		},
		{ // 4
			Q: `SELECT ":foo" FROM a WHERE first_name=1 AND last_name='NAME'`,
			R: `select ':foo' as __foo_ from a where first_name = 1 and last_name = 'NAME'`,
			T: nil,
		},
		{ // 5
			Q: `SELECT ":foo" FROM a WHERE first_name=1 AND last_name='NAME'`,
			R: `select ':foo' as __foo_ from a where first_name = 1 and last_name = 'NAME'`,
			T: map[string]interface{}{"last_name": nil},
		},
		{ // 6
			Q: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name = :first_name and middle_name = :middle_name or last_name = :last_name`,
			R: `select 'a:b:c' or first_name as _a_b_c__or_first_name, '::ABC:_:' as ___ABC____ from person where first_name = :first_name and middle_name = :middle_name or last_name = :last_name`,
			T: nil,
		},
		{ // 7
			Q: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name = :first_name and middle_name = :middle_name or last_name = :last_name`,
			R: `select 'a:b:c' or first_name as _a_b_c__or_first_name, '::ABC:_:' as ___ABC____ from person where first_name = :first_name and middle_name = :middle_name`,
			T: map[string]interface{}{"last_name": nil},
		},
	}

	for i, test := range table {
		qr, err := sqlx_.CompileQuery(test.Q, sqlx_.WithCompileQueryOptionTrimByColumn(test.T))
		if err != nil {
			t.Errorf("%d. got err %s, want err nil", i, err)
		}
		if qr != test.R {
			t.Errorf("%d. got %q, want %q", i, qr, test.R)
		}
	}
}

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
		qr := sqlx_.CompliantName(test.Q)
		if qr != test.R {
			t.Errorf("%d. got %q, want %q", i, qr, test.R)
		}
	}
}
