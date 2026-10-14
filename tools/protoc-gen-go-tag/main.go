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
// Run it directly to export tag.proto, which declares the field option field_tag, to the include directory DIR:
//
//	protoc-gen-go-tag -export_proto DIR
//
// See the README and documentation for protocol buffers to learn more:
//
//	https://developers.google.com/protocol-buffers/
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/searKing/golang/tools/protoc-gen-go-tag/ast"
)

// tagProtoPath is the import path of tag.proto.
const tagProtoPath = "github.com/searKing/golang/tools/protoc-gen-go-tag/tag/tag.proto"

//go:embed tag/tag.proto
var tagProto []byte

func main() {
	// protoc runs plugins without arguments.
	if len(os.Args) > 1 {
		cmd := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
		exportProto := cmd.String("export_proto", "", "export tag.proto to `DIR`/"+tagProtoPath+", to include DIR by protoc -I")
		_ = cmd.Parse(os.Args[1:])
		if *exportProto == "" || cmd.NArg() > 0 {
			cmd.Usage()
			os.Exit(2)
		}
		if err := exportTagProto(*exportProto); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(generate)
}

// generate generates Go code of proto files as protoc-gen-go does, then rewrites struct tags by field_tag options.
// gengo, the generator of protoc-gen-go, is internal to the protobuf module and may change without notice.
func generate(gen *protogen.Plugin) error {
	gen.SupportedFeatures = gengo.SupportedFeatures
	gen.SupportedEditionsMinimum = gengo.SupportedEditionsMinimum
	gen.SupportedEditionsMaximum = gengo.SupportedEditionsMaximum
	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}
		if err := ast.Rewrite(gen, f, gengo.GenerateFile(gen, f)); err != nil {
			return err
		}
	}
	return nil
}

// exportTagProto writes tag.proto to dir, by its import path.
func exportTagProto(dir string) error {
	path := filepath.Join(dir, filepath.FromSlash(tagProtoPath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, tagProto, 0o644)
}
