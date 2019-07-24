package _default

import (
	"github.com/searKing/golang/go/encoding/internal/tag"
	reflect2 "github.com/searKing/golang/go/reflect"
	"gopkg.in/yaml.v2"
	"reflect"
)

const TagDefault = "default"

// Convert wrapper of convertState
func Convert(val interface{}) error {
	return tag.Tag(val, func(val reflect.Value, tag reflect.StructTag) error {
		fn := newTypeConverter(func(val reflect.Value, tag reflect.StructTag) (isUserDefined bool, err error) {
			isUserDefined = false
			if !reflect2.IsEmptyValue(val) {
				return
			}
			defaultTag, ok := tag.Lookup(TagDefault)
			if !ok {
				return
			}
			return isUserDefined, yaml.Unmarshal([]byte(defaultTag), val.Addr().Interface())
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
