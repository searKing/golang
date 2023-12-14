// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/searKing/golang/go/reflect"
	strings_ "github.com/searKing/golang/go/strings"

	pb "github.com/searKing/golang/tools/protoc-gen-go-tag/tag"
)

type FieldInfo struct {
	FieldNameInProto string
	FieldNameInGo    string
	FieldTag         reflect.StructTag
	UpdateStrategy   pb.FieldTag_UpdateStrategy
}
type StructInfo struct {
	StructNameInProto string
	StructNameInGo    string
	FieldInfos        []FieldInfo
}

type FileInfo struct {
	FileName    string
	StructInfos []StructInfo
}

func (si *StructInfo) FindField(name string) (FieldInfo, bool) {
	for _, f := range si.FieldInfos {
		if f.FieldNameInGo == name {
			return f, true
		}
	}
	return FieldInfo{}, false
}

// WalkDescriptorProto returns all struct infos of dp， which contains FieldTag.
func WalkDescriptorProto(g *protogen.Plugin, dp *descriptorpb.DescriptorProto, typeNames []string) []StructInfo {
	var ss []StructInfo

	s := StructInfo{}
	s.StructNameInProto = dp.GetName()
	s.StructNameInGo = CamelCaseSlice(append(typeNames, CamelCase(dp.GetName())))

	for _, field := range dp.GetField() {
		var oneofS *StructInfo
		if field.OneofIndex != nil { // Special Case: oneof
			oneofS = &StructInfo{}
			oneofS.StructNameInProto = field.GetName()
			oneofS.StructNameInGo = CamelCaseSlice(append(typeNames, CamelCase(dp.GetName()), CamelCase(field.GetName())))
		}

		f := HandleFieldDescriptorProto(g, field)
		if f != nil {
			if oneofS != nil {
				oneofS.FieldInfos = append(oneofS.FieldInfos, *f)
				if len(oneofS.FieldInfos) > 0 {
					ss = append(ss, *oneofS)
				}
			} else {
				s.FieldInfos = append(s.FieldInfos, *f)
			}
		}
	}

	typeNames = append(typeNames, CamelCase(dp.GetName()))

	for _, decl := range dp.GetOneofDecl() {
		declS := HandleOneofDescriptorProto(g, decl, typeNames)
		if declS != nil {
			if decl.GetOptions() != nil {
				v := proto.GetExtension(decl.GetOptions(), pb.E_FieldTag)
				switch v := v.(type) {
				case *pb.FieldTag:
					tag := v.GetStructTag()
					tags, err := reflect.ParseStructTag(tag)
					if err != nil {
						g.Error(fmt.Errorf("failed to parse struct tag in oneof field extension: %w", err))
						// ignore this tag
						break
					}
					s.FieldInfos = append(s.FieldInfos, FieldInfo{
						FieldNameInProto: decl.GetName(),
						FieldNameInGo:    CamelCase(decl.GetName()),
						FieldTag:         *tags,
						UpdateStrategy:   v.GetUpdateStrategy(),
					})
				}
			}

			ss = append(ss, *declS)
		}
	}

	if len(s.FieldInfos) > 0 {
		ss = append(ss, s)
	}

	for _, nest := range dp.GetNestedType() {
		nestSs := WalkDescriptorProto(g, nest, typeNames)
		if len(nestSs) > 0 {
			ss = append(ss, nestSs...)
		}
	}
	return ss
}

func HandleOneofDescriptorProto(g *protogen.Plugin, dp *descriptorpb.OneofDescriptorProto, typeNames []string) *StructInfo {
	if dp == nil {
		return nil
	}
	s := StructInfo{}
	s.StructNameInProto = dp.GetName()
	s.StructNameInGo = "is" + CamelCaseSlice(append(typeNames, CamelCase(dp.GetName())))
	return &s
}

func HandleFieldDescriptorProto(g *protogen.Plugin, field *descriptorpb.FieldDescriptorProto) *FieldInfo {
	if field.GetOptions() == nil {
		return nil
	}

	v := proto.GetExtension(field.GetOptions(), pb.E_FieldTag)
	switch v := v.(type) {
	case *pb.FieldTag:
		tag := v.GetStructTag()
		tags, err := reflect.ParseStructTag(tag)
		if err != nil {
			g.Error(fmt.Errorf("failed to parse struct tag in field extension: %w", err))
			// ignore this tag
			return nil
		}
		return &FieldInfo{
			FieldNameInProto: field.GetName(),
			FieldNameInGo:    CamelCase(field.GetName()),
			FieldTag:         *tags,
			UpdateStrategy:   v.GetUpdateStrategy(),
		}
	}
	return nil
}

func Rewrite(g *protogen.Plugin) {
	var protoFiles []FileInfo

	for _, protoFile := range g.Request.GetProtoFile() {
		if !strings_.SliceContains(g.Request.GetFileToGenerate(), protoFile.GetName()) {
			continue
		}
		f := FileInfo{}
		f.FileName = protoFile.GetName()

		for _, messageType := range protoFile.GetMessageType() {
			ss := WalkDescriptorProto(g, messageType, nil)
			if len(ss) > 0 {
				f.StructInfos = append(f.StructInfos, ss...)
			}
		}
		if len(f.StructInfos) > 0 {
			protoFiles = append(protoFiles, f)
		}
	}
	// FIXME: always generate *.pb.go, to replace protoc-go, avoid "Tried to write the same file twice"
	//if len(protoFiles) == 0 {
	//	return
	//}
	// g.Response() will generate files, so skip this step
	//if len(g.Response().GetFile()) == 0 {
	//	return
	//}

	rewriter := NewGenerator(protoFiles, g)
	for _, f := range g.Response().GetFile() {
		rewriter.ParseGoContent(f)
	}
	rewriter.Generate()
}
