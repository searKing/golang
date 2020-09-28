// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "strings"

//	SqlJsonType: NullJson type name
//	ValueType: value type name
//	NilValue: nil value of map type
type SqlRender struct {
	SqlJsonType string // NullJson type name
	ValueType   string // value type name
	NilValue    string // nil value of map type

	CanAlias bool // use type {{.SqlJsonType}} = {{.ValueType}} instead of type {{.SqlJsonType}} {{.ValueType}}
}

func (r *SqlRender) ResetCanAlias() {
	if r.SqlJsonType != r.ValueType {
		r.CanAlias = false
		return
	}

	if strings.HasPrefix(strings.TrimSpace(r.SqlJsonType), "map") ||
		strings.HasPrefix(strings.TrimSpace(r.SqlJsonType), "[") {
		r.CanAlias = false
		return
	}
	r.CanAlias = true
	return
}
