// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"reflect"
)

// A field represents a single field found in a struct.
type field struct {
	sf reflect.StructField

	name      string
	nameBytes []byte                 // []byte(name)
	equalFold func(s, t []byte) bool // bytes.EqualFold or equivalent

	tag       bool
	index     []int
	typ       reflect.Type
	omitEmpty bool
	quoted    bool
}

// v[i][j][...], v[i] is v's ith field, v[i][j] is v[i]'s jth field
func ValueByStructFieldIndex(v reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}

// t[i][j][...], t[i] is t's ith field, t[i][j] is t[i]'s jth field
func TypeByStructFieldIndex(t reflect.Type, index []int) reflect.Type {
	for _, i := range index {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		t = t.Field(i).Type
	}
	return t
}

func IsFieldExported(sf reflect.StructField) bool {
	return sf.PkgPath == ""
}
