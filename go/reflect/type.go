package reflect

import (
	"bytes"
	"fmt"
	bytes_ "github.com/searKing/golang/go/bytes"
	"github.com/searKing/golang/go/container/traversal"
	"reflect"
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

func (thiz FieldTypeInfo) Middles() []interface{} {
	typ := thiz.structField.Type
	middles := []interface{}{}
	typ = FollowTypePointer(typ)
	if IsNilType(typ) {
		return nil
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	// Scan typ for fields to include.
	for i := 0; i < typ.NumField(); i++ {
		index := make([]int, len(thiz.index)+1)
		copy(index, thiz.index)
		index[len(thiz.index)] = i
		sf := typ.Field(i)
		middles = append(middles, FieldTypeInfo{
			structField: sf,
			index:       index,
		})
	}
	return middles
}

func (thiz FieldTypeInfo) Depth() int {
	return len(thiz.index)
}

func (thiz FieldTypeInfo) StructField() (reflect.StructField, bool) {
	if IsEmptyValue(reflect.ValueOf(thiz.structField)) {
		return thiz.structField, false
	}
	return thiz.structField, true
}

func (thiz FieldTypeInfo) Index() []int {
	return thiz.index
}

func (thiz FieldTypeInfo) String() string {
	if thiz.structField.Type == nil {
		return fmt.Sprintf("%+v", nil)
	}
	return fmt.Sprintf("%+v", thiz.structField.Type.String())
}

// Breadth First Search
func WalkTypeBFS(typ reflect.Type, parseFn func(info FieldTypeInfo) (goon bool)) {
	traversal.BFS(FieldTypeInfo{
		structField: reflect.StructField{
			Type: typ,
		},
	}, nil, func(ele interface{}, depth int) (gotoNextLayer bool) {
		return parseFn(ele.(FieldTypeInfo))
	})
}

// Wid First Search
func WalkTypeDFS(typ reflect.Type, parseFn func(info FieldTypeInfo) (goon bool)) {
	traversal.DFS(FieldTypeInfo{
		structField: reflect.StructField{
			Type: typ,
		},
	}, nil, func(ele interface{}, depth int) (gotoNextLayer bool) {
		return parseFn(ele.(FieldTypeInfo))
	})
}
func DumpTypeInfoDFS(t reflect.Type) string {
	dumpInfo := &bytes.Buffer{}
	first := true
	WalkTypeDFS(t, func(info FieldTypeInfo) (goon bool) {
		if first {
			first = false
			bytes_.NewIndent(dumpInfo, "", "\t", info.Depth())
		} else {
			bytes_.NewLine(dumpInfo, "", "\t", info.Depth())
		}
		dumpInfo.WriteString(fmt.Sprintf("%+v", info.String()))
		return true
	})
	return dumpInfo.String()
}
func DumpTypeInfoBFS(t reflect.Type) string {
	dumpInfo := &bytes.Buffer{}
	first := true
	WalkTypeBFS(t, func(info FieldTypeInfo) (goon bool) {
		if first {
			first = false
			bytes_.NewIndent(dumpInfo, "", "\t", info.Depth())
		} else {
			bytes_.NewLine(dumpInfo, "", "\t", info.Depth())
		}
		dumpInfo.WriteString(fmt.Sprintf("%+v", info.String()))
		return true
	})
	return dumpInfo.String()
}
