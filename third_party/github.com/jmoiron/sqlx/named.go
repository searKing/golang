// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlx

import (
	"errors"

	"vitess.io/vitess/go/vt/sqlparser"
)

//go:generate go-option -type "compileQuery"
type compileQuery struct {
	// Generate Alias
	// `SELECT t.a, b`
	// TO
	// `select t.a as t_a, b as b`,
	AliasWithSelect bool
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
	err = NamedUnbindVars(stmts, opt.AliasWithSelect)
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
func NamedUnbindVars(stmt sqlparser.Statement, withAlias bool) error {
	return sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node := node.(type) {
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
						return false, errors.New("mismatched cloumns and values")

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

		case *sqlparser.ColName, sqlparser.TableName:
			// Common node types that never contain SQLVals or ListArgs but create a lot of object
			// allocations.
			return false, nil
		}
		return true, nil
	}, stmt)
}
