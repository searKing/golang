// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import (
	"reflect"

	reflect2 "github.com/searKing/golang/go/reflect"
)

type structTagFunc struct {
	fields       []field
	fieldTagFunc []tagFunc
}

func (se *structTagFunc) handle(state *tagState, v reflect.Value, opts tagOpts) (isUserDefined bool) {
	isUserDefined = false

	for i, f := range se.fields {
		fv := reflect2.ValueByStructFieldIndex(v, f.index)
		if !fv.IsValid() && reflect2.IsEmptyValue(fv) {
			continue
		}
		field := v.FieldByIndex(se.fields[i].index)

		//判断是否为可取指，可导出字段
		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		// continue if a userDefined func has been called
		if isFieldTagFuncUserDefined := se.fieldTagFunc[i](state, fv, opts); isFieldTagFuncUserDefined {
			continue
		}

		tagFn := opts.TagHandler
		if tagFn != nil {
			if err := tagFn(field, se.fields[i].structTag); err != nil {
				state.error(&TaggerError{v.Type(), err})
				return
			}
		}
	}
	return
}

func newStructTagFunc(t reflect.Type) tagFunc {
	fields := cachedTypeFields(t)
	se := &structTagFunc{
		fields:       fields,
		fieldTagFunc: make([]tagFunc, len(fields)),
	}
	for i, f := range fields {
		se.fieldTagFunc[i] = typeTagFunc(reflect2.TypeByStructFieldIndex(t, f.index))
	}
	return se.handle
}
