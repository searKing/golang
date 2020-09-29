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

const PtrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const, sizeof *void

func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

// cmd/compile/internal/gc/dump.go
func IsZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.UnsafePointer:
		return v.Pointer() == 0
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return v.Len() == 0
	case reflect.Func:
		return v.IsNil()
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return true
		}
		break
	case reflect.Struct:
		break
	default:
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

func IsNilValue(v reflect.Value) (result bool) {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return
}

func FollowValuePointer(v reflect.Value) reflect.Value {
	v = reflect.Indirect(v)
	if v.Kind() == reflect.Ptr {
		return FollowValuePointer(v)
	}
	return v
}

// A field represents a single field found in a struct.
type FieldValueInfo struct {
	value       reflect.Value
	structField reflect.StructField
	index       []int
}

func (info FieldValueInfo) MiddleNodes() []interface{} {

	if !info.value.IsValid() {
		return nil
	}
	if IsNilType(info.value.Type()) {
		return nil
	}
	val := FollowValuePointer(info.value)
	if val.Kind() != reflect.Struct {
		return nil
	}

	var middles []interface{}
	// Scan typ for fields to include.
	for i := 0; i < val.NumField(); i++ {
		index := make([]int, len(info.index)+1)
		copy(index, info.index)
		index[len(info.index)] = i
		middles = append(middles, FieldValueInfo{
			value:       val.Field(i),
			structField: val.Type().Field(i),
			index:       index,
		})
	}
	return middles
}
func (info FieldValueInfo) Depth() int {
	return len(info.index)
}

func (info FieldValueInfo) Value() reflect.Value {
	return info.value
}

func (info FieldValueInfo) StructField() (reflect.StructField, bool) {
	if IsEmptyValue(reflect.ValueOf(info.structField)) {
		return info.structField, false
	}
	return info.structField, true
}

func (info FieldValueInfo) Index() []int {
	return info.index
}

func (info *FieldValueInfo) String() string {
	//if IsNilValue(info.value) {
	//	return fmt.Sprintf("%+v", nil)
	//}
	//info.value.String()
	//return fmt.Sprintf("%+v %+v", info.value.Type().String(), info.value)

	switch k := info.value.Kind(); k {
	case reflect.Invalid:
		return "<invalid value>"
	case reflect.String:
		return "[string: " + info.value.String() + "]"
	}
	// If you call String on a reflect.value of other type, it's better to
	// print something than to panic. Useful in debugging.
	return "[" + info.value.Type().String() + ":" + func() string {
		if info.value.CanInterface() && info.value.Interface() == nil {
			return "<nil value>"
		}
		return fmt.Sprintf(" %+v", info.value)
	}() + "]"
}

type FieldValueInfoHandler interface {
	Handler(info FieldValueInfo) (goon bool)
}
type FieldValueInfoHandlerFunc func(info FieldValueInfo) (goon bool)

func (f FieldValueInfoHandlerFunc) Handler(info FieldValueInfo) (goon bool) {
	return f(info)
}

func WalkValueDFS(val reflect.Value, handler FieldValueInfoHandler) {
	traversal.DepthFirstSearchOrder(FieldValueInfo{
		value: val,
	}, traversal.HandlerFunc(func(node interface{}, depth int) (goon bool) {
		return handler.Handler(node.(FieldValueInfo))
	}))
}

func WalkValueBFS(val reflect.Value, handler FieldValueInfoHandler) {
	traversal.BreadthFirstSearchOrder(FieldValueInfo{
		value: val,
	}, traversal.HandlerFunc(func(node interface{}, depth int) (goon bool) {
		return handler.Handler(node.(FieldValueInfo))
	}))
}

func DumpValueInfoDFS(v reflect.Value) string {
	dumpInfo := &bytes.Buffer{}
	first := true
	WalkValueDFS(v, FieldValueInfoHandlerFunc(func(info FieldValueInfo) (goon bool) {
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

func DumpValueInfoBFS(v reflect.Value) string {
	dumpInfo := &bytes.Buffer{}
	first := true
	WalkValueBFS(v, FieldValueInfoHandlerFunc(func(info FieldValueInfo) (goon bool) {
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
