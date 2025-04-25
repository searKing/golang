// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package union

import (
	_ "embed"
	"fmt"
	"go/ast"
	"strings"

	slices_ "github.com/searKing/golang/go/exp/slices"
	reflect_ "github.com/searKing/golang/go/reflect"
	strings_ "github.com/searKing/golang/go/strings"
)

//go:embed union.tmpl
var tmplUnion string

type TmplUnionRender struct {
	// Print the header and package clause.
	GoUnionToolName       string
	GoUnionToolArgs       []string
	GoUnionToolArgsJoined string

	PackageName string
	ImportPaths []string
	ValDecls    []string

	TargetTypeName               string // type name of target type
	TargetTypeImport             string // import path of target type
	TargetTypeGenericDeclaration string // the Generic type of the struct type
	TargetTypeGenericParams      string // the Generic params of the struct type

	FormatTypeName string        // The format FieldName of the struct type.
	Fields         []StructField // fields if target type is struct

	UnionInterfaceName string // union interface name of target type
	UnionStructName    string // union struct name of target type

	ApplyUnionsAsMemberFunction bool // ApplyUnions can be registered as UnionType's member function
	WithTargetTypeNameAsPrefix  bool // WithXXX() can be generated as {{UnionType}}WithXXX()
}

// Struct represents a declared constant.
type Struct struct {
	FileImports                  []string // The import path of the file contains the struct
	StructTypeImport             string   // The import path of StructTypeName.
	StructTypeName               string   // The StructTypeName of the struct.
	StructTypeGenericDeclaration string   // the Generic type of the struct type
	StructTypeGenericTypeParams  string   // the Generic params of the struct type
	IsStruct                     bool
	Fields                       []StructField
}

type StructField struct {
	FieldName                 string                // The FieldName of the struct field.
	FieldType                 string                // The FieldType of the struct field.
	FieldDocComment           *ast.CommentGroup     // The doc comment of the struct field.
	FieldLineComment          *ast.CommentGroup     // The line comment of the struct field.
	UnionTag                  reflect_.SubStructTag // The UnionTag of the struct field.
	FieldCanBeCompareWithNil  bool                  // The FieldType of the struct field can be compared with nil.
	FieldCanBeCompareWithZero bool                  // The FieldType of the struct field can be compared with zero.
}

func (t *TmplUnionRender) Complete() {
	t.GoUnionToolArgsJoined = strings.Join(t.GoUnionToolArgs, " ")
	t.ApplyUnionsAsMemberFunction = strings.TrimSpace(t.TargetTypeImport) == ""

	t.UnionInterfaceName = strings_.UpperCamelCaseSlice("union")

	importPath := strings.TrimSpace(t.TargetTypeImport)
	if importPath != "" {
		t.ImportPaths = append(t.ImportPaths, fmt.Sprintf("%q", importPath))
	}
	t.ImportPaths = slices_.Filter(t.ImportPaths)

	_, defaultValDecl := createValAndNameDecl(t.TargetTypeName)
	if defaultValDecl != "" {
		t.ValDecls = append(t.ValDecls, defaultValDecl)
	}

	t.FormatTypeName = strings_.ToUpperLeading(t.TargetTypeName)
}
