// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go-option Generates Go code using a package as a graceful options.
// Given the name of a type T
// go-option will create a new self-contained Go source file implementing
//	func apply(*Pill)
// The file is created in the same package and directory as the package that defines T.
// It has helpful defaults designed for use with go generate.
//
// For example, given this snippet,
//
//	package painkiller
//
//
//	type Pill struct{}
//
//
// running this command
//
//	go-option -type=Pill
//
// in the same directory will create the file pill_options.go, in package painkiller,
// containing a definition of
//
//	var _default_Pill_value = func() (val Pill) { return }()
//
//	// A PillOptions sets options.
//	type PillOptions interface {
//		apply(*Pill)
//	}
//
//	// EmptyPillOptions does not alter the configuration. It can be embedded
//	// in another structure to build custom options.
//	//
//	// This API is EXPERIMENTAL.
//	type EmptyPillOptions struct{}
//
//	func (EmptyPillOptions) apply(*Pill) {}
//
//	// PillOptionFunc wraps a function that modifies PillOptionFunc into an
//	// implementation of the PillOptions interface.
//	type PillOptionFunc func(*Number)
//
//	func (f PillOptionFunc) apply(do *Pill) {
//		f(do)
//	}
//
//	func (o *Pill) ApplyOptions(options ...PillOption) *Pill {
//		for _, opt := range options {
//			if opt == nil {
//				continue
//			}
//			opt.apply(o)
//		}
//		return o
//	}

//
// Typically this process would be run using go generate, like this:
//
//	//go:generate go-option -type=Pill
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must name a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -type flag accepts a comma-separated list of types so a single run can
// generate methods for multiple types. The default output file is t_string.go,
// where t is the lower-cased name of the first type listed. It can be overridden
// with the -output flag.
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
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
	typeInfos               = flag.String("type", "", "comma-separated list of type names; must be set")
	output                  = flag.String("output", "", "output file name; default srcdir/<type>_option.go")
	flagSkipPrivateFields   = flag.Bool("skip-unexported", false, "skip unexported Fields")
	flagSkipAnonymousFields = flag.Bool("skip-anonymous", false, "skip anonymous Fields")
	trimPrefix              = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	trim                    = flag.Bool("trim", false, "trim type names as prefix from the generated constant names")
	lineComment             = flag.Bool("linecomment", false, "use line comment text as printed text when present")
	buildTags               = flag.String("tags", "", "comma-separated list of build tags to apply")
	option                  = flag.Bool("option", true, "generate options for type names")
	config                  = flag.Bool("config", false, "generate completed config for type names")
	optionOnly              = flag.Bool("optiononly", false, "generate option, mute config; overwrite flags --config and --option; --optionOnly and --configOnly can not both be set")
	configOnly              = flag.Bool("configonly", false, "generate config, mute option; overwrite flags --config and --option; --optionOnly and --configOnly can not both be set")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of go-option:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-option [flags] -type T [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "For more information, see:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\thttps://godoc.org/github.com/searKing/golang/tools/go-option\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

const (
	goOptionsToolName = "go-option"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-option: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeInfos) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	// type <key, value> type <key, value>
	types := newTypeInfo(*typeInfos)
	if len(types) == 0 {
		flag.Usage()
		os.Exit(3)
	}

	if *optionOnly && *configOnly {
		flag.Usage()
		os.Exit(4)
	}

	if *optionOnly {
		*option = true
		*config = false
	}
	if *configOnly {
		*option = false
		*config = true
	}

	if !*option && !*config {
		log.Print("no op has been applied")
		return
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
	g := Generator{
		trimPrefix:  *trimPrefix,
		lineComment: *lineComment,
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
	//
	//if *code {
	//	// Run render for each type.
	//	g.generateConfig(dir, types...)
	//}

	var structs []Struct
	// Run inspect for each type.
	for _, typeInfo := range types {
		structs = append(structs, g.inspect(typeInfo))
	}

	// Run render for each type.
	if *option {
		g.generateOption(dir, structs...)
	}
	if *config {
		g.generateConfig(dir, structs...)
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
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.

	trimPrefix  string
	lineComment bool
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
func (g *Generator) Reset() {
	g.buf.Reset()
}

// Printf format & write to the buf in this generator
func (g *Generator) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) Render(text string, arg interface{}) {
	tmpl, err := template.New("go-option").Parse(text)
	if err != nil {
		panic(err)
	}
	_ = tmpl.Execute(&g.buf, arg)
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
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
			file:        file,
			pkg:         g.pkg,
			trimPrefix:  g.trimPrefix,
			lineComment: g.lineComment,
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
			ast.Inspect(file.file, file.genDecl)
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
	src, err := format.Source(g.buf.Bytes())
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

func (g *Generator) generateOption(dir string, structs ...Struct) {
	for _, val := range structs {
		g.generateOptionOneRun(dir, val)
	}
}

// generateOptionOneRun produces the Option method for the named type.
func (g *Generator) generateOptionOneRun(dir string, value Struct) {
	//The generated code is simple enough to write as a Printf format.
	tmplRender := TmplOptionRender{
		GoOptionToolName:             goOptionsToolName,
		GoOptionToolArgs:             os.Args[1:],
		PackageName:                  g.pkg.name,
		TargetTypeName:               value.StructTypeName,
		TargetTypeImport:             value.StructTypeImport,
		TrimmedTypeName:              value.trimmedStructTypeName,
		Fields:                       value.Fields,
		ApplyOptionsAsMemberFunction: false,
	}

	tmplRender.Complete()
	g.Reset()
	g.Render(tmplOption, tmplRender)

	// Format the output.
	src := g.format()

	target := g.goimport(src)

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_options.go", value.StructTypeName)
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := ioutil.WriteFile(outputName, target, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func (g *Generator) generateConfig(dir string, structs ...Struct) {
	for _, val := range structs {
		g.generateConfigOneRun(dir, val)
	}
}

// generateOptionOneRun produces the Option method for the named type.
func (g *Generator) generateConfigOneRun(dir string, value Struct) {
	//The generated code is simple enough to write as a Printf format.
	tmplRender := TmplConfigRender{
		GoOptionToolName:             goOptionsToolName,
		GoOptionToolArgs:             os.Args[1:],
		PackageName:                  g.pkg.name,
		TargetTypeName:               value.StructTypeName,
		TargetTypeImport:             value.StructTypeImport,
		TrimmedTypeName:              value.trimmedStructTypeName,
		ApplyOptionsAsMemberFunction: false,
	}

	tmplRender.Complete()
	g.Reset()
	g.Render(tmplConfig, tmplRender)

	// Format the output.
	src := g.format()

	target := g.goimport(src)

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s.config.go", value.StructTypeName)
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	if _, err := os.Stat(outputName); !os.IsNotExist(err) {
		actual, err := ioutil.ReadFile(outputName)
		if err != nil || (len(actual) > 0 && bytes.Compare(actual, target) != 0) {
			log.Fatalf("%[1]s already exists, remove or truncate(0) it before generate, as %[1]s will be overwritten", outputName)
		}
	}
	err := ioutil.WriteFile(outputName, target, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
