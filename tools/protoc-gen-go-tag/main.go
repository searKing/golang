// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// protoc-gen-go-tag is a plugin for the Google protocol buffer compiler to Generate
// Go code. Install it by building this program and making it accessible within
// your PATH with the name:
//
//	protoc-gen-go-tag
//
// The 'go' suffix becomes part of the argument for the protocol compiler,
// such that it can be invoked as:
//
//	protoc --go-tag_out=paths=source_relative:. path/to/astFile.proto
//
// This generates Go bindings for the protocol buffer defined by astFile.proto.
// With that input, the output will be written to:
//
//	path/to/astFile.pb.go
//
// See the README and documentation for protocol buffers to learn more:
//
//	https://developers.google.com/protocol-buffers/
package main

import (
	"flag"

	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/searKing/golang/tools/protoc-gen-go-tag/ast"
)

func main() {
	var (
		flags flag.FlagSet
	)
	// For Debug Only
	//{ // Dump
	//	in, err := io.ReadAll(os.Stdin)
	//	if err != nil {
	//		panic(err)
	//	}
	//	if err := os.WriteFile("in.pb", in, 0666); err != nil {
	//		panic(err)
	//	}
	//}
	//{ // Debug
	//	os.Stdin, _ = os.Open("in.pb")
	//	os.Stdout, _ = os.Create("out.pb")
	//}

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = gengo.SupportedFeatures
		var originFiles []*protogen.GeneratedFile
		for _, f := range gen.Files {
			if f.Generate {
				originFiles = append(originFiles, gengo.GenerateFile(gen, f))
			}
		}
		ast.Rewrite(gen)

		for _, f := range originFiles {
			f.Skip()
		}
		return nil
	})
}
