// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import "reflect"

// FieldByNames returns the struct field with the given names.
// It returns the zero Value if no field was found.
// It panics if v's Kind is not struct.
func FieldByNames(v reflect.Value, names ...string) (x reflect.Value, ok bool) {
	if !v.IsValid() || v.IsNil() {
		return reflect.ValueOf(nil), false
	}

	if len(names) == 0 {
		return reflect.ValueOf(nil), false
	}

	f := reflect.Indirect(v).FieldByName(names[0])
	if len(names) == 1 {
		if f.IsValid() {
			return f, true
		}
		return reflect.ValueOf(nil), false
	}
	return FieldByNames(f, names[1:]...)
}

// SetFieldByNames assigns x to the value v.
// It panics if CanSet returns false.
// As in Go, x's value must be assignable to type of v's son, grandson, etc
func SetFieldByNames(v reflect.Value, names []string, x reflect.Value) (ok bool) {
	if !v.IsValid() || v.IsNil() {
		return false
	}
	if len(names) == 0 {
		return false
	}

	f := reflect.Indirect(v).FieldByName(names[0])
	if len(names) == 1 {
		if f.IsValid() && f.Kind() == x.Kind() {
			f.Set(x)
			return true
		}
		return false
	}
	return SetFieldByNames(f, names[1:], x)
}
