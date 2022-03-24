package sqlx

import (
	"fmt"
	"strings"
)

// Pagination returns the "LIMIT %d, OFFSET %d"
// query := Pagination(0, 0)
// // "LIMIT 0, OFFSET 0"
func Pagination(limit, offset int) string {
	if limit < 0 || offset < 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

// TableColumns returns the []string{table.value1, table.value2 ...}
// query := Columns("table", "foo", "bar")
// // []string{"table.foo", "table.bar"}
func TableColumns(table string, cols ...string) []string {
	cols = ShrinkEmptyColumns(cols...)
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
	cols = ShrinkEmptyColumns(cols...)

	var namedCols []string
	for range cols {
		namedCols = append(namedCols, "?")
	}
	return namedCols
}

// ColumnsValues returns the []string{table.value1=:value1, table.value2=:value2 ...}
// query := ColumnsValues("table", "foo", "bar")
// // []string{"table.foo=?", "table.bar=?"}
func TableColumnsValues(cmp SqlCompare, table string, cols ...string) []string {
	cols = ShrinkEmptyColumns(cols...)

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
	//cols = ShrinkEmptyColumns(cols...)
	return strings.Join(TableColumns(table, cols...), ",")
}

// JoinTableColumnsWithAs concatenates the elements of cols to column1, column2, ...
// query := JoinTableColumnsWithAs("table", "foo", "bar")
// // "table.foo AS table.foo, table.bar AS table.bar"
func JoinTableColumnsWithAs(table string, cols ...string) string {
	//cols = ShrinkEmptyColumns(cols...)
	return strings.Join(ExpandAsColumns(TableColumns(table, cols...)...), ",")
}

// JoinTableValues concatenates the elements of values to :value1, :value2, ...
// query := JoinTableValues("foo", "bar")
// // "?,?"
// query := JoinTableValues()
// // "DEFAULT"
func JoinTableValues(cols ...string) string {
	cols = ShrinkEmptyColumns(cols...)
	if len(cols) == 0 {
		// https://dev.mysql.com/doc/refman/5.7/en/data-type-defaults.html
		return "DEFAULT"
	}
	return strings.Join(TableValues(cols...), ",")
}

// JoinTableColumnsValues concatenates the elements of values to table.value1=:value1, table.value2=:value2 ...
// query := JoinTableColumnsValues("table", "foo", "bar")
// // "table.foo=?, table.bar=?"
func JoinTableColumnsValues(table string, cols ...string) string {
	//cols = ShrinkEmptyColumns(cols...)
	return strings.Join(TableColumnsValues(SqlCompareEqual, table, cols...), ",")
}

// JoinTableCondition concatenates the elements of values to table.value1=:value1 AND table.value2=:value2 ...
// query := JoinTableCondition(SqlCompareEqual, SqlOperatorAnd, "table", "foo", "bar")
// // "table.foo=? AND table.bar=?"
func JoinTableCondition(cmp SqlCompare, operator SqlOperator, table string, cols ...string) string {
	//cols = ShrinkEmptyColumns(cols...)
	return strings.Join(TableColumnsValues(cmp, table, cols...), fmt.Sprintf(" %s ", operator.String()))
}
