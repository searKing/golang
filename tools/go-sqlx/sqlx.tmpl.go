// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//	StructType: NullJson type trimmedStructName
//	TableName: value type trimmedStructName
//	NilValue: nil value of map type
type SqlxRender struct {
	StructType string        // NullJson type trimmedStructName
	TableName  string        // value type trimmedStructName
	Fields     []StructField // struct Fields

	WithDao       bool   // generate with dao
	WithQueryInfo bool   // generate error with sql query executed
	NilValue      string // nil value of map type
}
