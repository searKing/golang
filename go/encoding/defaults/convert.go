// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package defaults

import (
	"reflect"

	"github.com/searKing/golang/go/encoding/internal/tag"
	reflect_ "github.com/searKing/golang/go/reflect"
)

const TagDefault = "default"

// Convert wrapper of convertState
func Convert(val interface{}, unmarshal func(data []byte, v interface{}) error) error {
	return tag.Tag(val, func(val reflect.Value, tag reflect.StructTag) error {
		fn := newTypeConverter(func(val reflect.Value, tag reflect.StructTag) (isUserDefined bool, err error) {
			isUserDefined = false
			if !reflect_.IsEmptyValue(val) {
				return
			}
			defaultTag, ok := tag.Lookup(TagDefault)
			if !ok {
				return
			}
			return isUserDefined, unmarshal([]byte(defaultTag), val.Addr().Interface())
		}, val.Type(), true)

		_, err := fn(val, tag)
		return err
	})
}

// Marshaler is the interface implemented by types that
// can marshal themselves into valid JSON.
type Converter interface {
	ConvertDefault(val reflect.Value, tag reflect.StructTag) error
}

var converterType = reflect.TypeOf(new(Converter)).Elem()
