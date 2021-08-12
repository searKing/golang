// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

import (
	"fmt"
	"strings"

	strings_ "github.com/searKing/golang/go/strings"
)

// Placeholders behaves like strings.Join([]string{"?",...,"?"}, ",")
func Placeholders(n int) string {
	return strings_.JoinRepeat("?", ",", n)
}

// Pagination returns the "LIMIT %d, OFFSET %d"
// query := Pagination(0, 0)
// // "LIMIT 0, OFFSET 0"
func Pagination(limit, offset int) string {
	if limit < 0 || offset < 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

// ExpandAsColumns expand columns with alias AS
// query := ExpandAsColumns("table.foo", "bar")
// // []string{"table.foo AS table_foo", "bar AS bar"}
func ExpandAsColumns(cols ...string) []string {
	cols = strings_.SliceTrimEmpty(cols...)
	var params []string
	for _, col := range cols {
		params = append(params, fmt.Sprintf("%[1]s AS %[2]s", col, CompliantName(col)))
	}
	return params
}

// TableColumns returns the []string{table.value1, table.value2 ...}
// query := Columns("table", "foo", "bar")
// // []string{"table.foo", "table.bar"}
func TableColumns(table string, cols ...string) []string {
	cols = strings_.SliceTrimEmpty(cols...)
	var namedCols []string
	for _, col := range cols {
		if table == "" {
			namedCols = append(namedCols, col)
		} else {
			namedCols = append(namedCols, fmt.Sprintf("%s.%s", table, col))
		}
	}
	return namedCols
}

// TableValues returns the []string{:value1, :value2 ...}
// query := TableValues("foo", "bar")
// // []string{"?", "?"}
func TableValues(cols ...string) []string {
	cols = strings_.SliceTrimEmpty(cols...)

	var namedCols []string
	for range cols {
		namedCols = append(namedCols, "?")
	}
	return namedCols
}

// TableColumnsValues returns the []string{table.value1=:value1, table.value2=:value2 ...}
// query := ColumnsValues("table", "foo", "bar")
// // []string{"table.foo=?", "table.bar=?"}
func TableColumnsValues(cmp string, table string, cols ...string) []string {
	cols = strings_.SliceTrimEmpty(cols...)
	var namedCols []string
	for _, col := range cols {
		if table == "" {
			namedCols = append(namedCols, fmt.Sprintf("%[1]s %[2]s ?", col, cmp))
		} else {
			namedCols = append(namedCols, fmt.Sprintf("%[1]s.%[2]s %[3]s ?", table, col, cmp))
		}
	}
	return namedCols
}

// JoinTableColumns concatenates the elements of cols to column1, column2, ...
// query := JoinTableColumns("table", "foo", "bar")
// // "table.foo, table.bar"
func JoinTableColumns(table string, cols ...string) string {
	//cols = strings_.SliceTrimEmpty(cols...)
	return strings.Join(TableColumns(table, cols...), ",")
}

// JoinTableColumnsWithAs concatenates the elements of cols to column1, column2, ...
// query := JoinTableColumnsWithAs("table", "foo", "bar")
// // "table.foo AS table.foo, table.bar AS table.bar"
func JoinTableColumnsWithAs(table string, cols ...string) string {
	//cols = strings_.SliceTrimEmpty(cols...)
	return strings.Join(ExpandAsColumns(TableColumns(table, cols...)...), ",")
}

// JoinColumns concatenates the elements of cols to column1, column2, ...
// query := JoinColumns("foo", "bar")
// // "foo,bar"
func JoinColumns(cols ...string) string {
	return JoinTableColumns("", cols...)
}

// JoinColumnsWithAs concatenates the elements of cols to column1, column2, ...
// query := JoinColumnsWithAs("foo", "bar")
// // "foo AS foo,bar AS bar"
func JoinColumnsWithAs(cols ...string) string {
	return JoinTableColumnsWithAs("", cols...)
}

// JoinTableValues concatenates the elements of values to :value1, :value2, ...
// query := JoinTableValues("foo", "bar")
// // "?,?"
// query := JoinTableValues()
// // ""
func JoinTableValues(cols ...string) string {
	cols = strings_.SliceTrimEmpty(cols...)
	if len(cols) == 0 {
		// https://dev.mysql.com/doc/refman/5.7/en/data-type-defaults.html
		// DEFAULT
		return ""
	}
	return strings.Join(TableValues(cols...), ",")
}

// JoinTableColumnsValues concatenates the elements of values to table.value1=:value1, table.value2=:value2 ...
// query := JoinTableColumnsValues("table", "foo", "bar")
// // "table.foo=?, table.bar=?"
func JoinTableColumnsValues(cmp string, table string, cols ...string) string {
	//cols = strings_.SliceTrimEmpty(cols...)
	return strings.Join(TableColumnsValues(cmp, table, cols...), ",")
}
