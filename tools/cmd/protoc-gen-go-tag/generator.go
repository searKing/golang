// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/golang/protobuf/protoc-gen-go/generator"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	fset = token.NewFileSet()
)

type Generator struct {
	// Ast goFiles to which this package contains.
	goFiles []GoFile

	// proto's astFile info
	protoFiles []FileInfo

	protoGenerator *generator.Generator
}

func NewGenerator(protoFiles []FileInfo) *Generator {
	return &Generator{
		protoFiles: protoFiles,
	}
}

// ParseGoContent analyzes the single package constructed from the patterns and tags.
// ParseGoContent exits if there is an error.
func (g *Generator) ParseGoContent(outerFile *pluginpb.CodeGeneratorResponse_File) {
	if outerFile == nil || outerFile.GetContent() == "" {
		return
	}
	const mode = parser.AllErrors | parser.ParseComments

	f, err := parser.ParseFile(fset, "", outerFile.GetContent(), mode)
	if err != nil {
		g.protoGenerator.Error(err, "failed to parse struct tag in field extension")
	}
	g.addGoFile(f, outerFile)
}

// addGoFile adds a type checked Package and its syntax goFiles to the generator.
func (g *Generator) addGoFile(astFile *ast.File, outerFile *pluginpb.CodeGeneratorResponse_File) {
	g.goFiles = append(g.goFiles, GoFile{
		goGenerator: g,
		astFile:     astFile,
		fset:        fset,
		protoFiles:  g.protoFiles,
		outerFile:   outerFile,
	})
}

// Generate produces the rewrite content to proto's Generator.
func (g *Generator) Generate() {
	for _, file := range g.goFiles {
		// Set the state for this run of the walker.
		if file.astFile != nil {
			ast.Inspect(file.astFile, file.genDecl)
		}

		if file.fileChanged {
			var buf bytes.Buffer
			err := format.Node(&buf, file.fset, file.astFile)
			if err != nil {
				g.protoGenerator.Error(err, "failed to format go content")
			}

			content := buf.String()
			file.outerFile.Content = &content
		}
	}
}
