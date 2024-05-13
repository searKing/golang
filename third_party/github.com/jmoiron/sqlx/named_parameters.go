// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlx

// NamedSelectArguments returns columns and arguments for SQL SELECT statements based on columns.
//
//	query := fmt.Sprintf("SELECT %s FROM foo", NamedSelectArguments("foo", "bar"))
//	// SELECT foo, bar FROM foo
//
//	query := fmt.Sprintf("SELECT %s FROM foo", NamedSelectArguments("foo"))
//	// SELECT foo FROM foo
//
//	query := fmt.Sprintf("SELECT %s FROM foo", NamedSelectArguments())
//	// SELECT * FROM foo
func NamedSelectArguments(cols ...string) (arguments string) {
	return NamedTableSelectArguments("", cols...)
}

// NamedSelectArgumentsWithAs returns columns and arguments for SQL SELECT statements based on columns with alias.
//
//	query := fmt.Sprintf("SELECT %s FROM foo", NamedSelectArgumentsWithAs("foo", "bar"))
//	// SELECT foo AS foo, bar AS bar FROM foo
//
//	query := fmt.Sprintf("SELECT %s FROM foo", NamedSelectArgumentsWithAs("foo"))
//	// SELECT foo AS foo FROM foo
//
//	query := fmt.Sprintf("SELECT %s FROM foo", NamedSelectArgumentsWithAs())
//	// SELECT * FROM foo
func NamedSelectArgumentsWithAs(cols ...string) (arguments string) {
	return NamedTableSelectArguments("", cols...)
}

// NamedInsertArgumentsCombined returns columns and arguments together
// for SQL INSERT statements based on columns.
//
//	query := fmt.Sprintf("INSERT INTO foo %s", NamedInsertArgumentsCombined("foo", "bar"))
//	// INSERT INTO foo (foo, bar) VALUES (:foo, :bar)
//
//	query := fmt.Sprintf("INSERT INTO foo %s", NamedInsertArgumentsCombined())
//	// INSERT INTO foo (foo, bar) DEFAULT VALUES
func NamedInsertArgumentsCombined(cols ...string) (arguments string) {
	return NamedTableInsertArgumentsCombined("", cols...)
}

// NamedInsertArguments returns columns and arguments for SQL INSERT statements based on columns.
//
//	columns, arguments := NamedInsertArguments("foo", "bar")
//	query := fmt.Sprintf("INSERT INTO foo (%s) VALUES (%s)", columns, arguments)
//	// INSERT INTO foo (foo, bar) VALUES (:foo, :bar)
//
// Deprecated: Use NamedInsertArgumentsCombined instead.
func NamedInsertArguments(cols ...string) (columns string, arguments string) {
	return NamedTableInsertArguments("", cols...)
}

// NamedUpdateArguments returns columns and arguments for SQL UPDATE statements based on columns.
//
//	statement := NamedUpdateArguments("foo", "bar")
//	query := fmt.Sprintf("UPDATE foo SET %s", statement)
//	// UPDATE foo SET foo=:foo, bar=:bar
func NamedUpdateArguments(cols ...string) (arguments string) {
	return NamedTableUpdateArguments("", cols...)
}

// NamedWhereArguments returns conditions for SQL WHERE statements based on columns.
//
//	query := fmt.Sprintf("SELECT * FROM foo WHERE %s", NamedWhereArguments(SqlCompareEqual, SqlOperatorAnd, "foo", "bar"))
//	// SELECT * FROM foo WHERE foo=:foo AND bar=:bar
//
//	query := fmt.Sprintf("SELECT * FROM foo WHERE %s", NamedWhereArguments(SqlCompareEqual, SqlOperatorAnd))
//	// SELECT * FROM foo WHERE TRUE
func NamedWhereArguments(cmp SqlCompare, operator SqlOperator, cols ...string) (arguments string) {
	return NamedTableWhereArguments(cmp, operator, "", cols...)
}
