// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TypeInfo for Parsing.
// Also includes a Lexical Analysis and Syntactic Analysis.

package main

import (
	"fmt"

	"github.com/searKing/golang/tools/pkg/ast/generic"
)

type tokenType int

const (
	tokenTypeParen tokenType = iota
	tokenTypeName
)

type _token struct {
	typ   tokenType
	value string
}

type typeInfo struct {
	// These fields are reset for each type being generated.
	Name            string // Name of the sync.Pool type.
	Import          string // import path of the sync.Pool type.
	valueType       string // The type of the value in sync.Pool.
	valueImport     string // import path of the sync.Pool's value.
	valueIsPointer  bool   // whether the value's type is ptr
	valueTypePrefix string // The type's prefix, such as []*[]
}

func newTypeInfo(input string) []typeInfo {
	var infos []typeInfo
	for _, info := range generic.New(input) {
		info_ := typeInfo{
			Name:   info.Name,
			Import: info.Import,
		}
		for i, template := range info.TemplateTypes {
			if i == 0 {
				info_.valueImport = template.Import
				info_.valueType = template.Type
				info_.valueIsPointer = template.IsPointer
				info_.valueTypePrefix = template.TypePrefix
				continue
			}
			panic(fmt.Sprintf("unexpected redundant #%d template type: %s, only 1 is expected", i, &template))
		}
		infos = append(infos, info_)
	}

	return infos
}
