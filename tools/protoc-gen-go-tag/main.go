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
// The parameter field_tag_number sets extension numbers of FieldOptions read as field_tag, in addition to
// 1391 of searking.tag.v1.field_tag, 65000 of the deprecated google.protobuf.field_tag by default:
//
//	protoc --go-tag_out=paths=source_relative,field_tag_number=60001:. path/to/astFile.proto
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
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/encoding/protowire"

	"github.com/searKing/golang/tools/protoc-gen-go-tag/ast"
)

// tagProtos are tag.proto files by their import paths.
var tagProtos = map[string]string{
	"searking/tag/v1/tag.proto":                                        "proto/searking/tag/v1/tag.proto",
	"github.com/searKing/golang/tools/protoc-gen-go-tag/tag/tag.proto": "tag/tag.proto", // deprecated
}

//go:embed proto/searking/tag/v1/tag.proto tag/tag.proto
var tagProtoFS embed.FS

func main() {
	// protoc runs plugins without arguments.
	if len(os.Args) > 1 {
		cmd := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
		exportProto := cmd.String("export_proto", "", "export tag.proto files to `DIR` by their import paths, to include DIR by protoc -I")
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

	opts, run := plugin()
	opts.Run(run)
}

// plugin returns options of protogen, parsing parameters of the plugin, and the function to run the plugin.
func plugin() (protogen.Options, func(gen *protogen.Plugin) error) {
	var flags flag.FlagSet
	var numbers fieldTagNumbers
	flags.Var(&numbers, "field_tag_number", "extension number of FieldOptions read as field_tag, may be repeated")
	return protogen.Options{ParamFunc: flags.Set}, func(gen *protogen.Plugin) error {
		if numbers == nil {
			numbers = fieldTagNumbers{ast.LegacyFieldTagNumber}
		}
		return generate(gen, numbers...)
	}
}

// fieldTagNumbers is a flag.Value of extension numbers, set repeatedly.
type fieldTagNumbers []protowire.Number

func (n *fieldTagNumbers) String() string {
	var s []string
	for _, v := range *n {
		s = append(s, strconv.Itoa(int(v)))
	}
	return strings.Join(s, ",")
}

func (n *fieldTagNumbers) Set(s string) error {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil || !protowire.Number(v).IsValid() {
		return fmt.Errorf("invalid extension number %q", s)
	}
	*n = append(*n, protowire.Number(v))
	return nil
}

// generate generates Go code of proto files as protoc-gen-go does, then rewrites struct tags by field_tag options,
// of searking.tag.v1.field_tag and fieldTagNumbers, and @gotags comments.
// gengo, the generator of protoc-gen-go, is internal to the protobuf module and may change without notice.
func generate(gen *protogen.Plugin, fieldTagNumbers ...protowire.Number) error {
	gen.SupportedFeatures = gengo.SupportedFeatures
	gen.SupportedEditionsMinimum = gengo.SupportedEditionsMinimum
	gen.SupportedEditionsMaximum = gengo.SupportedEditionsMaximum
	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}
		if err := ast.Rewrite(gen, f, gengo.GenerateFile(gen, f), fieldTagNumbers...); err != nil {
			return err
		}
	}
	return nil
}

// exportTagProto writes tag.proto files to dir, by their import paths.
func exportTagProto(dir string) error {
	for importPath, name := range tagProtos {
		b, err := tagProtoFS.ReadFile(name)
		if err != nil {
			return err
		}
		path := filepath.Join(dir, filepath.FromSlash(importPath))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, b, 0o644); err != nil {
			return err
		}
	}
	return nil
}
