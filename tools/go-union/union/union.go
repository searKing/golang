// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package union

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	slices_ "github.com/searKing/golang/go/exp/slices"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

var (
	typeInfos               = flag.String("type", "", "comma-separated list of type names; must be set")
	output                  = flag.String("output", "", "output file name; default srcdir/<type>_union.go")
	flagSkipPrivateFields   = flag.Bool("skip-unexported", false, "skip unexported Fields")
	flagSkipAnonymousFields = flag.Bool("skip-anonymous", false, "skip anonymous Fields")
	buildTags               = flag.String("tags", "", "comma-separated list of build tags to apply")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of go-union:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-union [flags] -type T [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "For more information, see:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\thttps://pkg.go.dev/github.com/searKing/golang/tools/go-union\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

const (
	goUnionsToolName = "go-union"
)

func Main() {
	log.SetFlags(0)
	log.SetPrefix("go-union: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeInfos) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// type <key, value> type <key, value>
	typs := newTypeInfo(*typeInfos)
	if len(typs) == 0 {
		flag.Usage()
		os.Exit(3)
	}

	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var dir string
	var g Generator
	// TODO(suzmue): accept other patterns for packages (directories, list of files, import paths, etc).
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags union applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	g.parsePackage(args, tags)

	var structs []Struct
	// Run inspect for each type.
	for _, typeInfo := range typs {
		structs = append(structs, g.inspect(typeInfo))
	}

	// Run render for each type.
	g.generateUnion(dir, structs...)
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
func (g *Generator) Reset() {
	g.buf.Reset()
}

// Printf format & write to the buf in this generator
func (g *Generator) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(&g.buf, format, args...)
}

func IsExported(arg reflect.Value) bool {
	return ast.IsExported(arg.String())
}

func (g *Generator) Render(text string, arg any) {
	tmpl, err := template.New("go-union").Funcs(template.FuncMap{"IsExported": IsExported}).Parse(text)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(&g.buf, arg)
	if err != nil {
		panic(err)
	}
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeInfo typeInfo
	structs  []Struct // Accumulator for constant structs of that type.
}

// Package holds a single parsed package and associated files and ast files.
type Package struct {
	// Name is the package name as it appears in the package source code.
	name string

	// Defs maps identifiers to the objects they define (including
	// package names, dots "." of dot-imports, and blank "_" identifiers).
	// For identifiers that do not denote objects (e.g., the package name
	// in package clauses, or symbolic variables t in t := x.(type) of
	// type switch headers), the corresponding objects are nil.
	//
	// For an embedded field, Defs returns the field *Var it defines.
	//
	// Invariant: Defs[id] == nil || Defs[id].Pos() == id.Pos()
	defs map[*ast.Ident]types.Object

	// Ast files to which this package contains.
	files []*File
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string, tags []string) {
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
	g.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file: file,
			pkg:  g.pkg,
		}
	}
}

// inspect scans files for the named type.
func (g *Generator) inspect(typeInfo typeInfo) Struct {
	// <key, value>
	structs := make([]Struct, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeInfo = typeInfo
		file.structs = nil
		if file.file != nil {
			importExps := slices_.MapFunc(file.file.Imports, func(e *ast.ImportSpec) string {
				if e.Path == nil {
					return ""
				}
				if e.Name == nil || e.Name.String() == "_" {
					return e.Path.Value
				}
				return e.Name.String() + " " + e.Path.Value
			})
			importExps = slices_.Filter(importExps)

			ast.Inspect(file.file, file.genDecl)
			for i := range file.structs {
				file.structs[i].FileImports = importExps
			}

			structs = append(structs, file.structs...)
		}
	}

	if len(structs) == 0 {
		log.Fatalf("no values defined for type %+v", typeInfo)
	}
	return structs[0]
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := imports.Process("", g.buf.Bytes(), nil)
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
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

// Value represents a declared constant.
type Value struct {
	typName        string // The name of the constant.
	trimmedTypName string // The name with trimmed prefix.
	str            string // The string representation given by the "go/constant" package.

	typImport string // import path of the atomic.Value type.
}

func (v *Value) String() string {
	return v.str
}

// Helpers

// createValAndNameDecl returns the pair of declarations for the run. The caller will add "var".
func createValAndNameDecl(typ string) (string, string) {
	defaultValName := fmt.Sprintf("_default_%s_value", typ)
	defaultValDecl := fmt.Sprintf("%s = func() (val %s) { return }()", defaultValName, typ)

	return defaultValName, defaultValDecl
}

func (g *Generator) generateUnion(dir string, structs ...Struct) {
	for _, val := range structs {
		g.generateUnionOneRun(dir, val)
	}
}

// generateUnionOneRun produces the Union method for the named type.
func (g *Generator) generateUnionOneRun(dir string, value Struct) {
	//The generated code is simple enough to write as a Printf format.
	tmplRender := TmplUnionRender{
		GoUnionToolName:              goUnionsToolName,
		GoUnionToolArgs:              os.Args[1:],
		PackageName:                  g.pkg.name,
		ImportPaths:                  value.FileImports,
		TargetTypeName:               value.StructTypeName,
		TargetTypeImport:             value.StructTypeImport,
		TargetTypeGenericDeclaration: value.StructTypeGenericDeclaration,
		TargetTypeGenericParams:      value.StructTypeGenericTypeParams,
		Fields:                       value.Fields,
		ApplyUnionsAsMemberFunction:  false,
	}

	tmplRender.Complete()
	g.Reset()
	g.Render(tmplUnion, tmplRender)

	// Format the output.
	src := g.format()

	target := g.goimport(src)

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_union.go", value.StructTypeName)
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := os.WriteFile(outputName, target, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
