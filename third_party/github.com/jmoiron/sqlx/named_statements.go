package sqlx

import (
	"fmt"
	"strings"

	strings_ "github.com/searKing/golang/go/strings"
)

type InsertOption int

const (
	// INSERT inserts new rows into an existing table.
	InsertOptionInsert InsertOption = iota
	// REPLACE works exactly like INSERT, except that if an old row in the table has the same value as a
	// new row for a PRIMARY KEY or a UNIQUE index, the old row is deleted before the new row is inserted.
	InsertOptionReplace InsertOption = iota
	// If you specify an ON DUPLICATE KEY UPDATE clause and a row to be inserted would cause a duplicate
	// value in a UNIQUE index or PRIMARY KEY, an UPDATE of the old row occurs.
	InsertOptionUpdate InsertOption = iota
	// If you use the IGNORE modifier, ignorable errors that occur while executing the INSERT statement
	// are ignored.
	InsertOptionIgnore InsertOption = iota
	// deprecated in MySQL 5.6. In MySQL 8.0, DELAYED is not supported
	// Deprecated: Use InsertOptionInsert instead.
	InsertOptionDelayed InsertOption = iota
)

// SimpleStatements is a simple render for simple SQL
type SimpleStatements struct {
	TableName      string
	Columns        []string
	Conditions     []string    // take effect only in WHERE clause, that exists in SELECT, UPDATE, DELETE, or UPDATE in INSERT
	Compare        SqlCompare  // take effect only in WHERE clause, that exists in SELECT, UPDATE, DELETE
	Operator       SqlOperator // take effect only in WHERE clause, that exists in SELECT, UPDATE, DELETE
	Limit          int         // take effect only in SELECT
	Offset         int         // take effect only in SELECT
	GroupByColumns []string    // take effect only in WHERE clause, that exists in SELECT, UPDATE, DELETE
	OrderByColumns []string    // take effect only in WHERE clause, that exists in SELECT, UPDATE, DELETE

	// INSERT
	InsertOption InsertOption // take effect only in INSERT
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
	if len(s.GroupByColumns) > 0 {
		query += fmt.Sprintf(" GROUP BY %s", JoinTableColumns(s.TableName, s.GroupByColumns...))
	}
	if len(s.OrderByColumns) > 0 {
		query += fmt.Sprintf(" ORDER BY %s", JoinTableColumns(s.TableName, s.OrderByColumns...))
	}
	return query + " " + strings.Join(appends, "")
}

// NamedInsertStatement returns a simple sql statement for SQL INSERT statements based on columns.
//
//	statement := SimpleStatements{
//		TableName: foo,
//		Columns: []string{"foo", "bar"},
//		Conditions: []string{"qux"},
//	}
//	query := statement.NamedInsertStatement(false)
//
//	// INSERT INTO foo (foo, bar, qux) VALUES (:foo, :bar, :qux)
//
//	query := statement.NamedInsertStatement(true)
//
//	// INSERT INTO foo (foo, bar, qux) VALUES (:foo, :bar, :qux) ON DUPLICATE KEY UPDATE foo=:foo, bar=:bar
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
	insertOption := s.InsertOption
	if len(s.Columns)+len(s.Conditions) > 0 && update {
		insertOption = InsertOptionUpdate
	}
	cols := strings_.SliceUnique(strings_.SliceCombine(s.Columns, s.Conditions)...)
	switch insertOption {
	case InsertOptionReplace:
		return fmt.Sprintf(`REPLACE INTO %s %s`,
			s.TableName,
			NamedInsertArgumentsCombined(cols...)) + " " + strings.Join(appends, "")
	case InsertOptionUpdate:
		return fmt.Sprintf(`INSERT INTO %s %s ON DUPLICATE KEY UPDATE %s`,
			s.TableName,
			NamedInsertArgumentsCombined(cols...),
			NamedUpdateArguments(s.Columns...)) + " " + strings.Join(appends, "")
	case InsertOptionIgnore:
		return fmt.Sprintf(`INSERT IGNORE INTO %s %s`,
			s.TableName,
			NamedInsertArgumentsCombined(cols...)) + " " + strings.Join(appends, "")
	case InsertOptionDelayed:
		return fmt.Sprintf(`INSERT DELAYED INTO %s %s`,
			s.TableName,
			NamedInsertArgumentsCombined(cols...)) + " " + strings.Join(appends, "")
	case InsertOptionInsert:
		fallthrough
	default:
		return fmt.Sprintf(`INSERT INTO %s %s`,
			s.TableName,
			NamedInsertArgumentsCombined(cols...)) + " " + strings.Join(appends, "")
	}
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
