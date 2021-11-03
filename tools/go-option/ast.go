// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"github.com/searKing/golang/go/reflect"
)

const (
	TagOption = "option"
)

func isPublicName(name string) bool {
	for _, c := range name {
		return unicode.IsUpper(c)
	}
	return false
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	// Token must be in IMPORT, CONST, TYPE, VAR
	if !ok || decl.Tok != token.TYPE {
		// We only care about const|var declarations.
		return true
	}
	// The trimmedStructName of the type of the constants or variables we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by structs.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		tspec := spec.(*ast.TypeSpec) // Guaranteed to succeed as this is TYPE.
		typ = tspec.Name.Name
		v := Struct{
			StructTypeName:   typ,
			StructTypeImport: f.typeInfo.Import,
		}
		if c := tspec.Comment; f.lineComment && c != nil && len(c.List) == 1 {
			v.trimmedStructTypeName = strings.TrimSpace(c.Text())
		} else {
			v.trimmedStructTypeName = strings.TrimPrefix(v.StructTypeName, f.trimPrefix)
		}
		if typ != f.typeInfo.Name {
			// This is not the type we're looking for.
			continue
		}
		sExpr, ok := tspec.Type.(*ast.StructType)
		if !ok {
			f.structs = append(f.structs, v)
			continue
		}

		for _, field := range sExpr.Fields.List {
			fieldName := ""
			fieldType := ""
			{
				ident, ok := field.Type.(*ast.Ident)
				if !ok {
					continue
				}
				fieldType = ident.String()
			}
			if len(field.Names) != 0 { // pick first exported Name
				for _, field := range field.Names {
					if !*flagSkipPrivateFields || isPublicName(field.Name) {
						fieldName = field.Name
						break
					}
				}
			} else { // anonymous field
				ident, ok := field.Type.(*ast.Ident)
				if !ok {
					continue
				}

				if !*flagSkipAnonymousFields {
					fieldName = ident.Name
				}
			}

			// nothing to process, continue with next line
			if fieldName == "" {
				continue
			}
			if field.Tag == nil {
				field.Tag = &ast.BasicLit{}
			}

			tags, err := reflect.ParseAstStructTag(field.Tag.Value)
			if err != nil {
				panic(err)
			}
			tagOption, _ := tags.Get(TagOption)
			if tagOption.Name == "-" {
				// ignore this field
				continue
			}
			v.Fields = append(v.Fields, StructField{
				FieldName: fieldName,
				FieldType: fieldType,
				TagName:   tagOption.Name,
			})
		}
		f.structs = append(f.structs, v)
	}
	return false
}
