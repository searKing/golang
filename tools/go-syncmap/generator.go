package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"

	strings_ "github.com/searKing/golang/tools/pkg/strings"
	"golang.org/x/tools/imports"
)

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

// Render generates the package, import, variables and String method for contiguous values.
func (g *Generator) Render(text string, arg TmplPackageRender) {
	tmpl, err := template.New("go-syncmap").Parse(text)
	if err != nil {
		panic(err)
	}

	err = tmpl.Funcs(template.FuncMap{"trim": strings.TrimSpace}).Execute(&g.buf, &arg)
	if err != nil {
		panic(err)
	}
}

// generate produces the sync.Map method for the named type.
func (g *Generator) generate(typeInfo typeInfo) (render TmplMapRender) {
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

	nilValName, nilValDecl := g.createValAndNameDecl(values[0])

	return TmplMapRender{
		MapTypeName:     values[0].mapName,
		MapImport:       values[0].mapImport,
		KeyTypeName:     strings_.LoadElse(values[0].keyIsPointer, "*", "") + values[0].keyTypePrefix + values[0].keyType,
		KeyTypeImport:   values[0].keyImport,
		ValueTypeName:   strings_.LoadElse(values[0].valueIsPointer, "*", "") + values[0].valueTypePrefix + values[0].valueType,
		ValueTypeImport: values[0].valueImport,
		ValueTypeNilVal: nilValName,
		ValueTypeNilDecl: strings_.LoadElseGet(values[0].valueIsPointer, "nil", func() string {
			return nilValDecl
		}),
	}
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

// goimport returns the goimport-ed contents of the Generator's buffer.
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
