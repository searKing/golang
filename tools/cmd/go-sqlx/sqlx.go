// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go-sqlx Generates Go code using a package as a generic template that implements sqlx.
// Given the StructName of a Struct type T
// go-sqlx will create a new self-contained Go source file and rewrite T's "db" tag of struct field
// The file is created in the same package and directory as the package that defines T.
// It has helpful defaults designed for use with go generate.
//
// For example, given this snippet,
//
// running this command
//
//	go-sqlx -type=Pill
//
// in the same directory will create the file pill_sqlx.go, in package painkiller,
// containing a definition of helper for sqlx
//
// Typically this process would be run using go generate, like this:
//
//	//go:generate go-sqlx -type=Pill
//	//go:generate go-sqlx -type=Pill --linecomment
//	//go:generate go-sqlx -type=Pill --linecomment --with-dao
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must trimmedStructName a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -type flag accepts a comma-separated list of types so a single run can
// generate methods for multiple types. The default flagOutput file is t_string.go,
// where t is the lower-cased trimmedStructName of the first type listed. It can be overridden
// with the -flagOutput flag.
//
package main // import "github.com/searKing/golang/tools/cmd/go-sqlx"

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

var (
	flagTypeInfos           = flag.String("type", "", "comma-separated list of type names; must be set")
	flagOutput              = flag.String("flagOutput", "", "flagOutput file trimmedStructName; default srcdir/<type>_sqlx.go")
	flagSkipPrivateFields   = flag.Bool("skip-unexported", false, "skip unexported Fields")
	flagSkipAnonymousFields = flag.Bool("skip-anonymous", false, "skip anonymous Fields")
	flagTrimprefix          = flag.String("trimprefix", "", "trim the `prefix` from the generated struct type names")
	flagLinecomment         = flag.Bool("linecomment", false, "use line comment text followed the generated struct type name by as printed text when present")
	flagWithDao             = flag.Bool("with-dao", false, "generate with dao")
	flagBuildTags           = flag.String("tags", "", "comma-separated list of build tags to apply")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of go-sqlx:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-sqlx [flags] -type T [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-sqlx [flags] -type T files... # Must be a single package\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-sqlx [flags] -type T,S [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "For more information, see:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\thttps://godoc.org/github.com/searKing/golang/tools/cmd/go-sqlx\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

const (
	goSqlxToolName = "go-sqlx"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-sqlx: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*flagTypeInfos) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// type <key, value> type <key, value>
	types := newTypeInfo(*flagTypeInfos)
	if len(types) == 0 {
		flag.Usage()
		os.Exit(3)
	}

	var tags []string
	if len(*flagBuildTags) > 0 {
		tags = strings.Split(*flagBuildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var dir string
	g := Generator{
		trimPrefix:  *flagTrimprefix,
		lineComment: *flagLinecomment,
	}
	// TODO(suzmue): accept other patterns for packages (directories, list of files, import paths, etc).
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
	} else {
		if len(tags) != 0 {
			log.Fatal("-tags option applies only to directories, not when files are specified")
		}
		dir = filepath.Dir(args[0])
	}

	g.parsePackage(args, tags)

	// Print the header and package clause.
	g.Printf("// Code generated by \"%s %s\"; DO NOT EDIT.\n", goSqlxToolName, strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("package %s", g.pkg.name)
	g.Printf("\n")

	// Run generate for each type.
	for _, typeInfo := range types {
		g.generate(typeInfo)
	}

	// Format the flagOutput.
	src := g.format()

	target := g.goimport(src)

	// Write to file.
	outputName := *flagOutput
	if outputName == "" {
		baseName := fmt.Sprintf("%s_sqlx.go", types[0].Name)
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := ioutil.WriteFile(outputName, target, 0644)
	if err != nil {
		log.Fatalf("writing flagOutput: %s", err)
	}
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
// the flagOutput for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated flagOutput.
	pkg *Package     // Package we are scanning.

	trimPrefix  string
	lineComment bool

	modified io.Reader // read an archive of modified files from standard input
}

// Printf format & write to the buf in this generator
func (g *Generator) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) Render(text string, arg SqlxRender) {
	tmpl, err := template.New("go-sqlx").Parse(text)
	if err != nil {
		panic(err)
	}

	err = tmpl.Funcs(template.FuncMap{"trim": strings.TrimSpace}).Execute(&g.buf, &arg)
	if err != nil {
		panic(err)
	}
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// Fset provides position information for Types, TypesInfo, and Syntax.
	// It is set only when Types is set.
	fset        *token.FileSet
	fileChanged bool

	// These Fields are reset for each type being generated.
	typeInfo typeInfo
	structs  []Struct // Accumulator for constant structs of that type.

	trimPrefix  string
	lineComment bool
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

	// Fset provides position information for Types, TypesInfo, and Syntax.
	// It is set only when Types is set.
	fset *token.FileSet
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
		fset:  pkg.Fset,
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file:        file,
			pkg:         g.pkg,
			trimPrefix:  g.trimPrefix,
			lineComment: g.lineComment,
		}
	}
}

// generate produces the String method for the named type.
func (g *Generator) generate(typeInfo typeInfo) {
	// <key, value>
	structs := make([]Struct, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeInfo = typeInfo
		file.fset = g.pkg.fset
		file.structs = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			structs = append(structs, file.structs...)
		}
		if file.fileChanged {
			var buf bytes.Buffer
			err := format.Node(&buf, g.pkg.fset, file.file)
			if err != nil {
				panic(err)
			}

			err = ioutil.WriteFile(g.pkg.fset.File(file.file.Pos()).Name(), buf.Bytes(), 0)
			if err != nil {
				panic(err)
			}
		}
	}

	if len(structs) == 0 {
		log.Fatalf("no structs defined for type %+v", typeInfo)
	}
	g.buildOneRun(structs[0])
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the flagOutput to see the error.
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

// Struct represents a declared constant.
type Struct struct {
	StructType string // The StructType of the struct.
	TableName  string // The TableName with trimmed prefix.
	Fields     []StructField
}

type StructField struct {
	FieldType string // The FieldType of the struct field.
	DbName    string
}

// Helpers

// declareNameVar declares the concatenated names
// strings representing the runs of structs
func (g *Generator) declareNameVar(run Struct) string {
	nilValName, _ := g.createValAndNameDecl(run)
	return nilValName
}

// createValAndNameDecl returns the pair of declarations for the run. The caller will add "var".
func (g *Generator) createValAndNameDecl(val Struct) (string, string) {
	goRep := strings.NewReplacer(".", "_", "{", "_", "}", "_")

	nilValName := fmt.Sprintf("_nil_%s_%s_value",
		val.StructType,
		goRep.Replace(val.TableName))
	nilValDecl := fmt.Sprintf("%s = func() (val %s) { return }()", nilValName, val.StructType)

	return nilValName, nilValDecl
}

// buildOneRun generates the variables and String method for a single run of contiguous structs.
func (g *Generator) buildOneRun(value Struct) {
	//The generated code is simple enough to write as a Printf format.
	sqlRender := SqlxRender{
		StructType: value.StructType,
		TableName:  value.TableName,
		Fields:     value.Fields,
		NilValue:   g.declareNameVar(value),
	}
	g.Render(tmplJson, sqlRender)
}

// Arguments to format are:
//	[1]: import path
const stringImport = `import "%s"
`
