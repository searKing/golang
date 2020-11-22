package main

import (
	"go/ast"
	"go/token"
	"strings"
)

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	// Token must be in IMPORT, CONST, TYPE, VAR
	if !ok || decl.Tok != token.TYPE {
		// We only care about const declarations.
		return true
	}
	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""
	// Loop over the elements of the declaration. Each element is a ValueSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		tspec := spec.(*ast.TypeSpec) // Guaranteed to succeed as this is TYPE.
		typ = tspec.Name.Name
		sExpr, ok := tspec.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		if sExpr.X.(*ast.Ident).Name == "sync" && sExpr.Sel.Name == "Map" {
			if typ != f.typeInfo.mapName {
				// This is not the type we're looking for.
				continue
			}
			v := Value{
				originalName: typ,
				str:          typ,

				keyImport:       f.typeInfo.keyImport,
				keyType:         f.typeInfo.keyType,
				keyIsPointer:    f.typeInfo.keyIsPointer,
				keyTypePrefix:   f.typeInfo.keyTypePrefix,
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
			v.mapName = v.name

			if strings.TrimSpace(v.keyType) == "" {
				v.keyType = "interface{}"
			}
			if strings.TrimSpace(v.valueType) == "" {
				v.valueType = "interface{}"
			}
			f.values = append(f.values, v)
		}

	}
	return false
}
