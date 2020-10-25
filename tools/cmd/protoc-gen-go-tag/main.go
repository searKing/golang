// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// protoc-gen-go-tag is a plugin for the Google protocol buffer compiler to Generate
// Go code. Install it by building this program and making it accessible within
// your PATH with the name:
//	protoc-gen-go-tag
//
// The 'go' suffix becomes part of the argument for the protocol compiler,
// such that it can be invoked as:
//	protoc --go-tag_out=paths=source_relative:. path/to/astFile.proto
//
// This generates Go bindings for the protocol buffer defined by astFile.proto.
// With that input, the output will be written to:
//	path/to/astFile.pb.go
//
// See the README and documentation for protocol buffers to learn more:
//	https://developers.google.com/protocol-buffers/
package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/generator"
	"github.com/searKing/golang/tools/cmd/protoc-gen-go-tag/gen"
)

func main() {
	// Begin by allocating a generator. The request and response structures are stored there
	// so we can do error handling easily - the response structure contains the field to
	// report failure.
	g := generator.New()

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}
	//ioutil.WriteFile("in.pb", data, 0666)
	//return
	//data, err := ioutil.ReadFile("in.pb")
	//if err != nil {
	//	g.Error(err, "reading input")
	//}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no goFiles to Generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the astFile that defines them.
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()
	gen.Rewrite(g)

	// Send back the results.
	data, err = proto.Marshal(g.Response)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
