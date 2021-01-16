// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go-enum is a tool to automate the creation of methods that satisfy such interfaces:
// 	fmt         ==>  fmt.Stringer
// 	binary      ==>  encoding.BinaryMarshaler and encoding.BinaryUnmarshaler
// 	json        ==>  encoding/json.MarshalJSON and encoding/json.UnmarshalJSON
// 	text        ==>  encoding.TextMarshaler and encoding.TextUnmarshaler
// 	sql         ==>  database/sql.Scanner and database/sql/driver.Valuer
// 	yaml        ==>  gopkg.in/yaml.v2:yaml.Marshaler and gopkg.in/yaml.v2:yaml.Unmarshaler
//
// Given the name of a (signed or unsigned) integer type T that has constants
// defined, stringer will create a new self-contained Go source file implementing
// 	fmt         ==>  fmt.Stringer
//		func (t T) String() string
// 	binary      ==>  encoding.BinaryMarshaler and encoding.BinaryUnmarshaler
//		func (t T) MarshalBinary() (data []byte, err error)
//		func (t *T) UnmarshalBinary(data []byte) error
// 	json        ==>  encoding/json.MarshalJSON and encoding/json.UnmarshalJSON
//		func (t T) MarshalJSON() ([]byte, error)
//		func (t *T) UnmarshalJSON(data []byte) error
// 	text        ==>  encoding.TextMarshaler and encoding.TextUnmarshaler
//		func (t T) MarshalText() ([]byte, error)
//		func (t *T) UnmarshalText(text []byte) error
// 	sql         ==>  database/sql.Scanner and database/sql/driver.Valuer
//		func (t T) Value() (driver.Value, error)
//		func (t *T) Scan(value interface{}) error
// 	yaml        ==>  gopkg.in/yaml.v2:yaml.Marshaler and gopkg.in/yaml.v2:yaml.Unmarshaler
//		func (t T) MarshalYAML() (interface{}, error)
//		func (t *T) UnmarshalYAML(unmarshal func(interface{}) error) error
//
// The file is created in the same package and directory as the package that defines T.
// It has helpful defaults designed for use with go generate.
//
// go-enum works best with constants that are consecutive values such as created using iota,
// but creates good code regardless. In the future it might also provide custom support for
// constant sets that are bit patterns.
//
// For example, given this snippet,
//
//	package painkiller
//
//	type Pill int
//
//	const (
//		Placebo Pill = iota
//		Aspirin
//		Ibuprofen
//		Paracetamol
//		Acetaminophen = Paracetamol
//	)
//
// running this command
//
//	go-enum -type=Pill
//
// in the same directory will create the file pill_string.go, in package painkiller,
// containing a definition of interfaces mentioned.
//
// That method will translate the value of a Pill constant to the string representation
// of the respective constant name, so that the call fmt.Print(painkiller.Aspirin) will
// print the string "Aspirin".
//
// Typically this process would be run using go generate, like this:
//
//	//go:generate go-enum -type=Pill
//
// If multiple constants have the same value, the lexically first matching name will
// be used (in the example, Acetaminophen will print as "Paracetamol").
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
// The -linecomment flag tells stringer to generate the text of any line comment, trimmed
// of leading spaces, instead of the constant name. For instance, if the constants above had a
// Pill prefix, one could write
//   PillAspirin // Aspirin
// to suppress it in the output.
package main // import "github.com/searKing/golang/tools/cmd/go-enum"

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/format"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

var (
	typeInfos string
	useAll    bool
	useNew    bool
	useString bool
	useBinary bool
	useText   bool
	useJson   bool
	useSql    bool
	useYaml   bool

	useContains     bool
	transformMethod string
	output          string
	trimprefix      string
	linecomment     bool
	buildTags       string
)

func ParseCommandLine(def bool) *flag.FlagSet {
	var commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	commandLine.StringVar(&typeInfos, "type", "", "comma-separated list of type names; must be set")
	commandLine.BoolVar(&useAll, "all", def, "if true, all interfaces will be implemented default. Default: true.")
	commandLine.BoolVar(&useNew, "new", def, "if true, the New will be implemented. Default: true.")
	commandLine.BoolVar(&useString, "string", def, "if true, the fmt.Stringer interface will be implemented. Default: true, you can use stringer instead.")
	commandLine.BoolVar(&useBinary, "binary", def, "if true, the encoding.BinaryMarshaler and encoding.BinaryUnmarshaler interface will be implemented. Default: true")
	commandLine.BoolVar(&useText, "text", def, "if true, the encoding.TextMarshaler and encoding.TextUnmarshaler interface will be implemented. Default: true")
	commandLine.BoolVar(&useJson, "json", def, "if true, the encoding/json.Marshaler and encoding/json.Unmarshaler interface will be implemented. Default: true")
	commandLine.BoolVar(&useSql, "sql", def, "if true, the database/sql.Scanner and database/sql/driver.Valuer interface will be implemented. Default: true")
	commandLine.BoolVar(&useYaml, "yaml", def, "if true, the gopkg.in/yaml.v2:yaml.Marshaler and gopkg.in/yaml.v2:yaml.Unmarshaler interface will be implemented. Default: true")

	commandLine.BoolVar(&useContains, "contains", def, "if true, the XXXSliceContains|XXXSliceContainsAny methods will be generated(XXX will be replaced by typename), such as strings.Contains|ContainsAny. Default: true")

	commandLine.StringVar(&transformMethod, "transform", "nop", "enum item name transformation method [nop, upper, lower, snake, upper_camel, lower_camel, kebab, dotted]. Default: nop")

	commandLine.StringVar(&output, "output", "", "output file name; default srcdir/<type>_enum.go")
	commandLine.StringVar(&trimprefix, "trimprefix", "", "trim the `prefix` from the generated constant names")
	commandLine.BoolVar(&linecomment, "linecomment", false, "use line comment text as printed text when present")
	commandLine.StringVar(&buildTags, "tags", "", "comma-separated list of build tags to apply")
	// Usage is a replacement usage function for the flags package.
	commandLine.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of go-enum:\n")
		_, _ = fmt.Fprintf(os.Stderr, "\tgo-enum [flags] -type T [directory]\n")
		_, _ = fmt.Fprintf(os.Stderr, "\tgo-enum [flags] -type T files... # Must be a single package\n")
		_, _ = fmt.Fprintf(os.Stderr, "\tgo-enum [flags] -type T,S [directory]\n")
		_, _ = fmt.Fprintf(os.Stderr, "For more information, see:\n")
		_, _ = fmt.Fprintf(os.Stderr, "\thttps://godoc.org/github.com/searKing/golang/tools/cmd/go-enum\n")
		_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
		commandLine.PrintDefaults()
	}

	commandLine.Parse(os.Args[1:])
	return commandLine
}

const (
	goJsonEnumToolName = "go-enum"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-enum: ")
	ParseCommandLine(true)
	if len(typeInfos) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if !useAll {
		ParseCommandLine(false)
	}

	// type <key, value> type <key, value>
	types := newTypeInfo(typeInfos)
	if len(types) == 0 {
		flag.Usage()
		os.Exit(3)
	}

	var tags []string
	if len(buildTags) > 0 {
		tags = strings.Split(buildTags, ",")
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
		trimPrefix:  trimprefix,
		lineComment: linecomment,
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
	g.Printf("// Code generated by \"%s %s\"; DO NOT EDIT.\n", goJsonEnumToolName, strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("package %s", g.pkg.name)
	g.Printf("\n")

	// Run generate for each type.
	for _, typeInfo := range types {
		g.generate(typeInfo)
	}

	// Format the output.
	src := g.format()

	target := g.goimport(src)

	// Write to file.
	outputName := output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_enum.go", types[0].Name)
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
func (g *Generator) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(&g.buf, format, args...)
}

// byValue lets us sort the constants into increasing order.
// We take care in the Less method to sort in signed or unsigned order,
// as appropriate.
type byValue []Value

func (b byValue) Len() int      { return len(b) }
func (b byValue) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byValue) Less(i, j int) bool {
	if b[i].valueInfo.signed {
		return int64(b[i].valueInfo.value) < int64(b[j].valueInfo.value)
	}
	return b[i].valueInfo.value < b[j].valueInfo.value
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeInfo typeInfo
	values   []Value // Accumulator for constant values of that type.

	trimPrefix  string
	lineComment bool // use line comment text as printed text when present
}

// Package holds a single parsed package and associated files and ast files.
type Package struct {
	// Name is the package trimmedTypeName as it appears in the package source code.
	name string

	// Defs maps identifiers to the objects they define (including
	// package names, dots "." of dot-imports, and blank "_" identifiers).
	// For identifiers that do not denote objects (e.g., the package trimmedTypeName
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

func (g *Generator) buildEnumRegenerateCheck(values []Value) {
	g.Printf("func _() {\n")
	g.Printf("\t// An \"invalid array index\" compiler error signifies that the constant values have changed.\n")
	g.Printf("\t// Re-run the stringer command to generate them again.\n")
	g.Printf("\tvar x [1]struct{}\n")
	for _, v := range values {
		g.Printf("\t_ = x[%s - %s]\n", v.nameInfo.originalName, v.valueInfo.str)
	}
	g.Printf("}\n")
}

// generate produces the String method for the named type.
func (g *Generator) generate(typeInfo typeInfo) {
	// <key, value>
	values := make([]Value, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeInfo = typeInfo
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %+v", typeInfo)
	}
	g.transformValueNames(values, transformMethod)
	// Generate code that will fail if the constants change value.
	for _, im := range checkImportPackages {
		g.Printf(stringImport, im)
	}

	if useNew {
		for _, im := range newImportPackages {
			g.Printf(stringImport, im)
		}
	}
	if useBinary {
		for _, im := range binaryImportPackages {
			g.Printf(stringImport, im)
		}
	}
	if useJson {
		for _, im := range jsonImportPackages {
			g.Printf(stringImport, im)
		}
	}
	if useText {
		for _, im := range textImportPackages {
			g.Printf(stringImport, im)
		}
	}
	if useYaml {
		for _, im := range yamlImportPackages {
			g.Printf(stringImport, im)
		}
	}
	if useSql {
		for _, im := range sqlImportPackages {
			g.Printf(stringImport, im)
		}
	}

	g.buildEnumRegenerateCheck(values)

	runs := splitIntoRuns(values)
	threshold := 10

	if useString {
		// The decision of which pattern to use depends on the number of
		// runs in the numbers. If there's only one, it's easy. For more than
		// one, there's a tradeoff between complexity and size of the data
		// and code vs. the simplicity of a map. A map takes more space,
		// but so does the code. The decision here (crossover at 10) is
		// arbitrary, but considers that for large numbers of runs the cost
		// of the linear scan in the switch might become important, and
		// rather than use yet another algorithm such as binary search,
		// we punt and use a map. In any case, the likelihood of a map
		// being necessary for any realistic example other than bitmasks
		// is very low. And bitmasks probably deserve their own analysis,
		// to be done some other day.
		switch {
		case len(runs) == 1:
			g.buildOneRun(runs, typeInfo)
		case len(runs) <= threshold:
			g.buildMultipleRuns(runs, typeInfo)
		default:
			g.buildMap(runs, typeInfo)
		}
	}

	if useNew {
		g.Printf(newTemplate, typeInfo.Name)
	}
	if useBinary {
		g.buildCheck(runs, typeInfo.Name, threshold)
		g.Printf(binaryTemplate, typeInfo.Name)
	}
	if useJson {
		g.buildCheck(runs, typeInfo.Name, threshold)
		g.Printf(jsonTemplate, typeInfo.Name)
	}
	if useText {
		g.buildCheck(runs, typeInfo.Name, threshold)
		g.Printf(textTemplate, typeInfo.Name)
	}
	if useYaml {
		g.buildCheck(runs, typeInfo.Name, threshold)
		g.Printf(yamlTemplate, typeInfo.Name)
	}
	if useSql {
		g.buildCheck(runs, typeInfo.Name, threshold)
		g.Printf(sqpTemplate, typeInfo.Name)
	}

	if useContains {
		g.Printf(containsTemplate, typeInfo.Name)
	}
}

// splitIntoRuns breaks the values into runs of contiguous sequences.
// For example, given 1,2,3,5,6,7 it returns {1,2,3},{5,6,7}.
// The input slice is known to be non-empty.
func splitIntoRuns(values []Value) [][]Value {
	// We use stable sort so the lexically first name is chosen for equal elements.
	sort.Stable(byValue(values))
	// Remove duplicates. Stable sort has put the one we want to print first,
	// so use that one. The String method won't care about which named constant
	// was the argument, so the first name for the given value is the only one to keep.
	// We need to do this because identical values would cause the switch or map
	// to fail to compile.
	j := 1
	for i := 1; i < len(values); i++ {
		if values[i].valueInfo.value != values[i-1].valueInfo.value {
			values[j] = values[i]
			j++
		}
	}
	values = values[:j]
	runs := make([][]Value, 0, 10)
	for len(values) > 0 {
		// One contiguous sequence per outer loop.
		i := 1
		for i < len(values) && values[i].valueInfo.value == values[i-1].valueInfo.value+1 {
			i++
		}
		runs = append(runs, values[:i])
		values = values[i:]
	}
	return runs
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

type NameInfo struct {
	originalName string // The name of the constant.
	trimmedName  string // The name with trimmed prefix.
}
type TypeInfo struct {
	originalName string // The name of the constant.
}
type ValueInfo struct {
	// The value is stored as a bit pattern alone. The boolean tells us
	// whether to interpret it as an int64 or a uint64; the only place
	// this matters is when sorting.
	// Much of the time the str field is all we need; it is printed
	// by Value.String.
	value  uint64 // Will be converted to int64 when needed.
	signed bool   // Whether the constant is a signed type.
	str    string // The string representation given by the "go/constant" package.
}

// Value represents a declared constant.
// const NAME TYPE = VALUE
type Value struct {
	// name
	nameInfo NameInfo
	// type
	typeInfo TypeInfo
	// value
	valueInfo ValueInfo
	// comment
	comment string
}

func (v *Value) String() string {
	return v.valueInfo.str
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	// Token must be in IMPORT, CONST, TYPE, VAR
	if !ok || decl.Tok != token.CONST {
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
		vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		if vspec.Type == nil && len(vspec.Values) > 0 {
			// "X = 1". With no type but a value, the constant is untyped.
			// Skip this vspec and reset the remembered type.
			typ = ""

			// If this is a simple type conversion, remember the type.
			// We don't mind if this is actually a call; a qualified call won't
			// be matched (that will be SelectorExpr, not Ident), and only unusual
			// situations will result in a function call that appears to be
			// a type conversion.
			// such as a.(b)
			ce, ok := vspec.Values[0].(*ast.CallExpr)
			if !ok {
				continue
			}
			id, ok := ce.Fun.(*ast.Ident)
			if !ok {
				continue
			}
			typ = id.Name
		}
		if vspec.Type != nil {
			// "X T". We have a type. Remember it.
			ident, ok := vspec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			typ = ident.Name
		}
		if typ != f.typeInfo.Name {
			// This is not the type we're looking for.
			continue
		}

		// We now have a list of names (from one line of source code) all being
		// declared with the desired type.
		// Grab their names and actual values and store them in f.values.
		for _, name := range vspec.Names {
			if name.Name == "_" {
				continue
			}

			// This dance lets the type checker find the values for us. It's a
			// bit tricky: look up the object declared by the trimmedTypeName, find its
			// types.Const, and extract its value.
			obj, ok := f.pkg.defs[name]
			if !ok {
				log.Fatalf("no value for constant %s", name)
			}
			info := obj.Type().Underlying().(*types.Basic).Info()
			if info&types.IsInteger == 0 {
				log.Fatalf("can't handle non-integer constant type %s", typ)
				return false
			}

			value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
			if value.Kind() != constant.Int {
				log.Fatalf("can't happen: constant is not an integer %s", name)
				return false
			}
			i64, isInt := constant.Int64Val(value)
			u64, isUint := constant.Uint64Val(value)
			if !isInt && !isUint {
				log.Fatalf("internal error: value of %s is not an integer: %s", name, value.String())
			}
			if !isInt {
				u64 = uint64(i64)
			}

			v := Value{
				nameInfo: NameInfo{
					originalName: name.Name,
				},
				typeInfo: TypeInfo{
					originalName: typ,
				},
				valueInfo: ValueInfo{
					value:  u64,
					signed: info&types.IsUnsigned == 0,
					str:    value.String(),
				},
			}
			if c := vspec.Comment; f.lineComment && c != nil {
				v.comment = c.Text()
			}
			if c := vspec.Comment; f.lineComment && c != nil && len(c.List) == 1 {
				v.nameInfo.trimmedName = strings.TrimSpace(c.Text())
			} else {
				v.nameInfo.trimmedName = strings.TrimPrefix(v.nameInfo.originalName, f.trimPrefix)
			}
			f.values = append(f.values, v)
		}
	}
	return false
}

// Helpers

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	default:
		// 2^32 is enough constants for anyone.
		return 32
	}
}

// declareIndexAndNameVars declares the index slices and concatenated names
// strings representing the runs of values.
func (g *Generator) declareIndexAndNameVars(runs [][]Value, typeName string) {
	var indexes, names []string
	for i, run := range runs {
		index, name := g.createIndexAndNameDecl(run, typeName, fmt.Sprintf("_%d", i))
		if len(run) != 1 {
			indexes = append(indexes, index)
		}
		names = append(names, name)
	}
	g.Printf("const (\n")
	for _, name := range names {
		g.Printf("\t%s\n", name)
	}
	g.Printf(")\n\n")

	if len(indexes) > 0 {
		g.Printf("var (")
		for _, index := range indexes {
			g.Printf("\t%s\n", index)
		}
		g.Printf(")\n\n")
	}
}

// declareIndexAndNameVar is the single-run version of declareIndexAndNameVars
func (g *Generator) declareIndexAndNameVar(run []Value, typeName string) {
	index, name := g.createIndexAndNameDecl(run, typeName, "")
	g.Printf("const %s\n", name)
	g.Printf("var %s\n", index)
}

// createIndexAndNameDecl returns the pair of declarations for the run. The caller will add "const" and "var".
func (g *Generator) createIndexAndNameDecl(run []Value, typeName string, suffix string) (string, string) {
	b := new(bytes.Buffer)
	indexes := make([]int, len(run))
	for i := range run {
		b.WriteString(run[i].nameInfo.trimmedName)
		indexes[i] = b.Len()
	}
	nameConst := fmt.Sprintf("_%s_name%s = %q", typeName, suffix, b.String())
	nameLen := b.Len()
	b.Reset()
	_, _ = fmt.Fprintf(b, "_%s_index%s = [...]uint%d{0, ", typeName, suffix, usize(nameLen))
	for i, v := range indexes {
		if i > 0 {
			_, _ = fmt.Fprintf(b, ", ")
		}
		_, _ = fmt.Fprintf(b, "%d", v)
	}
	_, _ = fmt.Fprintf(b, "}")
	return b.String(), nameConst
}

// declareNameVars declares the concatenated names string representing all the values in the runs.
func (g *Generator) declareNameVars(runs [][]Value, typeName string, suffix string) {
	g.Printf("const _%s_name%s = \"", typeName, suffix)
	for _, run := range runs {
		for i := range run {
			g.Printf("%s", run[i].typeInfo.originalName)
		}
	}
	g.Printf("\"\n")
}

// buildOneRun generates the variables and String method for a single run of contiguous values.
func (g *Generator) buildOneRun(runs [][]Value, typeInfo typeInfo) {
	values := runs[0]
	typeName := typeInfo.Name
	g.Printf("\n")
	g.declareIndexAndNameVar(values, typeName)
	// The generated code is simple enough to write as a Printf format.
	lessThanZero := ""
	if values[0].valueInfo.signed {
		lessThanZero = "i < 0 || "
	}
	if values[0].valueInfo.value == 0 { // Signed or unsigned, 0 is still 0.
		g.Printf(stringOneRun, typeName, usize(len(values)), lessThanZero)
	} else {
		g.Printf(stringOneRunWithOffset, typeName, values[0].String(), usize(len(values)), lessThanZero)
	}
}

// buildMultipleRuns generates the variables and String method for multiple runs of contiguous values.
// For this pattern, a single Printf format won't do.
func (g *Generator) buildMultipleRuns(runs [][]Value, typeInfo typeInfo) {
	typeName := typeInfo.Name
	g.Printf("\n")
	g.declareIndexAndNameVars(runs, typeName)
	g.Printf("func (i %s) String() string {\n", typeName)
	g.Printf("\tswitch {\n")
	for i, values := range runs {
		if len(values) == 1 {
			g.Printf("\tcase i == %s:\n", &values[0])
			g.Printf("\t\treturn _%s_name_%d\n", typeName, i)
			continue
		}
		g.Printf("\tcase %s <= i && i <= %s:\n", &values[0], &values[len(values)-1])
		if values[0].valueInfo.value != 0 {
			g.Printf("\t\ti -= %s\n", &values[0])
		}
		g.Printf("\t\treturn _%s_name_%d[_%s_index_%d[i]:_%s_index_%d[i+1]]\n",
			typeName, i, typeName, i, typeName, i)
	}
	g.Printf("\tdefault:\n")
	g.Printf("\t\treturn \"%s(\" + strconv.FormatInt(int64(i), 10) + \")\"\n", typeName)
	g.Printf("\t}\n")
	g.Printf("}\n")
}

// buildMap handles the case where the space is so sparse a map is a reasonable fallback.
// It's a rare situation but has simple code.
func (g *Generator) buildMap(runs [][]Value, typeInfo typeInfo) {
	typeName := typeInfo.Name
	g.Printf("\n")
	g.declareNameVars(runs, typeName, "")
	g.Printf("\nvar _%s_map = map[%s]string{\n", typeName, typeName)
	n := 0
	for _, values := range runs {
		for _, value := range values {
			g.Printf("\t%s: _%s_name[%d:%d],\n", &value, typeName, n, n+len(value.typeInfo.originalName))
			n += len(value.typeInfo.originalName)
		}
	}
	g.Printf("}\n\n")
	g.Printf(stringMap, typeName)
}
