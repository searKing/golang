// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"github.com/searKing/golang/go/reflect"
	strings_ "github.com/searKing/golang/go/strings"
)

const (
	TagSqlx = "db"
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
		sExpr, ok := tspec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		if typ != f.typeInfo.Name {
			// This is not the type we're looking for.
			continue
		}

		if sExpr.Fields.NumFields() <= 0 {
			panic(fmt.Errorf("%s has no Fields", typ))
		}

		v := Struct{
			StructType: typ,
		}
		if c := tspec.Comment; f.lineComment && c != nil && len(c.List) == 1 {
			v.TableName = strings.TrimSpace(c.Text())
		} else {
			v.TableName = strings_.SnakeCase(strings.TrimPrefix(typ, f.trimPrefix))
		}

		for _, field := range sExpr.Fields.List {
			fieldName := ""
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
			{
				tagName := strings_.SnakeCase(fieldName)
				var tagChanged bool
				{
					_, has := tags.Get(TagSqlx)
					if !has {
						tagChanged = true
						tags.SetName(TagSqlx, tagName)
						//tags.AddOptions(TagSqlx, "omitempty")
					}
				}

				if tagChanged {
					f.fileChanged = true
					field.Tag.Value = tags.AstString()
				}
			}
			tagSqlx, _ := tags.Get(TagSqlx)
			if tagSqlx.Name == "-" {
				// ignore this field
				continue
			}
			v.Fields = append(v.Fields, StructField{
				FieldName: fieldName,
				DbName:    tagSqlx.Name,
			})
		}
		if len(v.Fields) == 0 {
			panic(fmt.Errorf("%s has no Fields with tag %q", typ, TagSqlx))
		}
		f.structs = append(f.structs, v)
	}
	return false
}
