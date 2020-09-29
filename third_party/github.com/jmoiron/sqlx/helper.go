package sqlx

import (
	"fmt"
	"strings"
)

//go:generate go-enum -type SqlOperator -trimprefix=SqlOperator --transform=upper
type SqlOperator int

const (
	SqlOperatorAnd SqlOperator = 0
	SqlOperatorOr  SqlOperator = 1
	SqlOperatorNot SqlOperator = 2
)

//go:generate go-enum -type SqlCompare -trimprefix=SqlCompare --linecomment
type SqlCompare int

const (
	SqlCompareEqual            SqlCompare = iota //=
	SqlCompareNotEqual         SqlCompare = iota //<>
	SqlCompareGreaterThan      SqlCompare = iota //>
	SqlCompareLessThan         SqlCompare = iota //<
	SqlCompareGreatEqual       SqlCompare = iota //>=
	SqlCompareLessAndEqualThan SqlCompare = iota //<=
	SqlCompareLike             SqlCompare = iota //LIKE
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

// NamedColumns returns the []string{value1, value2 ...}
// query := NamedColumns("foo", "bar")
// // []string{"foo", "bar"}
func NamedColumns(cols ...string) []string {
	var params = make([]string, len(cols))
	copy(params, cols)
	return params
}

// NamedValues returns the []string{:value1, :value2 ...}
// query := NamedValues("foo", "bar")
// // []string{":foo", ":bar"}
func NamedValues(cols ...string) []string {
	var params = make([]string, len(cols))
	for i, col := range cols {
		params[i] = ":" + col
	}
	return params
}

// NamedColumnsValues returns the []string{value1=:value1, value2=:value2 ...}
// query := NamedColumnsValues("foo", "bar")
// // []string{"foo=:foo", bar=:bar"}
func NamedColumnsValues(cmp SqlCompare, cols ...string) []string {
	var params = make([]string, len(cols))
	for i, col := range cols {
		params[i] = fmt.Sprintf("%[1]s %[2]s :%[1]s", col, cmp)
	}
	return params
}

// JoinColumns concatenates the elements of cols to column1, column2, ...
// query := JoinColumns("foo", "bar")
// // "foo,bar"
func JoinColumns(cols ...string) string {
	return strings.Join(NamedColumns(cols...), ",")
}

// JoinTableColumns concatenates the elements of cols to column1, column2, ...
// query := JoinTableColumns("table", "foo", "bar")
// // "table.foo, table.bar"
func JoinTableColumns(table string, cols ...string) string {
	return strings.Join(NamedTableColumns(table, cols...), ",")
}

// JoinNamedValues concatenates the elements of values to :value1, :value2, ...
// query := JoinNamedValues("foo", "bar")
// // ":foo,:bar"
// query := JoinNamedValues()
// // "DEFAULT"
// Deprecated: Use NamedInsertArguments instead.
func JoinNamedValues(cols ...string) string {
	if len(cols) == 0 {
		// https://dev.mysql.com/doc/refman/5.7/en/data-type-defaults.html
		return "DEFAULT"
	}
	return strings.Join(NamedValues(cols...), ",")
}

// JoinNamedColumnsAndValues concatenates the elements of values to value1=:value1, value2=:value2 ...
// Deprecated: Use NamedUpdateArguments instead.
func JoinNamedColumnsValues(cols ...string) string {
	return strings.Join(NamedColumnsValues(SqlCompareEqual, cols...), ",")
}

// JoinNamedCondition concatenates the elements of values to value1=:value1 AND value2=:value2 ...
// query := JoinNamedCondition(SqlCompareEqual,SqlOperatorAnd,"foo", "bar")
// // "foo=:foo AND bar=:bar"
func JoinNamedCondition(cmp SqlCompare, operator SqlOperator, cols ...string) string {
	return strings.Join(NamedColumnsValues(cmp, cols...), fmt.Sprintf(" %s ", operator.String()))
}

// JoinNamedColumns concatenates the elements of cols to column1, column2, ...
// Deprecated: Use NamedInsertArguments instead.
func JoinNamedColumns(cols ...string) string {
	return JoinColumns(cols...)
}

// JoinNamedTableColumns concatenates the elements of cols in table to column1, column2, ...
func JoinNamedTableColumns(table string, cols ...string) string {
	return JoinTableColumns(table, cols...)
}
