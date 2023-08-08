// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import "go/ast"

// Indirect returns the ast.Expr that x points to.
// If x is an ast.StarExpr, Indirect returns a zero ast.Expr.
// If x is not an ast.StarExpr, Indirect returns x.
func Indirect(x ast.Expr, d int) (ast.Expr, int) {
	switch t := x.(type) {
	case *ast.StarExpr:
		return Indirect(t.X, d+1)
	}
	return x, d
}
