// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast_test

import (
	"go/ast"
	"reflect"
	"strconv"
	"strings"
	"testing"

	ast_ "github.com/searKing/golang/go/ast"
)

// TestIndirect tests the Indirect function in Go source code.
func TestIndirect(t *testing.T) {
	tests := []struct {
		name       string
		inputExpr  ast.Expr
		inputDepth int
		wantExpr   ast.Expr
		wantDepth  int
	}{
		{"star", &ast.StarExpr{}, 0, nil, 1},
		{"paren", &ast.ParenExpr{}, 0, &ast.ParenExpr{}, 0},
		{"star-star", &ast.StarExpr{X: &ast.StarExpr{}}, 0, nil, 2},
		{"star-star-paren", &ast.StarExpr{X: &ast.StarExpr{X: &ast.ParenExpr{}}}, 0, &ast.ParenExpr{}, 2},
		{"star-paren", &ast.StarExpr{X: &ast.ParenExpr{}}, 0, &ast.ParenExpr{}, 1},
		{"selector", &ast.SelectorExpr{}, 0, &ast.SelectorExpr{}, 0},
		{"index", &ast.IndexExpr{}, 0, &ast.IndexExpr{}, 0},
		{"slice", &ast.SliceExpr{}, 0, &ast.SliceExpr{}, 0},
		{"call", &ast.CallExpr{}, 0, &ast.CallExpr{}, 0},
		{"unary", &ast.UnaryExpr{}, 0, &ast.UnaryExpr{}, 0},
		{"binary", &ast.BinaryExpr{}, 0, &ast.BinaryExpr{}, 0},
		{"basicLit", &ast.BasicLit{}, 0, &ast.BasicLit{}, 0},
	}

	for i, tt := range tests {
		t.Run(strings.Join([]string{strconv.Itoa(i), tt.name}, ":"), func(t *testing.T) {
			gotExpr, gotDepth := ast_.Indirect(tt.inputExpr, tt.inputDepth)
			if !reflect.DeepEqual(gotExpr, tt.wantExpr) || gotDepth != tt.wantDepth {
				t.Errorf("Indirect(%v, %v) = (%v, %v), want (%v,%v)", tt.inputExpr, 0, gotExpr, gotDepth, tt.wantExpr, tt.wantDepth)
			}
		})
	}
}
