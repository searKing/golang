package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/searKing/golang/go/reflect"
)

const (
	TagSqlx = "db"
	TagJson = "json"
)

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
			continue
		}

		if sExpr.Fields.NumFields() <= 0 {
			continue
		}

		for _, field := range sExpr.Fields.List {

			fieldName := ""
			if len(field.Names) != 0 {
				for _, field := range field.Names {
					if !c.skipUnexportedFields || isPublicName(field.Name) {
						fieldName = field.Name
						break
					}
				}
			}




			if field.Tag != nil {

				tags, err := reflect.ParseStructTag(field.Tag.Value)
				if err != nil {
					panic(err)
				}
				tag, has := tags.Get(TagSqlx)
				if !has {
					tags.AddOptions(TagSqlx, field.Names)

				}
			}
			field.Tag

			i := field.Type.(*ast.Ident)
			fieldType := i.Name
			for _, name := range field.Names {
				fmt.Printf("\tField: name=%s type=%s\n", name.Name, fieldType)
			}
		}

		if sExpr.X.(*ast.Ident).Name == "atomic" && sExpr.Sel.Name == "Value" {
			if typ != f.typeInfo.Name {
				// This is not the type we're looking for.
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
				v.name = strings.TrimPrefix(v.originalName, f.trimPrefix)
			}
			v.eleName = v.name

			if strings.TrimSpace(v.valueType) == "" {
				v.valueType = "interface{}"
			}
			f.values = append(f.values, v)
		}

	}
	return false
}
