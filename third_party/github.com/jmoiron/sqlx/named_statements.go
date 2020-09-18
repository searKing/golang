package sqlx

import (
	"fmt"
)

// SimpleStatements is a simple render for simple SQL
type SimpleStatements struct {
	TableName  string
	Columns    []string
	Conditions []string
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
func (s SimpleStatements) NamedSelectStatement() string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE %s`,
		NamedSelectArguments(s.Columns...),
		s.TableName,
		NamedWhereArguments(s.Conditions...))
}

// NamedInsertStatement returns a simple sql statement for SQL INSERT statements based on columns.
//
//	statement := SimpleStatements{
//		TableName: foo,
//		Columns: []string{"foo", "bar"}
//		Conditions: []string{"thud", "grunt"}
//	}
//	query := statement.NamedInsertStatement(false)
//
//	// INSERT INTO foo (foo, bar, thud, grunt) VALUES (:foo, :bar, :thud, :grunt)
//
//	query := statement.NamedInsertStatement(true)
//
//	// INSERT INTO foo (foo, bar, thud, grunt) VALUES (:foo, :bar, :thud, :grunt) ON DUPLICATE KEY UPDATE thud=:thud, grunt=:grunt
//
//	statement := SimpleStatements{
//		TableName: foo,
//	}
//	query	:= statement.NamedSelectStatement(false)
//
//	// INSERT INTO foo DEFAULT VALUES
//
//	query	:= statement.NamedSelectStatement(true)
//
//	// INSERT INTO foo DEFAULT VALUES
func (s SimpleStatements) NamedInsertStatement(update bool) string {
	if update {
		if len(s.Columns)+len(s.Conditions) == 0 {
			return fmt.Sprintf(`INSERET INTO %s %s`,
				s.TableName,
				NamedInsertArgumentsCombined())
		}
		return fmt.Sprintf(`INSERET INTO %s %s ON DUPLICATE KEY UPDATE %s`,
			s.TableName,
			NamedInsertArgumentsCombined(s.Columns...),
			NamedUpdateArguments(s.Conditions...))
	}
	var cols []string
	cols = append(cols, s.Columns...)
	cols = append(cols, s.Conditions...)
	return fmt.Sprintf(`INSERET INTO %s %s`,
		s.TableName,
		NamedInsertArgumentsCombined(cols...))
}

// NamedUpdateStatement returns a simple sql statement for SQL UPDATE statements based on columns.
//
//	query := SimpleStatements{
//		TableName: foo,
//		Columns: []string{"foo", "bar"}
//		Conditions: []string{"thud", "grunt"}
//	}.NamedUpdateStatement()
//
//	// UPDATE foo SET foo=:foo, bar=:bar WHERE thud=:thud AND grunt=:grunt
//
//	query := SimpleStatements{
//		TableName: foo,
//		Columns: []string{"foo"}
//	}.NamedUpdateStatement()
//
//	// UPDATE foo SET foo=:foo WHERE TRUE
//
//	query := SimpleStatements{
//		TableName: foo,
//	}.NamedUpdateStatement()
//
//  // Malformed SQL
//	// UPDATE foo SET WHERE TRUE
func (s SimpleStatements) NamedUpdateStatement() string {
	return fmt.Sprintf(`UPDATE %s SET %s WHERE %s`,
		s.TableName,
		NamedUpdateArguments(s.Columns...),
		NamedWhereArguments(s.Conditions...))
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
func (s SimpleStatements) NamedDeleteStatement() string {
	return fmt.Sprintf(`DELETE FROM %s WHERE %s`,
		s.TableName,
		NamedWhereArguments(s.Conditions...))
}
