package sqlx

import (
	"fmt"
	"strings"
)

// NamedTableColumns returns the []string{table.value1, table.value2 ...}
// query := NamedColumns("table", "foo", "bar")
// // []string{"table.foo", "table.bar"}
func NamedTableColumns(table string, cols ...string) []string {
	var params = make([]string, len(cols))
	copy(params, cols)
	if table == "" {
		return params
	}
	for _, param := range params {
		param = fmt.Sprintf("%s.%s", table, param)
	}
	return params
}

// NamedTableValues returns the []string{:value1, :value2 ...}
// query := NamedTableValues("foo", "bar")
// // []string{":foo", ":bar"}
func NamedTableValues(cols ...string) []string {
	var params = make([]string, len(cols))
	for i, col := range cols {
		params[i] = ":" + col
	}
	return params
}

// NamedColumnsValues returns the []string{value1=:value1, value2=:value2 ...}
// query := NamedColumnsValues("foo", "bar")
// // []string{"foo=:foo", bar=:bar"}
func NamedTableColumnsValues(cmp SqlCompare, table string, cols ...string) []string {
	var params = make([]string, len(cols))
	for i, col := range cols {
		if table == "" {
			params[i] = fmt.Sprintf("%[1]s %[2]s :%[1]s", col, cmp)
		} else {
			params[i] = fmt.Sprintf("%[1]s.%[2]s %[3]s :%[2]s", table, col, cmp)
		}
	}
	return params
}

// JoinTableColumns concatenates the elements of cols to column1, column2, ...
// query := JoinTableColumns("table", "foo", "bar")
// // "table.foo, table.bar"
func JoinTableColumns(table string, cols ...string) string {
	return strings.Join(NamedTableColumns(table, cols...), ",")
}

// JoinNamedTableValues concatenates the elements of values to :value1, :value2, ...
// query := JoinNamedTableValues("foo", "bar")
// // ":foo,:bar"
// query := JoinNamedTableValues()
// // "DEFAULT"
func JoinNamedTableValues(cols ...string) string {
	if len(cols) == 0 {
		// https://dev.mysql.com/doc/refman/5.7/en/data-type-defaults.html
		return "DEFAULT"
	}
	return strings.Join(NamedTableValues(cols...), ",")
}

// JoinNamedTableColumnsValues concatenates the elements of values to table.value1=:value1, table.value2=:value2 ...
func JoinNamedTableColumnsValues(table string, cols ...string) string {
	return strings.Join(NamedTableColumnsValues(SqlCompareEqual, table, cols...), ",")
}

// JoinNamedTableCondition concatenates the elements of values to table.value1=:value1 AND table.value2=:value2 ...
// query := JoinNamedTableCondition(SqlCompareEqual,SqlOperatorAnd, "table", "foo", "bar")
// // "table.foo=:foo AND table.bar=:bar"
func JoinNamedTableCondition(cmp SqlCompare, operator SqlOperator, table string, cols ...string) string {
	return strings.Join(NamedTableColumnsValues(cmp, table, cols...), fmt.Sprintf(" %s ", operator.String()))
}

// JoinNamedTableColumns concatenates the elements of cols in table to column1, column2, ...
func JoinNamedTableColumns(table string, cols ...string) string {
	return JoinTableColumns(table, cols...)
}
