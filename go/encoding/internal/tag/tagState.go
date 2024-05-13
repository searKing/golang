// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import (
	"reflect"
	"sync"
)

// An convertState encodes JSON into a bytes.Buffer.
type tagState struct{}

func (_ *tagState) Reset() {
	return
}

var tagStatePool sync.Pool

func newTagState() *tagState {
	if v := tagStatePool.Get(); v != nil {
		e := v.(*tagState)
		e.Reset()
		return e
	}
	return new(tagState)
}

// defaultError is an error wrapper type for internal use only.
// Panics with errors are wrapped in defaultError so that the top-level recover
// can distinguish intentional panics from this package.
type defaultError struct{ error }

func (e *tagState) handle(v any, opts tagOpts) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if je, ok := r.(defaultError); ok {
				err = je.error
			} else {
				panic(r)
			}
		}
	}()
	e.reflectValue(reflect.ValueOf(v), opts)
	return nil
}

// error aborts the encoding by panicking with err wrapped in defaultError.
func (e *tagState) error(err error) {
	panic(defaultError{err})
}
func (e *tagState) reflectValue(v reflect.Value, opts tagOpts) {
	valueTaggeFunc(v)(e, v, opts)
}
