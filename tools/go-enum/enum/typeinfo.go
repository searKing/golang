// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TypeInfo for Parsing.
// Also includes a Lexical Analysis and Syntactic Analysis.

package enum

import (
	"fmt"

	"github.com/searKing/golang/tools/pkg/ast/generic"
)

type typeInfo struct {
	// These fields are reset for each type being generated.
	Name   string // Name of the atomic.Value type.
	Import string // import path of the atomic.Value type.
}

func newTypeInfo(input string) []typeInfo {
	var infos []typeInfo
	for _, info := range generic.New(input) {
		info_ := typeInfo{
			Name:   info.Name,
			Import: info.Import,
		}
		if len(info.TemplateTypes) > 0 {
			panic(fmt.Sprintf("unexpected redundant #%d template type: %+v, only 0 is expected", 0, info.TemplateTypes))
		}
		infos = append(infos, info_)
	}

	return infos
}
