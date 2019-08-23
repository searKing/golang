package reflect

import (
	"bytes"
	"fmt"
	bytes_ "github.com/searKing/golang/go/bytes"
	"github.com/searKing/golang/go/container/traversal"
	"reflect"
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

func IsZeroValue(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if v.IsNil() {
			return true
		}
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

func (thiz FieldValueInfo) Middles() []interface{} {

	if !thiz.value.IsValid() {
		return nil
	}
	if IsNilType(thiz.value.Type()) {
		return nil
	}
	val := FollowValuePointer(thiz.value)
	if val.Kind() != reflect.Struct {
		return nil
	}

	middles := []interface{}{}
	// Scan typ for fields to include.
	for i := 0; i < val.NumField(); i++ {
		index := make([]int, len(thiz.index)+1)
		copy(index, thiz.index)
		index[len(thiz.index)] = i
		middles = append(middles, FieldValueInfo{
			value:       val.Field(i),
			structField: val.Type().Field(i),
			index:       index,
		})
	}
	return middles
}
func (thiz FieldValueInfo) Depth() int {
	return len(thiz.index)
}

func (thiz FieldValueInfo) Value() reflect.Value {
	return thiz.value
}

func (thiz FieldValueInfo) StructField() (reflect.StructField, bool) {
	if IsEmptyValue(reflect.ValueOf(thiz.structField)) {
		return thiz.structField, false
	}
	return thiz.structField, true
}

func (thiz FieldValueInfo) Index() []int {
	return thiz.index
}

func (thiz *FieldValueInfo) String() string {
	//if IsNilValue(thiz.value) {
	//	return fmt.Sprintf("%+v", nil)
	//}
	//thiz.value.String()
	//return fmt.Sprintf("%+v %+v", thiz.value.Type().String(), thiz.value)

	switch k := thiz.value.Kind(); k {
	case reflect.Invalid:
		return "<invalid value>"
	case reflect.String:
		return "[string: " + thiz.value.String() + "]"
	}
	// If you call String on a reflect.value of other type, it's better to
	// print something than to panic. Useful in debugging.
	return "[" + thiz.value.Type().String() + ":" + func() string {
		if thiz.value.CanInterface() && thiz.value.Interface() == nil {
			return "<nil value>"
		}
		return fmt.Sprintf(" %+v", thiz.value)
	}() + "]"
}
func WalkValueDFS(val reflect.Value, parseFn func(info FieldValueInfo) (goon bool)) {
	traversal.BreadthFirstSearchOrder(FieldValueInfo{
		value: val,
	}, nil, func(ele interface{}, depth int) (gotoNextLayer bool) {
		return parseFn(ele.(FieldValueInfo))
	})
}

// Breadth First Search
func WalkValueBFS(val reflect.Value, parseFn func(info FieldValueInfo) (goon bool)) {
	traversal.BreadthFirstSearchOrder(FieldValueInfo{value: val},
		nil, func(ele interface{}, depth int) (gotoNextLayer bool) {
			return parseFn(ele.(FieldValueInfo))
		})
}

func DumpValueInfoDFS(v reflect.Value) string {
	dumpInfo := &bytes.Buffer{}
	first := true
	WalkValueDFS(v, func(info FieldValueInfo) (goon bool) {
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

func DumpValueInfoBFS(v reflect.Value) string {
	dumpInfo := &bytes.Buffer{}
	first := true
	WalkValueBFS(v, func(info FieldValueInfo) (goon bool) {
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
