// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package object

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

// DeepCopys returns a interface consisting of the deeply copying elements.
// just clone public&clonable elems - upper case - name & IsCloneable()==true
func DeepClone(obj interface{}) (copy interface{}) {
	if !IsNilable(obj) || !IsCloneable(obj) {
		return obj
	}
	v := reflect.ValueOf(obj)
	copyV := reflect.New(v.Type()).Elem()
	deepClones(obj, copyV.Addr().Interface())
	return copyV.Interface()
}

// deepClones provides the method to creates a deep copy of whatever is passed to
// it and returns the copy in an interface. The returned value will need to be
// asserted to the correct type.
func deepClones(origin, copy interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(origin); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(copy)
}

// IsCloneable Returns {@code true} if the arguments are a chan or func field
// or pointer to chan or func
// and {@code false} otherwise.
func IsCloneable(obj interface{}) bool {
	switch t := reflect.ValueOf(obj); t.Kind() {
	// All basic types are easy: they are predefined.
	case reflect.Bool:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	case reflect.Float32, reflect.Float64:
	case reflect.Complex64, reflect.Complex128:
	case reflect.String:
	case reflect.Interface:
	case reflect.Array:
	case reflect.Map:
	case reflect.Slice:
	case reflect.Struct:
	default:
		// If the field is a chan or func or pointer thereto, don't send it.
		// That is, treat it like an unexported field.
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Chan || t.Kind() == reflect.Func {
			return false
		}
		return true
	}
	return true
}
