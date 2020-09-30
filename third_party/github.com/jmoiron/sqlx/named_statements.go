package sqlx

import (
	"fmt"
	"strings"
)

// SimpleStatements is a simple render for simple SQL
type SimpleStatements struct {
	TableName  string
	Columns    []string
	Conditions []string    // take effect only in WHERE clause, that exists in SELECT, UPDATE, DELETE
	Compare    SqlCompare  // take effect only in WHERE clause, that exists in SELECT, UPDATE, DELETE
	Operator   SqlOperator // take effect only in WHERE clause, that exists in SELECT, UPDATE, DELETE
	Limit      int         // take effect only in SELECT
	Offset     int         // take effect only in SELECT
}

// NamedSelectStatement returns a simple sql statement for SQL SELECT statements based on columns.
//
//	query := SimpleStatements{
//		TableName: foo,
//		Columns: []string{"foo", "bar"}
//		Conditions: []string{"thud", "grunt"}
//	}.NamedSelectStatement()
//
//	// SELECT foo, bar FROM foo WHERE thud=:thud AND grunt=:grunt
//
//	query := SimpleStatements{
//		TableName: foo,
//	}.NamedSelectStatement()
//
//	// SELECT * FROM foo WHERE TRUE
func (s SimpleStatements) NamedSelectStatement(appends ...string) string {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s`,
		NamedSelectArguments(s.Columns...),
		s.TableName,
		NamedWhereArguments(s.Compare, s.Operator, s.Conditions...))

	if s.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", s.Limit, s.Offset)
	}
	return query + " " + strings.Join(appends, "")
}

// NamedInsertStatement returns a simple sql statement for SQL INSERT statements based on columns.
//
//	statement := SimpleStatements{
//		TableName: foo,
//		Columns: []string{"foo", "bar"}
//	}
//	query := statement.NamedInsertStatement(false)
//
//	// INSERT INTO foo (foo, bar) VALUES (:foo, :bar)
//
//	query := statement.NamedInsertStatement(true)
//
//	// INSERT INTO foo (foo, bar) VALUES (:foo, :bar) ON DUPLICATE KEY UPDATE foo=:foo, bar=:bar
//
//	statement := SimpleStatements{
//		TableName: foo,
//	}
//	query := statement.NamedSelectStatement(false)
//
//	// INSERT INTO foo DEFAULT VALUES
//
//	query := statement.NamedSelectStatement(true)
//
//	// INSERT INTO foo DEFAULT VALUES
func (s SimpleStatements) NamedInsertStatement(update bool, appends ...string) string {
	if len(s.Columns) > 0 && update {
		return fmt.Sprintf(`INSERT INTO %s %s ON DUPLICATE KEY UPDATE %s`,
			s.TableName,
			NamedInsertArgumentsCombined(s.Columns...),
			NamedUpdateArguments(s.Columns...)) + " " + strings.Join(appends, "")
	}
	return fmt.Sprintf(`INSERT INTO %s %s`,
		s.TableName,
		NamedInsertArgumentsCombined(s.Columns...)) + " " + strings.Join(appends, "")
}

// NamedUpdateStatement returns a simple sql statement for SQL UPDATE statements based on columns.
//
//	statement := SimpleStatements{
//		TableName: foo,
//		Columns: []string{"foo", "bar"},
//		Conditions: []string{"thud", "grunt"},
//		Operator: SqlOperatorAnd,
//	}
//	query := statement.NamedUpdateStatement(false)
//
//	// UPDATE foo SET foo=:foo, bar=:bar WHERE thud=:thud AND grunt=:grunt
//
//	query := statement.NamedUpdateStatement(true)
//
//	// INSERT INTO foo (foo, bar) VALUES (:foo, :bar) ON DUPLICATE KEY UPDATE foo=:foo, bar=:bar
//
//	statement := SimpleStatements{
//		TableName: foo,
//		Columns: []string{"foo"},
//	}
//	query := statement.NamedUpdateStatement(false)
//
//	// UPDATE foo SET foo=:foo WHERE TRUE
//
//	query := statement.NamedUpdateStatement(true)
//
//	// INSERT INTO foo (foo) VALUES (:foo) ON DUPLICATE KEY UPDATE foo=:foo
//
//	statement := SimpleStatements{
//		TableName: foo,
//	}
//	query := statement.NamedUpdateStatement(false)
//
//  // Malformed SQL
//	// UPDATE foo SET WHERE TRUE
//
//	query := statement.NamedUpdateStatement(true)
//
//	// INSERT INTO foo DEFAULT VALUES
func (s SimpleStatements) NamedUpdateStatement(insert bool, appends ...string) string {
	if insert {
		return s.NamedInsertStatement(true)
	}
	return fmt.Sprintf(`UPDATE %s SET %s WHERE %s`,
		s.TableName,
		NamedUpdateArguments(s.Columns...),
		NamedWhereArguments(s.Compare, s.Operator, s.Conditions...)) + " " + strings.Join(appends, "")
}

// NamedDeleteStatement returns a simple sql statement for SQL DELETE statements based on columns.
//
//	query := SimpleStatements{
//		TableName: foo,
//		Conditions: []string{"thud", "grunt"}
//	}.NamedUpdateStatement()
//
//	// DELETE FROM foo WHERE thud=:thud AND grunt=:grunt
//
//	query := SimpleStatements{
//		TableName: foo,
//		Columns: []string{"foo"}
//	}.NamedUpdateStatement()
//
//	// DELETE FROM foo WHERE TRUE
func (s SimpleStatements) NamedDeleteStatement(appends ...string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE %s`,
		s.TableName,
		NamedWhereArguments(s.Compare, s.Operator, s.Conditions...)) + " " + strings.Join(appends, "")
}
