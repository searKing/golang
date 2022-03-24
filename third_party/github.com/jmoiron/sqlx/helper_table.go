// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlx

import (
	"fmt"
	"strings"
)

// ShrinkEmptyColumns trim empty columns
func ShrinkEmptyColumns(cols ...string) []string {
	var params []string
	for _, col := range cols {
		if col == "" {
			continue
		}
		params = append(params, col)
	}
	return params
}

// ExpandAsColumns expand columns with alias AS
// query := ExpandAsColumns("foo", "bar")
// // []string{"foo AS foo", "bar AS bar"}
func ExpandAsColumns(cols ...string) []string {
	cols = ShrinkEmptyColumns(cols...)

	var params []string
	for _, col := range cols {
		params = append(params, fmt.Sprintf("%[1]s AS %[2]s", col, strings.ReplaceAll(col, ".", "_")))
	}
	return params
}

// NamedTableColumns returns the []string{table.value1, table.value2 ...}
// query := NamedColumns("table", "foo", "bar")
// // []string{"table.foo", "table.bar"}
func NamedTableColumns(table string, cols ...string) []string {
	cols = ShrinkEmptyColumns(cols...)
	return TableColumns(table, cols...)
}

// NamedTableValues returns the []string{:value1, :value2 ...}
// query := NamedTableValues("table", "foo", "bar")
// // []string{":table.foo", ":table.bar"}
func NamedTableValues(table string, cols ...string) []string {
	cols = ShrinkEmptyColumns(cols...)

	var namedCols []string
	for _, col := range cols {
		if table == "" {
			namedCols = append(namedCols, fmt.Sprintf(":%[1]s", col))
		} else {
			namedCols = append(namedCols, fmt.Sprintf(":%[1]s.%[2]s", table, col))
		}
	}
	return namedCols
}

// NamedColumnsValues returns the []string{table.value1=:value1, table.value2=:value2 ...}
// query := NamedColumnsValues("table", "foo", "bar")
// // []string{"table.foo=:table.foo", "table.bar=:table.bar"}
func NamedTableColumnsValues(cmp SqlCompare, table string, cols ...string) []string {
	cols = ShrinkEmptyColumns(cols...)

	var namedCols []string
	for _, col := range cols {
		if table == "" {
			namedCols = append(namedCols, fmt.Sprintf("%[1]s %[2]s :%[1]s", col, cmp))
		} else {
			namedCols = append(namedCols, fmt.Sprintf("%[1]s.%[2]s %[3]s :%[1]s.%[2]s", table, col, cmp))
		}
	}
	return namedCols
}

// JoinNamedTableValues concatenates the elements of values to :value1, :value2, ...
// query := JoinNamedTableValues("table", "foo", "bar")
// // ":table.foo,:table.bar"
// query := JoinNamedTableValues("table")
// // "DEFAULT"
func JoinNamedTableValues(table string, cols ...string) string {
	cols = ShrinkEmptyColumns(cols...)
	//if len(cols) == 0 {
	//	// https://dev.mysql.com/doc/refman/5.7/en/data-type-defaults.html
	//	return "DEFAULT"
	//}
	return strings.Join(NamedTableValues(table, cols...), ",")
}

// JoinNamedTableColumnsValues concatenates the elements of values to table.value1=:value1, table.value2=:value2 ...
// query := JoinNamedTableColumnsValues("table", "foo", "bar")
// // "table.foo=:table.foo, table.bar=:table.bar"
func JoinNamedTableColumnsValues(table string, cols ...string) string {
	//cols = ShrinkEmptyColumns(cols...)
	return strings.Join(NamedTableColumnsValues(SqlCompareEqual, table, cols...), ",")
}

// JoinNamedTableCondition concatenates the elements of values to table.value1=:value1 AND table.value2=:value2 ...
// query := JoinNamedTableCondition(SqlCompareEqual, SqlOperatorAnd, "table", "foo", "bar")
// // "table.foo=:table.foo AND table.bar=:table.bar"
func JoinNamedTableCondition(cmp SqlCompare, operator SqlOperator, table string, cols ...string) string {
	//cols = ShrinkEmptyColumns(cols...)
	return strings.Join(NamedTableColumnsValues(cmp, table, cols...), fmt.Sprintf(" %s ", operator.String()))
}

// JoinNamedTableColumns concatenates the elements of cols in table to column1, column2, ...
// query := JoinNamedTableColumns("table", "foo", "bar")
// // "table.foo, table.bar"
func JoinNamedTableColumns(table string, cols ...string) string {
	//cols = ShrinkEmptyColumns(cols...)
	return JoinTableColumns(table, cols...)
}

func JoinNamedTableColumnsWithAs(table string, cols ...string) string {
	//cols = ShrinkEmptyColumns(cols...)
	return JoinTableColumnsWithAs(table, cols...)
}
