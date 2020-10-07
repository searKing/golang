// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"

	"github.com/searKing/golang/go/reflect"
)

const (
	TagSqlx = "db"
	TagJson = "json"
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
	// The name of the type of the constants or variables we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
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
			continue
		}

		v := Value{
			originalName: typ,
			str:          typ,

			valueImport:     f.typeInfo.valueImport,
			valueType:       f.typeInfo.valueType,
			valueIsPointer:  f.typeInfo.valueIsPointer,
			valueTypePrefix: f.typeInfo.valueTypePrefix,
		}
		if c := tspec.Comment; f.lineComment && c != nil && len(c.List) == 1 {
			v.name = strings.TrimSpace(c.Text())
		} else {
			v.name = strings.TrimPrefix(typ, f.trimPrefix)
		}

		v.eleName = v.name

		if strings.TrimSpace(v.valueType) == "" {
			v.valueType = "interface{}"
		}
		f.values = append(f.values, v)
		for _, field := range sExpr.Fields.List {
			fieldName := ""
			if len(field.Names) != 0 {
				for _, field := range field.Names {
					if !*skipPrivateFields || isPublicName(field.Name) {
						fieldName = field.Name
						break
					}
				}
			} else if field.Names == nil { // anonymous field
				ident, ok := field.Type.(*ast.Ident)
				if !ok {
					continue
				}

				if !*skipPrivateFields {
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

			if field.Tag != nil {
				var tag string
				if field.Tag.Value != "" {
					var err error
					tag, err = strconv.Unquote(field.Tag.Value)
					if err != nil {
						panic(err)
					}
				}
				tags, err := reflect.ParseStructTag(tag)
				if err != nil {
					panic(err)
				}
				{
					_, has := tags.Get(TagSqlx)
					if !has {
						tags.AddOptions(TagSqlx, fieldName)
					}
				}
				{
					_, has := tags.Get(TagJson)
					if !has {
						tags.AddOptions(TagJson, fieldName)
					}
				}
				field.Tag.Value = tags.String()
			}
		}
	}
	return false
}
