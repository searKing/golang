// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	_ "unsafe"

	"github.com/searKing/golang/go/go/cmd/go/gopathload"
	"github.com/searKing/golang/go/go/cmd/go/modload"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

// Generator holds the state of the analysis. Primarily used to buffer
// the goimportName for format.Source.
type Generator struct {
	buf        bytes.Buffer // Accumulated goimportName.
	outputPkg  *Package
	importPkgs map[string]*Package // Package we are scanning.
	files      map[string]os.FileInfo
	seenPkgs   map[string]bool

	// config flags
	importPrefix string
	lineComment  bool
	globImport   string
	buildTag     string
}

func NewGenerator(importPrefix string, globImport string, tag string) *Generator {
	return &Generator{
		importPkgs:   map[string]*Package{},
		files:        map[string]os.FileInfo{},
		seenPkgs:     map[string]bool{},
		importPrefix: importPrefix,
		globImport:   globImport,
		buildTag:     tag,
	}
}

func (g *Generator) scanPackage(paths ...string) {
	g.parsePackageToOutput()

	for _, root := range paths {
		log.Println("scan", root)
		filepath.Walk(root, func(dir string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatalf("error: walk into %s", dir)
			}
			if !info.IsDir() {
				return nil
			}
			if g.seenPkgs[dir] {
				// Once the directory is visited, we can skip the rest of it.
				return filepath.SkipDir
			}
			g.seenPkgs[dir] = true
			log.Println("walk", dir)

			g.parsePackageToImport(dir)
			return nil
		})
	}
}

func loadPackage(patterns []string, tags ...string) *packages.Package {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		// TODO: Need to think about constants in test files. Maybe write type_string_test.go
		// in a separate pass? For later.
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}

	return pkgs[0]
}

// parsePackageToImport analyzes the single package constructed from the patterns and tags.
// parsePackageToImport exits if there is an error.
func (g *Generator) parsePackageToImport(dir string) {
	pkg := loadPackage([]string{dir}, g.buildTag)
	files, err := filepath.Glob(filepath.Join(dir, g.globImport))
	if err != nil {
		log.Fatalf("error: find glob %s under %s, got %v", g.globImport, dir, err)
	}
	if len(files) > 0 {
		g.addPackage(dir, pkg)
	}
}

func (g *Generator) parsePackageToOutput() {
	pkg := loadPackage([]string{*output}, g.buildTag)
	importPath := parseImportPath(*output)

	pkgInfo := &Package{
		name:         pkg.Name,
		importPath:   importPath,
		importPrefix: g.importPrefix,
		lineComment:  g.lineComment,
		globImport:   g.globImport,
		buildTag:     g.buildTag,
	}
	g.outputPkg = pkgInfo
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(dir string, pkg *packages.Package) {
	importPath := parseImportPath(dir)

	log.Println("import path", importPath)
	if importPath == g.outputPkg.importPath {
		// ignore self import
		return
	}
	pkgInfo := &Package{
		name:         pkg.Name,
		importPath:   importPath,
		importPrefix: g.importPrefix,
		lineComment:  g.lineComment,
		globImport:   g.globImport,
		buildTag:     g.buildTag,
	}
	g.importPkgs[dir] = pkgInfo
}
func (g *Generator) generate(toolName string, toolArgs ...string) {
	g.generateGoKeep(toolName, toolArgs...)
	g.generateGoImport(toolName, toolArgs...)
}

// generateGoKeep produces the String method for the named type.
func (g *Generator) generateGoKeep(toolName string, toolArgs ...string) {
	tmplName := "tmpl/import.tmpl"
	tmpl := importTmplProvider(tmplName)
	tmplInfo := &ImportTemplateInfo{
		GoImportToolName: toolName,
		GoImportToolArgs: toolArgs,
		BuildTag:         g.buildTag,
	}
	for dir, pkg := range g.importPkgs {
		g.buf.Reset()
		tmplInfo.ModuleName = pkg.Package()
		if err := tmpl().Execute(&g.buf, tmplInfo); err != nil {
			log.Fatalf("render template %s for %s: %v", tmplName, dir, err)
		}

		// Write to file.
		g.formatDumpTo(dir, *gokeepName+".go")
	}

}

func (g *Generator) generateGoImport(toolName string, toolArgs ...string) {
	g.buf.Reset()

	tmplName := "tmpl/import.tmpl"
	tmpl := importTmplProvider(tmplName)
	tmplInfo := &ImportTemplateInfo{
		GoImportToolName: toolName,
		GoImportToolArgs: toolArgs,
		BuildTag:         g.buildTag,
	}
	for _, pkg := range g.importPkgs {
		tmplInfo.ImportPaths = append(tmplInfo.ImportPaths, pkg.importPath)
	}

	tmplInfo.ModuleName = g.outputPkg.Package()
	if err := tmpl().Execute(&g.buf, tmplInfo); err != nil {
		log.Fatalf("render template %s for %s: %v", tmplName, *output, err)
	}

	g.formatDumpTo("", *goimportName+".go")
}

func (g *Generator) formatDumpTo(dir, file string) {
	// Format the goimportName.
	src := g.format()

	target := g.goimport(src)

	// Write to file.
	outputFile := filepath.Join(dir, strings.ToLower(file))

	err := ioutil.WriteFile(outputFile, target, 0644)
	if err != nil {
		log.Fatalf("writing goimportName: %s", err)
	}
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the goimportName to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}

func (g *Generator) goimport(src []byte) []byte {
	var opt = &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
		Fragment:  true,
	}

	res, err := imports.Process("", src, opt)
	if err != nil {
		log.Fatalf("process import: %s", err)
	}

	return res
}

func parseImportPath(dir string) string {
	var importFile func(file string) (srcdir, importPath string, err error)
	if modload.ModEnabled(dir) {
		importFile = modload.ImportFile
	} else {
		importFile = gopathload.ImportFile
	}
	var err error
	_, importPath, err := importFile(dir)
	if err != nil {
		panic(fmt.Errorf("parse package from %s, got %w", dir, err))
	}
	return importPath
}
