// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// go-union Generates Go code using a package as a union or sum type.
// Given the name of a Union type T,
// go-union will create a new self-contained Go source file implementing
//
//	func (m T) Union() any
//
// The file is created in the same package and directory as the package that defines T.
// It has helpful defaults designed for use with go generate.
//
// For example, given this snippet,
//
// running this command
//
//	go-union -type=Pill
//
// in the same directory will create the file pill_union.go, in package painkiller,
// containing a definition of
//
//	func (u Pill[T]) Union() any
//
// Typically this process would be run using go generate, like this:
//
//	//go:generate go-union -type=Pill
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must name a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -type flag accepts a comma-separated list of types so a single run can
// generate methods for multiple types. The default output file is t_string.go,
// where t is the lower-cased name of the first type listed. It can be overridden
// with the -output flag.
package main

import "github.com/searKing/golang/tools/go-union/union"

func main() {
	union.Main()
}
