package sqlx

import (
	"strings"
)

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
	if len(cols) == 0 {
		return "*"
	}
	return joinColumns(cols)
}

// NamedInsertArguments returns columns and arguments for SQL INSERT statements based on columns.
//
//	columns, arguments := NamedInsertArguments("foo", "bar")
//	query := fmt.Sprintf("INSERT INTO foo (%s) VALUES (%s)", columns, arguments)
//	// INSERT INTO foo (foo, bar) VALUES (:foo, :bar)
func NamedInsertArguments(cols ...string) (columns string, arguments string) {
	return joinColumns(cols), joinNamedValues(cols)
}

// NamedUpdateArguments returns columns and arguments for SQL UPDATE statements based on columns.
//
//	statement := NamedUpdateArguments("foo", "bar")
//	query := fmt.Sprintf("UPDATE foo SET %s", statement)
//	// UPDATE foo SET foo=:foo, bar=:bar
func NamedUpdateArguments(cols ...string) (arguments string) {
	return joinNamedColumnsAndValues(cols)
}

// NamedWhereArguments returns conditions for SQL WHERE statements based on columns.
//
//	query := fmt.Sprintf("SELECT * FROM foo WHERE %s", NamedWhereArguments("foo", "bar"))
//	// SELECT * FROM foo WHERE foo=:foo AND bar=:bar
//
//	query := fmt.Sprintf("SELECT * FROM foo WHERE %s", NamedWhereArguments())
//	// SELECT * FROM foo WHERE TRUE
func NamedWhereArguments(cols ...string) (arguments string) {
	if len(cols) == 0 {
		return "TRUE"
	}
	return joinNamedCondition(cols)
}

// joinColumns concatenates the elements of cols to column1, column2, ...
func joinColumns(cols []string) string {
	return strings.Join(cols, ",")
}

// joinNamedValues concatenates the elements of values to :value1, :value2, ...
func joinNamedValues(cols []string) string {
	if len(cols) == 0 {
		return ""
	}
	return ":" + strings.Join(cols, ", :")
}

// JoinNamedColumns concatenates the elements of cols to column1, column2, ...
// Deprecated: Use NamedInsertArguments instead.
func JoinNamedColumns(cols []string) string {
	return joinColumns(cols)
}

// joinNamedColumnsAndValues concatenates the elements of values to value1=:value1, value2=:value2 ...
func joinNamedColumnsAndValues(cols []string) string {
	if len(cols) == 0 {
		return ""
	}

	for i, col := range cols {
		cols[i] = col + "=:" + col
	}
	return strings.Join(cols, ", ")
}

// joinNamedCondition concatenates the elements of values to value1=:value1 AND value2=:value2 ...
func joinNamedCondition(cols []string) string {
	for i, col := range cols {
		cols[i] = col + "=:" + col
	}
	return strings.Join(cols, " AND ")
}

// JoinNamedValues concatenates the elements of values to :value1, :value2, ...
// Deprecated: Use NamedInsertArguments instead.
func JoinNamedValues(cols []string) string {
	return joinNamedValues(cols)
}

// JoinNamedColumnsAndValues concatenates the elements of values to value1=:value1, value2=:value2 ...
// Deprecated: Use NamedUpdateArguments instead.
func JoinNamedColumnsAndValues(cols []string) string {
	return joinNamedColumnsAndValues(cols)
}

// JoinNamedCondition concatenates the elements of values to value1=:value1 AND value2=:value2 ...
// Deprecated: Use NamedWhereArguments instead.
func JoinNamedCondition(cols []string) string {
	return joinNamedCondition(cols)
}
