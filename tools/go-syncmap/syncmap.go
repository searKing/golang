// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go-syncmap Generates Go code using a package as a generic template for sync.Map.
// Given the name of a sync.Map type T , and the name of a type Key and Value
// go-syncmap will create a new self-contained Go source file implementing
// from Go version 1.9 onward
//	func (m *T) Store(key Key, value Value)
//	func (m *T) LoadOrStore(key Key, value Value) (Value, bool)
//	func (m *T) Load(key Key) (Value, bool)
//	func (m *T) Delete(key Key)
//	func (m *T) Range(f func(key Key, value Value) bool
// from Go version 1.15 onward
//	func (m *T) LoadAndDelete(key Key) (Value, bool)
//
// The file is created in the same package and directory as the package that defines T, Key and Value.
// It has helpful defaults designed for use with go generate.
//
// For example, given this snippet,
//
//	package painkiller
//
//	import "sync"
//
//	type Pill sync.Map
//
//
// running this command
//
//	go-syncmap -type Pill<int,time.Time>
//
// in the same directory will create the file pill_syncmap.go, in package painkiller,
// containing a definition of
//
// from Go version 1.9 onward
//	func (m *Pill) Store(key int, value time.Time)
//	func (m *Pill) LoadOrStore(key int, value time.Time) (time.Time, bool)
//	func (m *Pill) Load(key int) (time.Time, bool)
//	func (m *Pill) Delete(key int)
//	func (m *Pill) Range(f func(key int, value time.Time) bool
//
// from Go version 1.15 onward
//	func (m *Pill) LoadAndDelete(key int) (string, bool)
//
// Typically this process would be run using go generate, like this:
//
//	//go:generate go-syncmap -type Pill<int, string>
//	//go:generate go-syncmap -type Pill<int, time.Time>
//	//go:generate go-syncmap -type Pill<int, encoding/json.Token>
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
package main // import "github.com/searKing/golang/tools/go-syncmap"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	typeInfos   = flag.String("type", "", "comma-separated list of type names; must be set")
	output      = flag.String("output", "", "output file name; default srcdir/<type>_syncmap.go")
	trimprefix  = flag.String("trimprefix", "", "trim the `prefix` from the generated constant names")
	linecomment = flag.Bool("linecomment", false, "use line comment text as printed text when present")
	buildTags   = flag.String("tags", "", "comma-separated list of build tags to apply")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of go-syncmap:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-syncmap [flags] -type T [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-syncmap [flags] -type T<K,V> [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-syncmap [flags] -type T<K,V> files... # Must be a single package\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-syncmap [flags] -type T,S [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-syncmap [flags] -type T<K,V>,S<K,V> [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "For more information, see:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\thttps://godoc.org/github.com/searKing/golang/tools/go-syncmap\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

const (
	goSyncMapToolName = "go-syncmap"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("go-syncmap: ")
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

	var render TmplPackageRender
	render.GenerateToolArgs = goSyncMapToolName
	render.GenerateToolArgs = strings.Join(os.Args[1:], " ")
	render.WithSyncMapMethod = WithSyncMapMethod
	render.WithMethodLoadAndDelete = WithMethodLoadAndDelete
	render.PackageName = g.pkg.name

	// Run generate for each type.
	for _, typeInfo := range types {
		render.MapRenders = append(render.MapRenders, g.generate(typeInfo))
	}

	// buildOneRun generates the variables and String method for a single run of contiguous values.
	//The generated code is simple enough to write as a Printf format.
	g.Render(tmplPackageCode, render)

	// Format the output.
	src := g.format()

	target := g.goimport(src)

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("%s_syncmap.go", types[0].mapName)
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err := ioutil.WriteFile(outputName, target, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
