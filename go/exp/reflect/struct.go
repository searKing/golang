// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"reflect"

	reflect_ "github.com/searKing/golang/go/reflect"
)

// FieldByNames returns the struct field with the given names.
// It returns the zero Value if no field was found.
// It panics if v's Kind is not struct or x's Kind is not X.
func FieldByNames[V any, X any](v V, names ...string) (x X, ok bool) {
	val, ok := reflect_.FieldByNames(reflect.ValueOf(v), names...)
	if ok {
		return val.Interface().(X), ok
	}
	var zero X
	return zero, false
}

// SetFieldByNames assigns x to the value v.
// It panics if CanSet returns false.
// As in Go, x's value must be assignable to type of v's son, grandson, etc
func SetFieldByNames[V any, X any](v any, names []string, x X) (ok bool) {
	return reflect_.SetFieldByNames(reflect.ValueOf(v), names, reflect.ValueOf(x))
}
