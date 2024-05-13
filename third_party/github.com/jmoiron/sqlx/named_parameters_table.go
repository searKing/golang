// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlx

import (
	"fmt"
)

// NamedTableSelectArguments returns columns and arguments for SQL SELECT statements based on columns.
//
//	query := fmt.Sprintf("SELECT %s FROM foo", NamedTableSelectArguments("table", "foo", "bar"))
//	// SELECT table.foo AS table.foo, table.bar AS table.bar FROM table
//
//	query := fmt.Sprintf("SELECT %s FROM table", NamedTableSelectArguments("table", "foo"))
//	// SELECT table.foo AS table.foo FROM table
//
//	query := fmt.Sprintf("SELECT %s FROM table", NamedTableSelectArguments("table"))
//	// SELECT * FROM table
func NamedTableSelectArguments(table string, cols ...string) (arguments string) {
	cols = ShrinkEmptyColumns(cols...)

	if len(cols) == 0 {
		return "*"
	}
	return JoinTableColumnsWithAs(table, cols...)
}

// NamedTableInsertArgumentsCombined returns columns and arguments together
// for SQL INSERT statements based on columns.
//
//	query := fmt.Sprintf("INSERT INTO table %s", NamedTableInsertArgumentsCombined("table", "foo", "bar"))
//	// INSERT INTO table (table.foo, table.bar) VALUES (:table.foo, :table.bar)
//
//	query := fmt.Sprintf("INSERT INTO table %s", NamedTableInsertArgumentsCombined("table"))
//	// INSERT INTO table (foo, bar) DEFAULT VALUES
func NamedTableInsertArgumentsCombined(table string, cols ...string) (arguments string) {
	cols = ShrinkEmptyColumns(cols...)

	//if len(cols) == 0 {
	//	// https://docs.microsoft.com/en-us/sql/t-sql/statements/insert-transact-sql?view=sql-server-ver15#d-inserting-data-into-a-table-with-columns-that-have-default-values
	//	return `DEFAULT VALUES`
	//}
	if len(cols) == 0 {
		return `VALUES(DEFAULT)`
	}
	return fmt.Sprintf(`(%s) VALUES (%s)`, JoinTableColumns(table, cols...), JoinNamedTableValues(table, cols...))
}

// NamedTableInsertArguments returns columns and arguments for SQL INSERT statements based on columns.
//
//	columns, arguments := NamedTableInsertArguments("table", "foo", "bar")
//	query := fmt.Sprintf("INSERT INTO table (%s) VALUES (%s)", columns, arguments)
//	// INSERT INTO table (table.foo, table.bar) VALUES (:table.foo, :table.bar)
//
// Deprecated: Use NamedTableInsertArgumentsCombined instead.
func NamedTableInsertArguments(table string, cols ...string) (columns string, arguments string) {
	return JoinTableColumns(table, cols...), JoinNamedTableValues(table, cols...)
}

// NamedTableUpdateArguments returns columns and arguments for SQL UPDATE statements based on columns.
//
//	statement := NamedTableUpdateArguments("table", "foo", "bar")
//	query := fmt.Sprintf("UPDATE table SET %s", statement)
//	// UPDATE foo SET table.foo=:table.foo, table.bar=:table.bar
func NamedTableUpdateArguments(table string, cols ...string) (arguments string) {
	return JoinNamedTableColumnsValues(table, cols...)
}

// NamedTableWhereArguments returns conditions for SQL WHERE statements based on columns.
//
//	query := fmt.Sprintf("SELECT * FROM table WHERE %s", NamedTableWhereArguments(SqlCompareEqual, SqlOperatorAnd, "table", "foo", "bar"))
//	// SELECT * FROM table WHERE table.foo=:table.foo AND table.bar=:table.bar
//
//	query := fmt.Sprintf("SELECT * FROM table WHERE %s", NamedTableWhereArguments(SqlCompareEqual, SqlOperatorAnd, "table"))
//	// SELECT * FROM table WHERE TRUE
func NamedTableWhereArguments(cmp SqlCompare, operator SqlOperator, table string, cols ...string) (arguments string) {
	cols = ShrinkEmptyColumns(cols...)

	if len(cols) == 0 {
		return "TRUE"
	}
	return JoinNamedTableCondition(cmp, operator, table, cols...)
}
