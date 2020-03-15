// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import "reflect"

type TaggerError struct {
	Type reflect.Type
	Err  error
}

func (e *TaggerError) Error() string {
	return "default: error calling TagDefault for type " + e.Type.String() + ": " + e.Err.Error()
}

// An UnsupportedTypeError is returned by Marshal when attempting
// to handle an unsupported structTag type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "default: unsupported type: " + e.Type.String()
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return "default: unsupported structTag: " + e.Str
}
