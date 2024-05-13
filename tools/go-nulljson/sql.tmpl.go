// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "strings"

// SqlJsonType: NullJson type name
// ValueType: value type name
// NilValue: nil value of map type
type SqlRender struct {
	SqlJsonType string // NullJson type name
	valueImport string // import path of the atomic.Value's value.
	ValueType   string // value type name
	NilValue    string // nil value of map type

	ProtoJson bool // generate codec of proto by protojson, instead of json

	CanAlias bool // use type {{.SqlJsonType}} = {{.ValueType}} instead of type {{.SqlJsonType}} {{.ValueType}}
}

func (r *SqlRender) ResetCanAlias() {
	if strings.HasPrefix(strings.TrimSpace(r.ValueType), "map") ||
		strings.HasPrefix(strings.TrimSpace(r.ValueType), "[") {
		r.CanAlias = false
		return
	}
	if r.valueImport != "" {
		r.CanAlias = false
		return
	}

	r.CanAlias = true
	return
}
