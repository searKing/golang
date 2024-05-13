// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"bytes"
	"reflect"
	"strings"
	_ "unsafe" // for go:linkname

	"github.com/mitchellh/mapstructure"
	maps_ "github.com/searKing/golang/go/exp/maps"
	json_ "github.com/searKing/golang/third_party/github.com/spf13/viper/json"
	"github.com/searKing/golang/third_party/google.golang.org/protobuf/encoding/protojson"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
)

//go:linkname decode github.com/spf13/viper.decode
func decode(input any, config *mapstructure.DecoderConfig) error

//go:linkname defaultDecoderConfig github.com/spf13/viper.defaultDecoderConfig
func defaultDecoderConfig(output any, opts ...viper.DecoderConfigOption) *mapstructure.DecoderConfig

// DecodeProtoJsonHook if set, will be called before any decoding and any
// type conversion (if WeaklyTypedInput is on). This lets you modify
// the values before they're set down onto the resulting struct.
//
// If an error is returned, the entire decode will fail with that
// error.
func DecodeProtoJsonHook(v proto.Message) viper.DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.TagName = "json" // trick of protobuf, which generates json tag only
		c.WeaklyTypedInput = true
		c.ZeroFields = false
		c.Result = v
		if c.ZeroFields {
			c.DecodeHook = UnmarshalProtoMessageHookFunc(nil)
		} else {
			// v as default
			c.DecodeHook = UnmarshalProtoMessageHookFunc(v)
		}
	}
}

// UnmarshalKey takes a single key and unmarshalls it into a Struct.
// use protojson to decode if rawVal is proto.Message
func UnmarshalKey(key string, rawVal any, opts ...viper.DecoderConfigOption) error {
	return UnmarshalKeyViper(viper.GetViper(), key, rawVal, opts...)
}

func UnmarshalKeyViper(v *viper.Viper, key string, rawVal any, opts ...viper.DecoderConfigOption) error {
	return UnmarshalKeysViper(v, strings.Split(key, "."), rawVal, opts...)
}

func UnmarshalKeys(keys []string, rawVal any, opts ...viper.DecoderConfigOption) error {
	return UnmarshalKeysViper(viper.GetViper(), keys, rawVal)
}

func UnmarshalKeysViper(v *viper.Viper, keys []string, rawVal any, opts ...viper.DecoderConfigOption) error {
	if val, ok := rawVal.(proto.Message); ok {
		opts = append([]viper.DecoderConfigOption{DecodeProtoJsonHook(val)}, opts...)
	}

	if v == nil {
		v = viper.GetViper()
	}
	if len(keys) == 0 {
		return v.Unmarshal(rawVal, opts...)
	}
	c := maps_.NestedMap[string](v.AllSettings())
	val, has := c.Load(keys)
	if !has {
		return nil
	}
	return decode(val, defaultDecoderConfig(rawVal, opts...))
}

// Unmarshal unmarshalls the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
// use protojson to decode if rawVal is proto.Message
func Unmarshal(rawVal any, opts ...viper.DecoderConfigOption) error {
	return UnmarshalKeys(nil, rawVal, opts...)
}

func UnmarshalViper(v *viper.Viper, rawVal any, opts ...viper.DecoderConfigOption) error {
	return UnmarshalKeysViper(v, nil, rawVal, opts...)
}

// UnmarshalExact unmarshals the config into a Struct, erroring if a field is nonexistent
// in the destination struct.
// use protojson to decode if rawVal is proto.Message
func UnmarshalExact(rawVal any, opts ...viper.DecoderConfigOption) error {
	return UnmarshalExactViper(viper.GetViper(), rawVal, opts...)
}

func UnmarshalKeysExactViper(v *viper.Viper, keys []string, rawVal any, opts ...viper.DecoderConfigOption) error {
	if val, ok := rawVal.(proto.Message); ok {
		opts = append([]viper.DecoderConfigOption{DecodeProtoJsonHook(val)}, opts...)
	}

	if v == nil {
		v = viper.GetViper()
	}
	if len(keys) == 0 {
		return v.UnmarshalExact(rawVal, opts...)
	}
	c := maps_.NestedMap[string](v.AllSettings())
	val, has := c.Load(keys)
	if !has {
		return nil
	}
	config := defaultDecoderConfig(rawVal, opts...)
	config.ErrorUnused = true
	return decode(val, config)
}

func UnmarshalExactViper(v *viper.Viper, rawVal any, opts ...viper.DecoderConfigOption) error {
	return UnmarshalKeysExactViper(v, nil, rawVal, opts...)
}

// UnmarshalProtoMessageHookFunc returns a DecodeHookFunc that converts
// root struct to config.ViperProto.
// Trick of protobuf, which generates json tag only
// def is the default value of dst
func UnmarshalProtoMessageHookFunc(def proto.Message) mapstructure.DecodeHookFunc {
	return func(src reflect.Type, dst reflect.Type, data any) (any, error) {
		dataProto, ok := reflect.New(dst).Interface().(proto.Message)
		if !ok {
			return data, nil
		}

		// trick(json): error decoding '': json: unsupported type: map[interface {}]interface {}
		dataBytes, err := json_.Marshal(data)
		if err != nil {
			return nil, err
		}

		// trick: transfer data to the same format as def, that is proto.Message
		if def == nil {
			err = protojson.Unmarshal(dataBytes, dataProto)
			if err != nil {
				return nil, err
			}
			return dataProto, nil
		}

		// trick: transfer data to the same format as def, that is proto.Message
		// TODO replace merge trick below of merge option for protojson.Unmarshal
		// See https://github.com/protocolbuffers/protobuf/issues/8263
		defBytes, err := protojson.Marshal(def,
			protojson.WithMarshalUseProtoNames(true), // compatible with TextName
		)
		if err != nil {
			return nil, err
		}

		v := viper.New()
		v.SetConfigType("json")
		err = v.MergeConfig(bytes.NewReader(defBytes))
		if err != nil {
			return nil, err
		}
		err = v.MergeConfig(bytes.NewReader(dataBytes))
		if err != nil {
			return nil, err
		}

		// fix(json): error decoding '': json: unsupported type: map[interface {}]interface {}
		allBytes, err := json_.Marshal(v.AllSettings())
		if err != nil {
			return nil, err
		}
		err = protojson.Unmarshal(allBytes, dataProto)
		if err != nil {
			return nil, err
		}
		return dataProto, nil
	}
}
