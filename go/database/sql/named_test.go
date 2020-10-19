// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql_test

import (
	"testing"

	"github.com/searKing/golang/go/database/sql"
)

func TestCompileQuery(t *testing.T) {
	table := []struct {
		Q, R string
	}{
		{
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?), (?, ?, ?, ?) ON DUPLICATE KEY UPDATE a=?, b=?`,
			R: `insert into foo(a, b, c, d) values (:a, :b, :c, :d), (:a, :b, :c, :d) on duplicate key update a = :a, b = :b`,
		},
		{
			Q: `UPDATE foo SET foo=?, bar=? WHERE thud=? AND grunt=?`,
			R: `update foo set foo = :foo, bar = :bar where thud = :thud and grunt = :grunt`,
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
		qr, err := sql.CompileQuery(test.Q)
		if err != nil {
			t.Errorf("#%d. got err %s, want err nil", i, err)
		}
		if qr != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}

func TestWithCompileQueryOptionAliasWithSelect(t *testing.T) {
	table := []struct {
		Q, R string
		A    bool
	}{
		{
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?), (?, ?, :c, :d)`,
			R: `insert into foo(a, b, c, d) values (:a, :b, :c, :d), (:a, :b, :c, :d)`,
			A: false,
		},
		{
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`,
			R: `insert into foo(a, b, c, d) values (:a, :b, :c, :d)`,
			A: true,
		},
		{
			Q: `UPDATE foo SET foo=?, bar=? WHERE thud=? AND grunt=?`,
			R: `update foo set foo = :foo, bar = :bar where thud = :thud and grunt = :grunt`,
			A: false,
		},
		{
			Q: `UPDATE foo SET foo=?, bar=? WHERE thud=? AND grunt=?`,
			R: `update foo set foo = :foo, bar = :bar where thud = :thud and grunt = :grunt`,
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
		qr, err := sql.CompileQuery(test.Q, sql.WithCompileQueryOptionAliasWithSelect(test.A))
		if err != nil {
			t.Errorf("#%d. got err %s, want err nil", i, err)
		}
		if qr != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}
func TestWithCompileQueryOptionArgument(t *testing.T) {
	table := []struct {
		Q, R string
		T    interface{}
	}{
		{ // 0
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?), (?, ?, ?, ?)`,
			R: `insert into foo(a, b, c, d) values (:a, :b, :c, :d), (:a, :b, :c, :d)`,
			T: nil,
		},
		{ // 1
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`,
			R: `insert into foo(a, b, d) values (:a, :b, :d)`,
			T: []string{"a", "b", "d"},
		},
		{ // 2
			Q: `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?), (?, ?, ?, ?)`,
			R: `insert into foo(a, b, c, d) values (:a, :b, :c, :d), (:a, :b, :c, :d)`,
			T: []string{"c"},
		},
		{ // 3
			Q: `UPDATE foo SET foo = ?, bar = ? WHERE thud = ? AND grunt = ?`,
			R: `update foo set foo = :foo, bar = :bar where thud = :thud and grunt = :grunt`,
			T: nil,
		},
		{ // 4
			Q: `UPDATE foo SET foo = ?, bar =? WHERE thud = ? AND grunt = ?`,
			R: `update foo set bar = :bar where grunt = :grunt`,
			T: []string{"bar", "grunt"},
		},
		{ // 5
			Q: `UPDATE foo SET foo = ?, bar =? WHERE thud = ? AND grunt = ?`,
			R: `update foo set bar = :bar where grunt = :grunt`,
			T: map[string]bool{"foo": false, "bar": true, "thud": false, "grunt": true},
		},
		{ // 6
			Q: `UPDATE foo SET foo = ?, bar =? WHERE thud = ? AND grunt = ?`,
			R: `update foo set bar = :bar where grunt = :grunt`,
			T: struct {
				Foo   bool `db:"foo"`
				Bar   bool `db:"bar"`
				Thud  bool `db:"thud"`
				Grunt bool `db:"grunt"`
			}{
				Foo:   false,
				Bar:   true,
				Thud:  false,
				Grunt: true,
			},
		},
		{ // 7
			Q: `SELECT t.a, b FROM t WHERE first_name =:hehe AND middle_name = ? OR last_name = ?`,
			R: `select t.a as t_a, b as b from t where first_name = :first_name and middle_name = :middle_name or last_name = :last_name`,
			T: nil,
		},
		{ // 8
			Q: `SELECT t.a, b FROM t WHERE first_name =:hehe AND middle_name = ? OR last_name = ?`,
			R: `select t.a as t_a, b as b from t where first_name = :first_name or last_name = :last_name`,
			T: []string{"first_name", "last_name"},
		},
		{ // 9
			Q: `SELECT ":foo" FROM a WHERE first_name = 1 AND last_name = 'NAME'`,
			R: `select ':foo' as __foo_ from a where first_name = 1 and last_name = 'NAME'`,
			T: nil,
		},
		{ // 10
			Q: `SELECT ":foo" FROM a WHERE first_name = 1 AND last_name = 'NAME'`,
			R: `select ':foo' as __foo_ from a where first_name = 1 and last_name = 'NAME'`,
			T: []string{"first_name", "last_name"},
		},
		{ // 11
			Q: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name =:first_name and middle_name = :middle_name or last_name =:last_name`,
			R: `select 'a:b:c' or first_name as _a_b_c__or_first_name, '::ABC:_:' as ___ABC____ from person where first_name = :first_name and middle_name = :middle_name or last_name = :last_name`,
			T: nil,
		},
		{ // 12
			Q: `SELECT 'a:b:c' || first_name, '::ABC:_:' FROM person WHERE first_name =:first_name and middle_name =:middle_name or last_name =:last_name`,
			R: `select 'a:b:c' or first_name as _a_b_c__or_first_name, '::ABC:_:' as ___ABC____ from person where first_name = :first_name or last_name = :last_name`,
			T: []string{"first_name", "last_name"},
		},
	}

	for i, test := range table {
		qr, err := sql.CompileQuery(test.Q, sql.WithCompileQueryOptionArgument(test.T))
		if err != nil {
			t.Errorf("#%d. got err %s, want err nil", i, err)
		}
		if qr != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
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
		qr := sql.CompliantName(test.Q)
		if qr != test.R {
			t.Errorf("#%d. got %q, want %q", i, qr, test.R)
		}
	}
}
