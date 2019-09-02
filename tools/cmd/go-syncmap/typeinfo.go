// Copyright 2019 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TypeInfo for Parsing.
// Also includes a Lexical Analysis and Syntactic Analysis.

package main

import (
	"fmt"
	generic2 "github.com/searKing/golang/tools/common/ast/generic"
)

type typeInfo struct {
	// These fields are reset for each type being generated.
	mapName        string // Name of the sync.Map type.
	mapImport      string // import path of the sync.Map type.
	keyType        string // The type of the key in sync.Map.
	keyImport      string // import path of the sync.Map's key.
	keyIsPointer   bool   // whether the value's type is ptr
	valueType      string // The type of the value in sync.Map.
	valueImport    string // import path of the sync.Map's value.
	valueIsPointer bool   // whether the value's type is ptr
}

func newTypeInfo(input string) []typeInfo {
	var infos []typeInfo
	for _, info := range generic2.New(input) {
		info_ := typeInfo{
			mapName:   info.Name,
			mapImport: info.Import,
		}
		for i, template := range info.TemplateTypes {
			if i == 0 {
				info_.keyImport = template.Import
				info_.keyType = template.Type
				info_.keyIsPointer = template.IsPointer
				continue
			}
			if i == 1 {
				info_.valueImport = template.Import
				info_.valueType = template.Type
				info_.valueIsPointer = template.IsPointer
				continue
			}
			panic(fmt.Sprintf("unexpected redundant #%d template type: %s, only 2 is expected", i, &template))
		}
		infos = append(infos, info_)
	}

	return infos
}
