// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
	"go/ast"
	"go/token"

	"google.golang.org/protobuf/types/pluginpb"

	"github.com/searKing/golang/go/reflect"
	strings_ "github.com/searKing/golang/go/strings"

	"github.com/searKing/golang/tools/protoc-gen-go-tag/tag"
)

// GoFile holds a single parsed astFile and associated data.
type GoFile struct {
	goGenerator *Generator

	astFile *ast.File // Parsed AST.
	// Fset provides position information for Types, TypesInfo, and Syntax.
	// It is set only when Types is set.
	fset *token.FileSet

	// proto's astFile info
	protoFiles []FileInfo
	outerFile  *pluginpb.CodeGeneratorResponse_File

	fileChanged bool
}

func (g *GoFile) FoundProtoMessage(typ string) (StructInfo, bool) {
	for _, p := range g.protoFiles {
		for _, s := range p.StructInfos {
			if s.StructNameInGo == typ {
				return s, true
			}
		}
	}
	return StructInfo{}, false
}

// genDecl processes one declaration clause.
func (g *GoFile) genDecl(node ast.Node) bool {
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

		protoStruct, has := g.FoundProtoMessage(typ)
		if !has {
			// This is not the type we're looking for.
			continue
		}

		if sExpr.Fields.NumFields() <= 0 {
			g.goGenerator.protoGenerator.Error(fmt.Errorf("miss struct fields: %w", fmt.Errorf("%s has no Fields", typ)))
			continue
		}

		// Handle comment
		//if c := tspec.Comment; c != nil && len(c.List) == 1 {}

		for _, field := range sExpr.Fields.List {
			fieldName := ""
			if len(field.Names) != 0 { // pick first exported Name
				for _, field := range field.Names {
					if ast.IsExported(field.Name) {
						fieldName = field.Name
						break
					}
				}
			} else { // anonymous field
				ident, ok := field.Type.(*ast.Ident)
				if !ok {
					continue
				}
				fieldName = ident.Name
			}

			// nothing to process, continue with next line
			if fieldName == "" {
				continue
			}
			protoField, has := protoStruct.FindField(fieldName)
			if !has {
				continue
			}
			protoTags := protoField.FieldTag

			if field.Tag == nil {
				field.Tag = &ast.BasicLit{}
			}
			goTagFieldValue := field.Tag.Value

			goTags, err := reflect.ParseAstStructTag(goTagFieldValue)
			if err != nil {
				g.goGenerator.protoGenerator.Error(fmt.Errorf("malformed struct tag in field extension: %w", err))
				continue
			}

			// rewrite tags
			{
				for _, protoTag := range protoTags.Tags() {
					goTag, _ := goTags.Get(protoTag.Key)
					goTag.Key = protoTag.Key
					if protoTag.Name != "" {
						goTag.Name = protoTag.Name
					}

					switch protoField.UpdateStrategy {
					case tag.FieldTag_replace:
						goTag.Options = nil
					default:
					}
					goTag.AddOptions(protoTag.Options...)
					_ = goTags.Set(goTag)
				}

				// the order rule of struct tags is: protobuf, json, other tags in the same order written in *.proto.
				var keys = []string{"protobuf", "json"}
				keys = append(keys, strings_.SliceTrim(goTags.OrderKeys(), "protobuf", "json")...)
				newGoTag := goTags.SelectAstString(keys...)
				if newGoTag != goTagFieldValue {
					g.fileChanged = true
					field.Tag.Value = newGoTag
				}
			}
		}
	}
	return false
}
