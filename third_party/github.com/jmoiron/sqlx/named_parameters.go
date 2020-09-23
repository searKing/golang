package sqlx

import (
	"fmt"
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
	return JoinColumns(cols...)
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
	//if len(cols) == 0 {
	//	// https://docs.microsoft.com/en-us/sql/t-sql/statements/insert-transact-sql?view=sql-server-ver15#d-inserting-data-into-a-table-with-columns-that-have-default-values
	//	return `DEFAULT VALUES`
	//}
	if len(cols) == 0 {
		return `VALUES(DEFAULT)`
	}
	return fmt.Sprintf(`(%s) VALUES (%s)`, JoinColumns(cols...), JoinNamedValues(cols...))
}

// NamedInsertArguments returns columns and arguments for SQL INSERT statements based on columns.
//
//	columns, arguments := NamedInsertArguments("foo", "bar")
//	query := fmt.Sprintf("INSERT INTO foo (%s) VALUES (%s)", columns, arguments)
//	// INSERT INTO foo (foo, bar) VALUES (:foo, :bar)
// Deprecated: Use NamedInsertArgumentsCombined instead.
func NamedInsertArguments(cols ...string) (columns string, arguments string) {
	return JoinColumns(cols...), JoinNamedValues(cols...)
}

// NamedUpdateArguments returns columns and arguments for SQL UPDATE statements based on columns.
//
//	statement := NamedUpdateArguments("foo", "bar")
//	query := fmt.Sprintf("UPDATE foo SET %s", statement)
//	// UPDATE foo SET foo=:foo, bar=:bar
func NamedUpdateArguments(cols ...string) (arguments string) {
	return JoinNamedColumnsValues(cols...)
}

// NamedWhereArguments returns conditions for SQL WHERE statements based on columns.
//
//	query := fmt.Sprintf("SELECT * FROM foo WHERE %s", NamedWhereArguments("foo", "bar"))
//	// SELECT * FROM foo WHERE foo=:foo AND bar=:bar
//
//	query := fmt.Sprintf("SELECT * FROM foo WHERE %s", NamedWhereArguments())
//	// SELECT * FROM foo WHERE TRUE
func NamedWhereArguments(operator SqlOperator, cols ...string) (arguments string) {
	if len(cols) == 0 {
		return "TRUE"
	}
	return JoinNamedCondition(operator, cols...)
}
