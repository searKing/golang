package sqlx

import (
	"fmt"
)

// TableSelectArguments returns columns and arguments for SQL SELECT statements based on columns.
//
//	query := fmt.Sprintf("SELECT %s FROM foo", TableSelectArguments("table", "foo", "bar"))
//	// SELECT table.foo, table.bar FROM table
//
//	query := fmt.Sprintf("SELECT %s FROM table", TableSelectArguments(table", "foo"))
//	// SELECT table.foo FROM table
//
//	query := fmt.Sprintf("SELECT %s FROM table", TableSelectArguments())
//	// SELECT * FROM table
func TableSelectArguments(table string, cols ...string) (arguments string) {
	if len(cols) == 0 {
		return "*"
	}
	return JoinTableColumns(table, cols...)
}

// TableInsertArgumentsCombined returns columns and arguments together
// for SQL INSERT statements based on columns.
//
//	query := fmt.Sprintf("INSERT INTO table %s", TableInsertArgumentsCombined("table", "foo", "bar"))
//	// INSERT INTO table (table.foo, table.bar) VALUES (?, ?)
//
//	query := fmt.Sprintf("INSERT INTO table %s", TableInsertArgumentsCombined("table"))
//	// INSERT INTO table (foo, bar) DEFAULT VALUES
func TableInsertArgumentsCombined(table string, cols ...string) (arguments string) {
	//if len(cols) == 0 {
	//	// https://docs.microsoft.com/en-us/sql/t-sql/statements/insert-transact-sql?view=sql-server-ver15#d-inserting-data-into-a-table-with-columns-that-have-default-values
	//	return `DEFAULT VALUES`
	//}
	if len(cols) == 0 {
		return `VALUES(DEFAULT)`
	}
	return fmt.Sprintf(`(%s) VALUES (%s)`, JoinTableColumns(table, cols...), JoinTableValues(cols...))
}

// TableInsertArguments returns columns and arguments for SQL INSERT statements based on columns.
//
//	columns, arguments := TableInsertArguments("table", "foo", "bar")
//	query := fmt.Sprintf("INSERT INTO table (%s) VALUES (%s)", columns, arguments)
//	// INSERT INTO table (table.foo, table.bar) VALUES (?, ?)
// Deprecated: Use TableInsertArgumentsCombined instead.
func TableInsertArguments(table string, cols ...string) (columns string, arguments string) {
	return JoinTableColumns(table, cols...), JoinTableValues(cols...)
}

// TableUpdateArguments returns columns and arguments for SQL UPDATE statements based on columns.
//
//	statement := TableUpdateArguments("table", "foo", "bar")
//	query := fmt.Sprintf("UPDATE table SET %s", statement)
//	// UPDATE foo SET table.foo=?, table.bar=?
func TableUpdateArguments(table string, cols ...string) (arguments string) {
	return JoinTableColumnsValues(table, cols...)
}

// TableWhereArguments returns conditions for SQL WHERE statements based on columns.
//
//	query := fmt.Sprintf("SELECT * FROM table WHERE %s", TableWhereArguments(SqlCompareEqual, SqlOperatorAnd, "table", "foo", "bar"))
//	// SELECT * FROM table WHERE table.foo=? AND table.bar=?
//
//	query := fmt.Sprintf("SELECT * FROM table WHERE %s", TableWhereArguments(SqlCompareEqual, SqlOperatorAnd, "table"))
//	// SELECT * FROM table WHERE TRUE
func TableWhereArguments(cmp SqlCompare, operator SqlOperator, table string, cols ...string) (arguments string) {
	if len(cols) == 0 {
		return "TRUE"
	}
	return JoinTableCondition(cmp, operator, table, cols...)
}
