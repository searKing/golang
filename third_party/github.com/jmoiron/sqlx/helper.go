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

// NamedColumns returns the []string{value1, value2 ...}
// query := NamedColumns("foo", "bar")
// // []string{"foo", "bar"}
func NamedColumns(cols ...string) []string {
	var params []string
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
func NamedColumnsValues(cols ...string) []string {
	var params = make([]string, len(cols))
	for i, col := range cols {
		params[i] = col + "=:" + col
	}
	return params
}

// JoinColumns concatenates the elements of cols to column1, column2, ...
// query := JoinColumns("foo", "bar")
// // "foo,bar"
func JoinColumns(cols ...string) string {
	return strings.Join(NamedColumns(cols...), ",")
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
	return strings.Join(NamedColumnsValues(cols...), ",")
}

// JoinNamedCondition concatenates the elements of values to value1=:value1 AND value2=:value2 ...
// query := JoinNamedCondition(SqlOperatorAnd,"foo", "bar")
// // "foo=:foo AND bar=:bar"
func JoinNamedCondition(operator SqlOperator, cols ...string) string {
	return strings.Join(NamedColumnsValues(cols...), fmt.Sprintf(" %s ", operator.String()))
}

// JoinNamedColumns concatenates the elements of cols to column1, column2, ...
// Deprecated: Use NamedInsertArguments instead.
func JoinNamedColumns(cols []string) string {
	return JoinColumns(cols...)
}
