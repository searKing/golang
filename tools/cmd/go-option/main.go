// Copyright 2023 The searKing Author. All rights reserved.
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
// // go:generate go-option -type "Pill"
// type Pill[T comparable] struct {
// 	//  This is Name doc comment
// 	Name      string //  This is Name line comment
// 	Age       string `option:",short"`
// 	Address   string `option:"-"`
// 	NameAlias string `option:"Title,"`
//
// 	arrayType     [5]T
// 	funcType      func()
// 	interfaceType interface{}
// 	mapType       map[string]int
// 	sliceType     []int64
// }
//
// running this command
//
//	go-option -type=Pill
//
// in the same directory will create the file pill_options.go, in package painkiller,
// containing a definition of
//

// //  A PillOption sets options.
// type PillOption[T comparable] interface {
// 	apply(*Pill[T])
// }
//
// //  EmptyPillOption does not alter the configuration. It can be embedded
// //  in another structure to build custom options.
// //
// //  This API is EXPERIMENTAL.
// type EmptyPillOption[T comparable] struct{}
//
// func (EmptyPillOption[T]) apply(*Pill[T]) {}
//
// //  PillOptionFunc wraps a function that modifies Pill[T] into an
// //  implementation of the PillOption[T comparable] interface.
// type PillOptionFunc[T comparable] func(*Pill[T])
//
// func (f PillOptionFunc[T]) apply(do *Pill[T]) {
// 	f(do)
// }
//
// //  ApplyOptions call apply() for all options one by one
// func (o *Pill[T]) ApplyOptions(options ...PillOption[T]) *Pill[T] {
// 	for _, opt := range options {
// 		if opt == nil {
// 			continue
// 		}
// 		opt.apply(o)
// 	}
// 	return o
// }
//
// //  WithPillName sets Name in Pill[T].
// //  This is Name doc comment
// //  This is Name line comment
// func WithPillName[T comparable](v string) PillOption[T] {
// 	return PillOptionFunc[T](func(o *Pill[T]) {
// 		o.Name = v
// 	})
// }
//
// //  WithAge sets Age in Pill[T].
// func WithAge[T comparable](v string) PillOption[T] {
// 	return PillOptionFunc[T](func(o *Pill[T]) {
// 		o.Age = v
// 	})
// }
//
// //  WithPillTitle sets NameAlias in Pill[T].
// func WithPillTitle[T comparable](v string) PillOption[T] {
// 	return PillOptionFunc[T](func(o *Pill[T]) {
// 		o.NameAlias = v
// 	})
// }
//
// //  WithPillArrayType sets arrayType in Pill[T].
// func WithPillArrayType[T comparable](v [5]T) PillOption[T] {
// 	return PillOptionFunc[T](func(o *Pill[T]) {
// 		o.arrayType = v
// 	})
// }
//
// //  WithPillInterfaceType sets interfaceType in Pill[T].
// func WithPillInterfaceType[T comparable](v interface{}) PillOption[T] {
// 	return PillOptionFunc[T](func(o *Pill[T]) {
// 		o.interfaceType = v
// 	})
// }
//
// //  WithPillMapType appends mapType in Pill[T].
// func WithPillMapType[T comparable](m map[string]int) PillOption[T] {
// 	return PillOptionFunc[T](func(o *Pill[T]) {
// 		if o.mapType == nil {
// 			o.mapType = m
// 			return
// 		}
// 		for k, v := range m {
// 			o.mapType[k] = v
// 		}
// 	})
// }
//
// //  WithPillMapTypeReplace sets mapType in Pill[T].
// func WithPillMapTypeReplace[T comparable](v map[string]int) PillOption[T] {
// 	return PillOptionFunc[T](func(o *Pill[T]) {
// 		o.mapType = v
// 	})
// }
//
// //  WithPillSliceType appends sliceType in Pill[T].
// func WithPillSliceType[T comparable](v ...int64) PillOption[T] {
// 	return PillOptionFunc[T](func(o *Pill[T]) {
// 		o.sliceType = append(o.sliceType, v...)
// 	})
// }
//
// //  WithPillSliceTypeReplace sets sliceType in Pill[T].
// func WithPillSliceTypeReplace[T comparable](v ...int64) PillOption[T] {
// 	return PillOptionFunc[T](func(o *Pill[T]) {
// 		o.sliceType = v
// 	})
// }

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
package main

import (
	"github.com/searKing/golang/tools/go-option/option"
)

func main() {
	option.Main()
}
