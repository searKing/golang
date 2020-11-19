// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Arguments to format are:
//	StructType: NullJson type trimmedStructName
//	TableName: value type trimmedStructName
//	NilValue: nil value of map type
const tmplJson = `
{{ $package_scope := . }}
import (
{{- if .WithDao }}
	"context"
{{- end }}
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
{{- if .WithDao }}
	"strings"
{{- end }}

	reflect_ "github.com/searKing/golang/go/reflect"
	sqlx_ "github.com/searKing/golang/third_party/github.com/jmoiron/sqlx"
	sql_ "github.com/searKing/golang/go/database/sql"


{{- if .WithDao }}
	"github.com/jmoiron/sqlx"
{{- end }}
)

// {{.StructType}} represents an orm of table {{.TableName}}.
// {{.StructType}} implements the Scanner interface so
// it can be used as a scan destination, similar to sql.NullString.

// TableName returns table's name
func (_ {{.StructType}}) TableName() string {
	return "{{.TableName}}"
}

{{- range .Fields }}
// Column{{$package_scope.StructType}} return column name in db, from struct tag db:"{{.DbName}}"
func (_ {{$package_scope.StructType}})Column{{.FieldName}}() string{
	return "{{.DbName}}"
}

// TableColumn{{$package_scope.StructType}} return column name with TableName
// "{{$package_scope.TableName}}.{{.DbName}}"
func (_ {{$package_scope.StructType}})TableColumn{{.FieldName}}() string{
	// avoid runtime cost of fmt.Sprintf
	// return fmt.Sprintf("%s.%s", a.TableName(), a.Column{{$package_scope.StructType}}())
	return "{{$package_scope.TableName}}.{{.DbName}}"
}

// MapColumn{{$package_scope.StructType}} return column name with TableName
// "{{$package_scope.TableName}}_{{.DbName}}"
func (m {{$package_scope.StructType}})MapColumn{{.FieldName}}() string{
	// avoid runtime cost of fmt.Sprintf
	// return fmt.Sprintf("%s_%s", m.TableName(), m.Column{{$package_scope.StructType}}())
	// return "{{$package_scope.TableName}}_{{.DbName}}"
	return sql_.CompliantName(m.TableColumn{{.FieldName}}())
}
{{- end }}

func (m {{.StructType}}) Column(col {{.StructType}}Field) string {
	return col.ColumnName()
}

func (m {{.StructType}}) TableColumn(col {{.StructType}}Field) string {
	return fmt.Sprintf("%s.%s", m.TableName(), m.Column(col))
}

func (m {{.StructType}}) MapColumn(col {{.StructType}}Field) string {
	return sql_.CompliantName(m.TableColumn(col))
	//return fmt.Sprintf("%s_%s", m.TableName(), m.Column(col))
}

// MarshalMap marshal themselves into or append a valid map
func (m {{.StructType}}) MarshalMap(valueByCol map[string]interface{}) map[string]interface{} {
	if valueByCol == nil {
		valueByCol = map[string]interface{}{}
	}
{{- range .Fields}}	
	valueByCol[m.TableColumn({{$.StructType}}Field{{.FieldName}})] = m.{{.FieldName}}
{{- end}}
	return valueByCol
}

// UnmarshalMap is the interface implemented by types
// that can unmarshal a map description of themselves.
func (m *{{.StructType}}) UnmarshalMap(valueByCol map[string]interface{}) error {
	for col, val := range valueByCol {
		switch col {
{{- range .Fields}}	
		case m.MapColumn({{$.StructType}}Field{{.FieldName}}):
			// for sql.Scanner
			v := reflect.ValueOf(&m.{{.FieldName}})
			if v.Type().NumMethod() > 0 && v.CanInterface() {
				if u, ok := v.Interface().(sql.Scanner); ok {
					if err := u.Scan(val); err != nil {
						return fmt.Errorf("unmarshal col %q, got %w", col, err)
					}
					break
				}
			}

			data, err := json.Marshal(val)
			if err != nil {
				return fmt.Errorf("marshal col %q, got %w", col, err)
			}
			err = json.Unmarshal(data, &m.{{.FieldName}})
			if err != nil {
				return fmt.Errorf("unmarshal col %q, got %w", col, err)
			}
{{- end}}
		}
	}
	return nil
}
// 列名
type {{.StructType}}Field string

const (
{{- range .Fields }}
	{{$.StructType}}Field{{.FieldName}}    {{$.StructType}}Field = "{{.DbName}}"
{{- end }}
)

func (f {{.StructType}}Field) FieldName() string {
	switch f {
{{- range .Fields}}
	case {{$.StructType}}Field{{.FieldName}}:
		return "{{.FieldName}}"
{{- end}}
	}
	return string(f)
}


func (f {{.StructType}}Field) ColumnName() string {
	switch f {
{{- range .Fields}}
	case {{$.StructType}}Field{{.FieldName}}:
		return "{{.DbName}}"
{{- end}}
	}
	return string(f)
}

func (a {{.StructType}}) ColumnEditor() *{{.StructType}}Columns {
	return &{{.StructType}}Columns{
		arg: a,
	}
}

// columns

type {{.StructType}}Columns struct {
	arg  {{.StructType}}
	cols []string
}

func (c {{.StructType}}Columns) Columns(cols ...string) []string {
	return append(c.cols, cols...)
}

func (c *{{.StructType}}Columns) AppendColumn(col {{.StructType}}Field, forceAppend bool) *{{.StructType}}Columns {
	if forceAppend {
		c.cols = append(c.cols, col.ColumnName())
		return c
	}
	var zero = reflect_.IsZeroValue(reflect.ValueOf(c.arg).FieldByName(col.FieldName()))
	if !zero {
		c.cols = append(c.cols, col.ColumnName())
	}
	return c
}

func (c *{{.StructType}}Columns) AppendAll(forceAppend bool) *{{.StructType}}Columns {
	return c.
{{- range .Fields}}
		AppendColumn({{$.StructType}}Field{{.FieldName}}, forceAppend).
{{- end}}
		self()
}

func (c *{{.StructType}}Columns) self() *{{.StructType}}Columns {
	return c
}

{{- if .WithDao }}

// DAO

func (arg {{.StructType}}) Add{{.StructType}}(ctx context.Context, db *sqlx.DB, update bool) error {
	query := sqlx_.SimpleStatements{
		TableName: arg.TableName(),
		Columns:   arg.ColumnEditor().AppendAll(false).Columns(),
	}.NamedInsertStatement(update)

	_, err := db.NamedExecContext(ctx, query, arg)
	if err != nil {
{{- if .WithQueryInfo }}
		return fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return err
{{- end}}
	}
	return nil
}

func (arg {{.StructType}}) Add{{.StructType}}WithTx(ctx context.Context, tx *sqlx.Tx, update bool) error {
	query := sqlx_.SimpleStatements{
		TableName: arg.TableName(),
		Columns:   arg.ColumnEditor().AppendAll(false).Columns(),
	}.NamedInsertStatement(update)

	_, err := tx.NamedExecContext(ctx, query, arg)
	if err != nil {
{{- if .WithQueryInfo }}
		return fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return err
{{- end}}
	}
	return nil
}

func (arg {{.StructType}}) Delete{{.StructType}}(ctx context.Context, db *sqlx.DB, conds []string) error {
	query := sqlx_.SimpleStatements{
		TableName:  arg.TableName(),
		Conditions: conds, // WHERE 条件
	}.NamedDeleteStatement()

	_, err := db.NamedExecContext(ctx, query, arg)
	if err != nil {
{{- if .WithQueryInfo }}
		return fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return err
{{- end}}
	}
	return nil
}

func (arg {{.StructType}}) Delete{{.StructType}}WithTx(ctx context.Context, tx *sqlx.Tx, conds []string) error {
	query := sqlx_.SimpleStatements{
		TableName:  arg.TableName(),
		Conditions: conds, // WHERE 条件
	}.NamedDeleteStatement()

	_, err := tx.NamedExecContext(ctx, query, arg)
	if err != nil {
{{- if .WithQueryInfo }}
		return fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return err
{{- end}}
	}
	return nil
}

func (arg {{.StructType}}) Update{{.StructType}}(ctx context.Context, db *sqlx.DB, cols []string, conds []string, insert bool) error {
	query := sqlx_.SimpleStatements{
		TableName:  arg.TableName(),
		Columns:    cols,  // 要查询或修改的列名
		Conditions: conds, // WHERE 条件
	}.NamedUpdateStatement(insert)

	_, err := db.NamedExecContext(ctx, query, arg)
	if err != nil {
{{- if .WithQueryInfo }}
		return fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return err
{{- end}}
	}

	return nil
}

func (arg {{.StructType}}) Update{{.StructType}}WithTx(ctx context.Context, tx *sqlx.Tx, cols []string, conds []string, insert bool) error {
	query := sqlx_.SimpleStatements{
		TableName:  arg.TableName(),
		Columns:    cols,  // 要查询或修改的列名
		Conditions: conds, // WHERE 条件
	}.NamedUpdateStatement(insert)

	_, err := tx.NamedExecContext(ctx, query, arg)
	if err != nil {
{{- if .WithQueryInfo }}
		return fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return err
{{- end}}
	}

	return nil
}

func (arg {{.StructType}}) Get{{.StructType}}(ctx context.Context, db *sqlx.DB, cols []string, conds []string) ({{.StructType}}, error) {
	query := sqlx_.SimpleStatements{
		TableName:  arg.TableName(),
		Columns:    cols,
		Conditions: conds,
	}.NamedSelectStatement()

	// Check that invalid preparations fail
	ns, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
{{- if .WithQueryInfo }}
		return {{.StructType}}{}, fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return {{.StructType}}{}, err
{{- end}}
	}

	defer ns.Close()

	var dest {{.StructType}}
	err = ns.GetContext(ctx, &dest, arg)
	if err != nil {
		//if errors.Cause(err) == sql.ErrNoRows {
		//	return dest, nil
		//}
{{- if .WithQueryInfo }}
		return {{.StructType}}{}, fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return {{.StructType}}{}, err
{{- end}}
	}
	return dest, nil
}

func (arg {{.StructType}}) Get{{.StructType}}sByQuery(ctx context.Context, db *sqlx.DB, query string) ([]{{.StructType}}, error) {
	// Check that invalid preparations fail
	ns, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
{{- if .WithQueryInfo }}
		return nil, fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return nil, err
{{- end}}
	}

	defer ns.Close()

	var dest []{{.StructType}}
	err = ns.SelectContext(ctx, &dest, arg)
	if err != nil {
{{- if .WithQueryInfo }}
		return nil, fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return  nil, err
{{- end}}
	}
	return dest, nil
}

func (arg {{.StructType}}) Get{{.StructType}}s(ctx context.Context, db *sqlx.DB, cols []string, conds []string, likeConds []string, orderByCols []string) ([]{{.StructType}}, error) {
	query := sqlx_.SimpleStatements{
		TableName:  arg.TableName(),
		Columns:    cols,
		Conditions: conds,
		Compare:    sqlx_.SqlCompareEqual,
		Operator:   sqlx_.SqlOperatorAnd,
	}.NamedSelectStatement()
	if len(likeConds) > 0 {
		query += " AND "
		query += sqlx_.NamedWhereArguments(sqlx_.SqlCompareLike, sqlx_.SqlOperatorAnd, likeConds...)
	}
	if len(orderByCols) > 0 {
		query += fmt.Sprintf(" ORDER BY %s", sqlx_.JoinTableColumns(arg.TableName(), orderByCols...))
	}

	dest, err := arg.Get{{.StructType}}sByQuery(ctx, db, query)

	if err != nil {
{{- if .WithQueryInfo }}
		return nil, fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return nil, err
{{- end}}
	}
	return dest, nil
}


func (arg {{.StructType}}) Get{{.StructType}}sTemplate(ctx context.Context, db *sqlx.DB, limit, offset int) ([]{{.StructType}}, error) {
	query := fmt.Sprintf("SELECT %s FROM %s"+
		//" JOIN %s ON %s"+
		" %s"+ // WHERE
		" %s"+ // GROUP BY
		" %s"+ // ORDER BY
		" %s", // LIMIT
		func() string { // SELECT
			cols := sqlx_.ShrinkEmptyColumns(
				sqlx_.JoinNamedTableColumnsWithAs(arg.TableName(), arg.ColumnEditor().
{{- range .Fields}}
					AppendColumn({{$.StructType}}Field{{.FieldName}}, true).
{{- end}}
					Columns()...))

			if len(cols) == 0 {
				return "*"
			}
			return strings.Join(cols, " , ")
		}(),                        // WHERE
		arg.TableName(), // FROM
		// <other table's name'>,      // JOIN
		func() string { // WHERE
			cols := sqlx_.ShrinkEmptyColumns(
				// =
				sqlx_.JoinNamedTableCondition(sqlx_.SqlCompareEqual, sqlx_.SqlOperatorAnd,
					arg.TableName(),
					arg.ColumnEditor().
{{- range .Fields}}
						AppendColumn({{$.StructType}}Field{{.FieldName}}, true).
{{- end}}
						Columns()...),
				// <>
				sqlx_.JoinNamedTableCondition(sqlx_.SqlCompareNotEqual, sqlx_.SqlOperatorAnd,
					arg.TableName(),
					arg.ColumnEditor().
						// cols
						Columns()...),
				// LIKE
				sqlx_.JoinNamedTableCondition(sqlx_.SqlCompareLike, sqlx_.SqlOperatorAnd,
					arg.TableName(),
					arg.ColumnEditor().
						// cols
						Columns()...),
			)

			if len(cols) == 0 {
				return ""
			}

			return "WHERE " + strings.Join(cols, " "+sqlx_.SqlOperatorAnd.String()+" ")
		}(), // WHERE
		func() string { // GROUP BY
			cols := sqlx_.ShrinkEmptyColumns(
				sqlx_.JoinNamedTableColumnsWithAs(arg.TableName(), arg.ColumnEditor().
					// cols
					Columns()...))

			if len(cols) == 0 {
				return ""
			}
			return "GROUP BY " + strings.Join(cols, " , ")
		}(),  
		func() string { // ORDER BY
			cols := sqlx_.ShrinkEmptyColumns(
				sqlx_.JoinNamedTableColumnsWithAs(arg.TableName(), arg.ColumnEditor().
					// cols
					Columns()...))

			if len(cols) == 0 {
				return ""
			}
			return "ORDER BY " + strings.Join(cols, " , ")
		}(),  
		// LIMIT
		sqlx_.Pagination(limit, offset))


	// Check that invalid preparations fail
	ns, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
{{- if .WithQueryInfo }}
		return nil, fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return nil, err
{{- end}}
	}

	defer ns.Close()

	rows, err := ns.QueryxContext(ctx, arg.MarshalMap(nil))
	if err != nil {
{{- if .WithQueryInfo }}
		return nil, fmt.Errorf("%w, sql %q", err, query)
{{- else }}
		return nil, err
{{- end}}
	}

	var resps []{{.StructType}}
	for rows.Next() {
		row := make(map[string]interface{})
		err := rows.MapScan(row)
		if err != nil {
{{- if .WithQueryInfo }}
			return nil, fmt.Errorf("%w, sql %q", err, query)
{{- else }}
			return nil, err
{{- end}}
		}
		for k, v := range row {
			if b, ok := v.([]byte); ok {
				row[k] = string(b)
			}
		}

		resp := {{.StructType}}{}
		err = resp.UnmarshalMap(row)
		if err != nil {
{{- if .WithQueryInfo }}
			return nil, fmt.Errorf("%w, sql %q", err, query)
{{- else }}
			return nil, err
{{- end}}
		}
		resps = append(resps, resp)
	}
	return resps, nil
}
{{- end}}
`
