// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

import (
	"errors"
	"reflect"

	reflect_ "github.com/searKing/golang/go/reflect"
	"vitess.io/vitess/go/vt/sqlparser"
)

const (
	TagDb = "db"
)

//go:generate go-option -type "compileQuery"
type compileQuery struct {
	// Generate Alias
	// `SELECT t.a, b`
	// TO
	// `select t.a as t_a, b as b`,
	AliasWithSelect bool

	// Argument by column name
	// nil: keep all
	// []string: keep if exist
	// map[string]interface{} : keep if exist and none zero
	// struct{} tag is `db:"{{col_name}}"`: keep if exist and none zero
	// take effect in WHERE|INSERT|UPDATE, ignore if multi rows
	Argument interface{}
}

// CompileQuery compiles a unbound query (using the '?' bind var) into an NamedQuery.
// WithCompileQueryOptionAliasWithSelect, default true
// SELECT t.a, b FROM t
// 	TO
// select t.a as t_a, b as b
//
// INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)
// TO
// insert into foo(a, b, c, d) values (:a, :b, :c, :d)
func CompileQuery(sql string, opts ...CompileQueryOption) (query string, err error) {
	var opt compileQuery
	opt.ApplyOptions(WithCompileQueryOptionAliasWithSelect(true))
	opt.ApplyOptions(opts...)
	stmts, err := sqlparser.Parse(sql)
	if err != nil {
		return "", err
	}
	err = NamedUnbindVars(stmts, opt.AliasWithSelect, opt.Argument)
	if err != nil {
		return "", err
	}
	// Generate query while simultaneously resolving values.
	buf := sqlparser.NewTrackedBuffer(nil)
	stmts.Format(buf)
	return buf.String(), nil
}

// CompliantName returns a compliant id name
// that can be used for a bind or as var.
// replace special runes with '_'
// a.b -> a_b
func CompliantName(in string) string {
	return sqlparser.NewColIdent(in).CompliantName()
}

// NamedUnbindVars rewrites unbind vars to named vars referenced in the statement.
// Ideally, this should be done only once.
func NamedUnbindVars(stmt sqlparser.Statement, withAlias bool, arg interface{}) error {
	return sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node := node.(type) {
		case *sqlparser.Where:
			TrimWhere(node, arg)
			break
		case *sqlparser.Update:
			if err := TrimUpdate(node, arg); err != nil {
				return false, err
			}
			break
		case *sqlparser.AliasedExpr:
			if !withAlias {
				break
			}
			buf := sqlparser.NewTrackedBuffer(nil)
			node.Expr.Format(buf)
			node.As = sqlparser.NewColIdent(CompliantName(buf.String()))
		case *sqlparser.UpdateExpr:
			val, ok := node.Expr.(*sqlparser.SQLVal)
			if !ok {
				break
			}
			if val.Type != sqlparser.ValArg {
				break
			}
			buf := sqlparser.NewTrackedBuffer(nil)
			node.Name.Format(buf)
			val.Val = []byte(":" + buf.String())
		case *sqlparser.ComparisonExpr:
			val, ok := node.Right.(*sqlparser.SQLVal)
			if !ok {
				break
			}
			if val.Type != sqlparser.ValArg {
				break
			}
			buf := sqlparser.NewTrackedBuffer(nil)
			node.Left.Format(buf)
			val.Val = []byte(":" + buf.String())

		case *sqlparser.Insert:
			valTuples, ok := node.Rows.(sqlparser.Values)
			if !ok {
				break
			}
			for _, vals := range valTuples {
				for i, val := range vals {
					val, ok := val.(*sqlparser.SQLVal)
					if !ok {
						continue
					}
					if val.Type != sqlparser.ValArg {
						continue
					}
					if i > len(node.Columns) {
						return false, errors.New("mismatched columns and values")

					}
					val.Val = []byte(":" + node.Columns[i].String())
					//val.Val = []byte(":" + node.Columns[i].CompliantName())
				}
			}
			for _, updateExpr := range node.OnDup {
				val, ok := updateExpr.Expr.(*sqlparser.SQLVal)
				if !ok {
					continue
				}
				if val.Type != sqlparser.ValArg {
					continue
				}
				buf := sqlparser.NewTrackedBuffer(nil)
				updateExpr.Name.Format(buf)

				val.Val = []byte(":" + buf.String())
			}
			if err := TrimInsert(node, arg); err != nil {
				return false, err
			}

		case *sqlparser.ColName, sqlparser.TableName:
			// Common node types that never contain SQLVals or ListArgs but create a lot of object
			// allocations.
			return false, nil
		}
		return true, nil
	}, stmt)
}

func TrimWhere(where *sqlparser.Where, arg interface{}) {
	if arg == nil || isNilValue(arg) {
		return
	}
	where.Expr = TrimExpr(where.Expr, arg)
}

func TrimInsert(insert *sqlparser.Insert, arg interface{}) error {
	if arg == nil || isNilValue(arg) {
		return nil
	}
	colTuples := insert.Columns
	valTupleRows, ok := insert.Rows.(sqlparser.Values)
	if !ok {
		return nil
	}
	if len(valTupleRows) > 1 {
		return nil
	}
	valTuples := valTupleRows[0]
	if len(colTuples) != len(valTuples) {
		return errors.New("mismatched columns and values")
	}
	var filteredCols sqlparser.Columns
	var filteredVals sqlparser.ValTuple

	for i := 0; i < len(colTuples); i++ {
		col := colTuples[i]
		val := valTuples[i]

		v, ok := val.(*sqlparser.SQLVal)
		if !ok {
			filteredCols = append(filteredCols, col)
			filteredVals = append(filteredVals, v)
			continue
		}
		if v.Type != sqlparser.ValArg {
			filteredCols = append(filteredCols, col)
			filteredVals = append(filteredVals, v)
			continue
		}
		if exist(arg, col.String()) {
			filteredCols = append(filteredCols, col)
			filteredVals = append(filteredVals, v)
			continue
		}
	}
	insert.Columns = filteredCols
	valTupleRows = nil
	insert.Rows = append(valTupleRows, filteredVals)

	var filteredUpdates []*sqlparser.UpdateExpr
	for _, updateExpr := range insert.OnDup {
		val, ok := updateExpr.Expr.(*sqlparser.SQLVal)
		if !ok {
			filteredUpdates = append(filteredUpdates, &sqlparser.UpdateExpr{
				Name: updateExpr.Name,
				Expr: updateExpr.Expr,
			})
			continue
		}
		if val.Type != sqlparser.ValArg {
			filteredUpdates = append(filteredUpdates, &sqlparser.UpdateExpr{
				Name: updateExpr.Name,
				Expr: updateExpr.Expr,
			})
			continue
		}
		buf := sqlparser.NewTrackedBuffer(nil)
		updateExpr.Name.Format(buf)
		if exist(arg, buf.String()) {
			filteredUpdates = append(filteredUpdates, &sqlparser.UpdateExpr{
				Name: updateExpr.Name,
				Expr: updateExpr.Expr,
			})
			continue
		}
	}
	insert.OnDup = filteredUpdates
	return nil
}

func TrimUpdate(update *sqlparser.Update, arg interface{}) error {
	if arg == nil || isNilValue(arg) {
		return nil
	}

	var filteredUpdates []*sqlparser.UpdateExpr
	for _, updateExpr := range update.Exprs {
		val, ok := updateExpr.Expr.(*sqlparser.SQLVal)
		if !ok {
			filteredUpdates = append(filteredUpdates, &sqlparser.UpdateExpr{
				Name: updateExpr.Name,
				Expr: updateExpr.Expr,
			})
			continue
		}
		if val.Type != sqlparser.ValArg {
			filteredUpdates = append(filteredUpdates, &sqlparser.UpdateExpr{
				Name: updateExpr.Name,
				Expr: updateExpr.Expr,
			})
			continue
		}
		buf := sqlparser.NewTrackedBuffer(nil)
		updateExpr.Name.Format(buf)

		if exist(arg, buf.String()) {
			filteredUpdates = append(filteredUpdates, &sqlparser.UpdateExpr{
				Name: updateExpr.Name,
				Expr: updateExpr.Expr,
			})
			continue
		}
	}
	update.Exprs = filteredUpdates
	return nil
}

func TrimExpr(expr sqlparser.Expr, arg interface{}) sqlparser.Expr {
	if arg == nil || isNilValue(arg) {
		return nil
	}
	switch expr := expr.(type) {
	case *sqlparser.ComparisonExpr:
		buf := sqlparser.NewTrackedBuffer(nil)
		expr.Left.Format(buf)
		if exist(arg, buf.String()) {
			return expr
		}
		val, ok := expr.Right.(*sqlparser.SQLVal)
		if !ok {
			return expr
		}
		if val.Type != sqlparser.ValArg {
			return expr
		}
		// arg this node
		return nil
	case *sqlparser.AndExpr:
		rightExpr := TrimExpr(expr.Right, arg)
		leftExpr := TrimExpr(expr.Left, arg)
		if leftExpr == nil && rightExpr == nil {
			return nil
		}
		if leftExpr == nil {
			return rightExpr
		}
		if rightExpr == nil {
			return leftExpr
		}
		expr.Left = leftExpr
		expr.Right = rightExpr
		return expr
	case *sqlparser.OrExpr:
		rightExpr := TrimExpr(expr.Right, arg)
		leftExpr := TrimExpr(expr.Left, arg)
		if leftExpr == nil && rightExpr == nil {
			return nil
		}
		if leftExpr == nil {
			return rightExpr
		}
		if rightExpr == nil {
			return leftExpr
		}
		expr.Left = leftExpr
		expr.Right = rightExpr
		return expr
	case *sqlparser.XorExpr:
		rightExpr := TrimExpr(expr.Right, arg)
		leftExpr := TrimExpr(expr.Left, arg)
		if leftExpr == nil && rightExpr == nil {
			return nil
		}
		if leftExpr == nil {
			return rightExpr
		}
		if rightExpr == nil {
			return leftExpr
		}
		expr.Left = leftExpr
		expr.Right = rightExpr
		return expr
	case *sqlparser.NotExpr:
		expr.Expr = TrimExpr(expr.Expr, arg)
		if expr.Expr == nil {
			return nil
		}
		return expr
	default:
		return expr
	}
}

// exist returns true if col exist and not zero
func exist(arg interface{}, col string) bool {
	if arg == nil {
		return false
	}
	switch arg := arg.(type) {
	case map[string]interface{}:
		val, has := arg[col]
		if !has {
			return false
		}
		var zero = reflect_.IsZeroValue(reflect.ValueOf(val))
		return !zero
	case []string:
		for _, t := range arg {
			if t == col {
				return true
			}
		}
		return false
	default:
	}

	val := reflect.ValueOf(arg)
	t := val.Type()
	if val.Kind() == reflect.Ptr {
		if val.IsNil() { // create the structure if it's nil
			s := reflect.New(val.Type().Elem())
			val.Set(s)
			val = s
		}

		val = val.Elem()
		t = t.Elem()
	}

	if t.Kind() == reflect.Map {
		for _, e := range val.MapKeys() {
			if e.Kind() != reflect.String {
				continue
			}
			if e.String() != col {
				continue
			}
			var zero = reflect_.IsZeroValue(val.MapIndex(e))
			return !zero
		}
		return false
	}

	if t.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// unwrap any payloads
		if payload := field.Tag.Get(TagDb); payload != "" {
			if payload == col {
				var zero = reflect_.IsZeroValue(val.FieldByName(field.Name))
				return !zero
			}
		}
	}
	return false
}

func isNilValue(i interface{}) bool {
	valueOf := reflect.ValueOf(i)
	kind := valueOf.Kind()
	isNullable := kind == reflect.Ptr || kind == reflect.Array || kind == reflect.Slice
	return isNullable && valueOf.IsNil()
}
