// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"bytes"
	"fmt"
	"reflect"

	bytes_ "github.com/searKing/golang/go/bytes"
	"github.com/searKing/golang/go/container/traversal"
)

// nil, unknown type
func IsNilType(v reflect.Type) (result bool) {
	return v == nil
}
func FollowTypePointer(v reflect.Type) reflect.Type {
	if IsNilType(v) {
		return v
	}
	if v.Kind() == reflect.Ptr {
		return FollowTypePointer(v.Elem())
	}
	return v
}

// A field represents a single field found in a struct.
type FieldTypeInfo struct {
	structField reflect.StructField
	index       []int
}

func (info FieldTypeInfo) MiddleNodes() []interface{} {
	typ := info.structField.Type
	var middles []interface{}
	typ = FollowTypePointer(typ)
	if IsNilType(typ) {
		return nil
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	// Scan typ for fields to include.
	for i := 0; i < typ.NumField(); i++ {
		index := make([]int, len(info.index)+1)
		copy(index, info.index)
		index[len(info.index)] = i
		sf := typ.Field(i)
		middles = append(middles, FieldTypeInfo{
			structField: sf,
			index:       index,
		})
	}
	return middles
}

func (info FieldTypeInfo) Depth() int {
	return len(info.index)
}

func (info FieldTypeInfo) StructField() (reflect.StructField, bool) {
	if IsEmptyValue(reflect.ValueOf(info.structField)) {
		return info.structField, false
	}
	return info.structField, true
}

func (info FieldTypeInfo) Index() []int {
	return info.index
}

func (info FieldTypeInfo) String() string {
	if info.structField.Type == nil {
		return fmt.Sprintf("%+v", nil)
	}
	return fmt.Sprintf("%+v", info.structField.Type.String())
}

type FieldTypeInfoHandler interface {
	Handler(info FieldTypeInfo) (goon bool)
}
type FieldTypeInfoHandlerFunc func(info FieldTypeInfo) (goon bool)

func (f FieldTypeInfoHandlerFunc) Handler(info FieldTypeInfo) (goon bool) {
	return f(info)
}

// Breadth First Search
func WalkTypeBFS(typ reflect.Type, handler FieldTypeInfoHandler) {
	traversal.BreadthFirstSearchOrder(FieldTypeInfo{
		structField: reflect.StructField{
			Type: typ,
		},
	}, traversal.HandlerFunc(func(node interface{}, depth int) (goon bool) {
		return handler.Handler(node.(FieldTypeInfo))
	}))
}

// Wid First Search
func WalkTypeDFS(typ reflect.Type, handler FieldTypeInfoHandler) {
	traversal.DepthFirstSearchOrder(FieldTypeInfo{
		structField: reflect.StructField{
			Type: typ,
		},
	}, traversal.HandlerFunc(func(node interface{}, depth int) (goon bool) {
		return handler.Handler(node.(FieldTypeInfo))
	}))
}

func DumpTypeInfoDFS(t reflect.Type) string {
	dumpInfo := &bytes.Buffer{}
	first := true
	WalkTypeDFS(t, FieldTypeInfoHandlerFunc(func(info FieldTypeInfo) (goon bool) {
		if first {
			first = false
			bytes_.NewIndent(dumpInfo, "", "\t", info.Depth())
		} else {
			bytes_.NewLine(dumpInfo, "", "\t", info.Depth())
		}
		dumpInfo.WriteString(fmt.Sprintf("%+v", info.String()))
		return true
	}))
	return dumpInfo.String()
}

func DumpTypeInfoBFS(t reflect.Type) string {
	dumpInfo := &bytes.Buffer{}
	first := true
	WalkTypeBFS(t, FieldTypeInfoHandlerFunc(func(info FieldTypeInfo) (goon bool) {
		if first {
			first = false
			bytes_.NewIndent(dumpInfo, "", "\t", info.Depth())
		} else {
			bytes_.NewLine(dumpInfo, "", "\t", info.Depth())
		}
		dumpInfo.WriteString(fmt.Sprintf("%+v", info.String()))
		return true
	}))
	return dumpInfo.String()
}
