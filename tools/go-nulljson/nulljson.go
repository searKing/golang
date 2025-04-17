// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go-nulljson Generates Go code using a package as a generic template that implements database/sql.Scanner and database/sql/driver.Valuer.
// Given the name of a NullJson type T , and the name of a type Value
// go-nulljson will create a new self-contained Go source file implementing
//
//	func (m *T) Scan(src interface{}) error
//	func (m *T) Value() (driver.Value, error)
//
// The file is created in the same package and directory as the package that defines T, Key.
// It has helpful defaults designed for use with go generate.
//
// For example, given this snippet,
//
// running this command
//
//	go-nulljson -type=Pill<time.Time>
//
// in the same directory will create the file pill_nulljson.go, in package painkiller,
// containing a definition of
//
//	func (m *Pill) Scan(src interface{}) error
//	func (m *Pill) Value() (driver.Value, error)
//
// Typically this process would be run using go generate, like this:
//
//	//go:generate go-nulljson -type=Pill<int>
//	//go:generate go-nulljson -type=Pill<*string>
//	//go:generate go-nulljson -type=Pill<time.Time>
//	//go:generate go-nulljson -type=Pill<*encoding/json.Token>
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
// Deprecated: Use [github.com/searKing/golang/go/exp/database/sql.NullJson[T] or sql.Json[T]] instead.
// For more information, see:
// https://github.com/searKing/golang/blob/master/go/exp/database/sql/null_json.go
// https://github.com/searKing/golang/blob/master/go/exp/database/sql/json.go
//
// This package is frozen and no new functionality will be added.
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

	strings_ "github.com/searKing/golang/tools/pkg/strings"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

var (
	typeInfos   = flag.String("type", "", "comma-separated list of type names; must be set")
	output      = flag.String("output", "", "output file name; default srcdir/<type>_nulljson.go")
	trimprefix  = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	linecomment = flag.Bool("linecomment", false, "use line comment text as printed text when present")
	nullable    = flag.Bool("nullable", true, "generate nullable sql field, similar to sql.NullString")
	protojson   = flag.Bool("protojson", false, "generate codec of proto by protojson, instead of json")
	buildTags   = flag.String("tags", "", "comma-separated list of build tags to apply")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of go-nulljson:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-nulljson [flags] -type T [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-nulljson [flags] -type T<V> [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-nulljson [flags] -type T<V> files... # Must be a single package\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-nulljson [flags] -type T,S [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-nulljson [flags] -type T<V>,S<V> [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "For more information, see:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\thttps://pkg.go.dev/github.com/searKing/golang/tools/go-nulljson\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

const (
	goNullJsonToolName = "go-nulljson"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-nulljson: ")
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
		trimPrefix:  *trimprefix,
		lineComment: *linecomment,
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
	g.Printf("// Code generated by \"%s %s\"; DO NOT EDIT.\n", goNullJsonToolName, strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("// Install %[1]s by `go install github.com/searKing/golang/tools/%[1]s@latest`", goNullJsonToolName)
	g.Printf("\n")
	g.Printf("//")
	g.Printf("\n")
	g.Printf("// Deprecated: Use [github.com/searKing/golang/go/exp/database/sql.NullJson[T]] instead.")
	g.Printf("\n")
	g.Printf("// For more information, see:")
	g.Printf("\n")
	if *nullable {
		g.Printf("// https://github.com/searKing/golang/blob/master/go/exp/database/sql/null_json.go")
	} else {
		g.Printf("// https://github.com/searKing/golang/blob/master/go/exp/database/sql/json.go")
	}
	g.Printf("\n")

	g.Printf("package %s", g.pkg.name)
	g.Printf("\n")
	g.Printf(stringImport, "database/sql")
	g.Printf(stringImport, "database/sql/driver")
	g.Printf(stringImport, "encoding/json")
	g.Printf(stringImport, "fmt")
	g.Printf(stringImport, "time")
	if *protojson {
		g.Printf(stringImport, "google.golang.org/protobuf/encoding/protojson")
		g.Printf(stringImport, "google.golang.org/protobuf/proto")
	}

	// Run generate for each type.
	for _, typeInfo := range types {
		g.generate(typeInfo)
	}

	// Format the output.
	src := g.format()

	target := g.goimport(src)

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_nulljson.go", types[0].Name)
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := ioutil.WriteFile(outputName, target, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
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

// Printf format & write to the buf in this generator
func (g *Generator) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) Render(text string, arg SqlRender) {
	tmpl, err := template.New("go-nulljson").Parse(text)
	if err != nil {
		panic(err)
	}
	_ = tmpl.Execute(&g.buf, &arg)
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeInfo typeInfo
	values   []Value // Accumulator for constant values of that type.

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

// generate produces the String method for the named type.
func (g *Generator) generate(typeInfo typeInfo) {
	// <key, value>
	values := make([]Value, 0, 100)
	if typeInfo.Name != "" {
		values = append(values, Value{
			eleImport:       typeInfo.Import,
			eleName:         typeInfo.Name,
			valueImport:     typeInfo.valueImport,
			valueType:       typeInfo.valueType,
			valueIsPointer:  typeInfo.valueIsPointer,
			valueTypePrefix: typeInfo.valueTypePrefix,
		})
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %+v", typeInfo)
	}
	g.buildOneRun(values[0])
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
	eleImport string // import path of the atomic.Value type.
	eleName   string // Name of the atomic.Value type.

	valueImport     string // import path of the atomic.Value's value.
	valueType       string // The type of the value in atomic.Value.
	valueIsPointer  bool   // whether the value's type is ptr
	valueTypePrefix string // The type's prefix, such as []*[]
}

func (val *Value) ValueFullType() string {
	return strings_.LoadElse(val.valueIsPointer, "*", "") + val.valueTypePrefix + val.valueType
}

// Helpers

// declareNameVar declares the concatenated names
// strings representing the runs of values
func (g *Generator) declareNameVar(run Value) string {
	nilValName, nilValDecl := g.createValAndNameDecl(run)
	g.Printf("var %s\n", nilValDecl)
	return nilValName
}

// createValAndNameDecl returns the pair of declarations for the run. The caller will add "var".
func (g *Generator) createValAndNameDecl(val Value) (string, string) {
	goRep := strings.NewReplacer(".", "_", "{", "_", "}", "_")

	nilValName := fmt.Sprintf("_nil_%s_%s_value",
		val.eleName,
		goRep.Replace(val.valueType))
	nilValDecl := fmt.Sprintf("%s = func() (val %s) { return }()", nilValName, val.ValueFullType())

	return nilValName, nilValDecl
}

// buildOneRun generates the variables and String method for a single run of contiguous values.
func (g *Generator) buildOneRun(value Value) {
	//values := run
	g.Printf("\n")
	if strings.TrimSpace(value.eleImport) != "" {
		g.Printf(stringImport, value.eleImport)
	}
	if strings.TrimSpace(value.valueImport) != "" {
		g.Printf(stringImport, value.valueImport)
	}

	// Generate code that will fail if the constants change value.
	g.Printf("func _() {\n")
	g.Printf("\t// An \"cannot convert %s literal (type %s) to type atomic.Value\" compiler error signifies that the base type have changed.\n", value.eleName, value.eleName)
	g.Printf("\t// Re-run the go-nulljson command to generate them again.\n")
	g.Printf("\t	var val %s\n", value.eleName)
	g.Printf("\t	_ = (sql.Scanner)(&val)\n")
	g.Printf("\t	_ = (driver.Valuer)(&val)\n")
	g.Printf("}\n")

	//The generated code is simple enough to write as a Printf format.
	sqlRender := SqlRender{
		SqlJsonType: value.eleName,
		ValueType:   value.ValueFullType(),
		NilValue: strings_.LoadElseGet(value.valueIsPointer, "nil", func() string {
			return g.declareNameVar(value)
		}),
		valueImport: value.valueImport,
	}
	sqlRender.ResetCanAlias()
	var text string

	if *nullable {
		text = tmplNullJson
	} else {
		text = tmplJson
	}

	g.Render(text, sqlRender)
}

// Arguments to format are:
//
//	[1]: import path
const stringImport = `import "%s"
`
