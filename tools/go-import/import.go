// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go-import Performs auto import of non go files.
// Given the directory to be imported
// go-import will create gokeep.go Go source files and a new self-contained goimport.go Go source file.
// The gokeep.go file is created in the same package and directory as the cwd package.
// The goimport.go file is created in the package and directory under directories to be imported,
// It has helpful defaults designed for use with go generate.
//
// For example, given this snippet,
//
//	package painkiller
//
// running this command
//
//	go-import /dirs_to_be_force_imported
//
// in the same directory will create the file goimport.go,
// and in /dirs_to_be_force_imported will create the file gokeep.go
//
//
// Typically, this process would be run using go generate, like this:
//
//	//go:generate go-import
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must name a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -tag flag accepts a build tag string.
//
package main // import "github.com/searKing/golang/tools/go-import"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	importPrefix = flag.String("prefix", "", "root import path prefix of package, auto set if param is empty")
	globImport   = flag.String("glob", "*", "glob name filter of force import files")
	gokeepName   = flag.String("gokeep", "gokeep", "file name without a suffix to be imported by goimport")
	goimportName = flag.String("goimport", "goimport", "file name without a suffix to import gokeep")
	output       = flag.String("output", ".", "output file dir of goimport.go")
	buildTag     = flag.String("tag", "", "comma-separated list of build tags to apply")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of go-import:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-import [flags] [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgo-import [flags] -tag T [directory]\n")
	_, _ = fmt.Fprintf(os.Stderr, "For more information, see:\n")
	_, _ = fmt.Fprintf(os.Stderr, "\thttps://godoc.org/github.com/searKing/golang/tools/go-import\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

const (
	goImportToolName = "go-import"
)

type Import struct {
	GoImportToolName string
	GoImportToolArgs string
	ModuleName       string
	ImportPaths      []string
	BuildTags        []string
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", goImportToolName))
	flag.Usage = Usage
	flag.Parse()
	//if len(*importPrefix) == 0 {
	//	flag.Usage()
	//	os.Exit(2)
	//}

	// We accept either one directory, or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in the current directory.
		args = []string{"."}
	}

	// Parse the package once.
	g := NewGenerator(*importPrefix, *globImport, *buildTag)
	var dirs []string
	for _, arg := range args {
		if isDirectory(arg) {
			dirs = append(dirs, arg)
		} else {
			if len(*buildTag) != 0 {
				log.Fatal("-tag option applies only to directories, not when files are specified")
			}
			dirs = append(dirs, filepath.Dir(arg))
		}
	}
	if len(dirs) == 0 {
		dirs = append(dirs, ".")
		//log.Fatal("no dirs or files are specified in params")
	}

	g.scanPackage(dirs...)
	g.generate(goImportToolName, os.Args[1:]...)
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
