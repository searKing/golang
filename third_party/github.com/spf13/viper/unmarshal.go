// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	json_ "github.com/searKing/golang/go/encoding/json"
	strings_ "github.com/searKing/golang/go/strings"
	"github.com/searKing/golang/third_party/google.golang.org/protobuf/encoding/protojson"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
)

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
		c.DecodeHook = UnmarshalProtoMessageHookFunc()
	}
}

// UnmarshalKey takes a single key and unmarshalls it into a Struct.
// use protojson to decode if rawVal is proto.Message
func UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if key == "" {
		return Unmarshal(rawVal, opts...)
	}
	if val, ok := rawVal.(proto.Message); ok {
		opts = append([]viper.DecoderConfigOption{DecodeProtoJsonHook(val)}, opts...)
	}
	return viper.UnmarshalKey(key, rawVal, opts...)
}

func UnmarshalKeyViper(v *viper.Viper, key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if v == nil {
		return UnmarshalKey(key, rawVal, opts...)
	}
	if key == "" {
		return UnmarshalViper(v, rawVal, opts...)
	}
	if val, ok := rawVal.(proto.Message); ok {
		opts = append([]viper.DecoderConfigOption{DecodeProtoJsonHook(val)}, opts...)
	}
	return v.UnmarshalKey(key, rawVal, opts...)
}

func UnmarshalKeys(keys []string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return UnmarshalKey(strings.Join(strings_.SliceTrimEmpty(keys...), "."), rawVal, opts...)
}

func UnmarshalKeysViper(v *viper.Viper, keys []string, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	return UnmarshalKeyViper(v, strings.Join(strings_.SliceTrimEmpty(keys...), "."), rawVal, opts...)
}

// Unmarshal unmarshalls the config into a Struct. Make sure that the tags
// on the fields of the structure are properly set.
// use protojson to decode if rawVal is proto.Message
func Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if val, ok := rawVal.(proto.Message); ok {
		opts = append([]viper.DecoderConfigOption{DecodeProtoJsonHook(val)}, opts...)
	}
	return viper.Unmarshal(rawVal, opts...)
}

func UnmarshalViper(v *viper.Viper, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if v == nil {
		return Unmarshal(rawVal, opts...)
	}
	if val, ok := rawVal.(proto.Message); ok {
		opts = append([]viper.DecoderConfigOption{DecodeProtoJsonHook(val)}, opts...)
	}
	return v.Unmarshal(rawVal, opts...)
}

// UnmarshalExact unmarshals the config into a Struct, erroring if a field is nonexistent
// in the destination struct.
// use protojson to decode if rawVal is proto.Message
func UnmarshalExact(rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if val, ok := rawVal.(proto.Message); ok {
		opts = append([]viper.DecoderConfigOption{DecodeProtoJsonHook(val)}, opts...)
	}
	return viper.UnmarshalExact(rawVal, opts...)
}

func UnmarshalExactViper(v *viper.Viper, rawVal interface{}, opts ...viper.DecoderConfigOption) error {
	if v == nil {
		return UnmarshalExact(rawVal, opts...)
	}
	if val, ok := rawVal.(proto.Message); ok {
		opts = append([]viper.DecoderConfigOption{DecodeProtoJsonHook(val)}, opts...)
	}
	return v.UnmarshalExact(rawVal, opts...)
}

// UnmarshalProtoMessageHookFunc returns a DecodeHookFunc that converts
// root struct to config.ViperProto.
// Trick of protobuf, which generates json tag only
func UnmarshalProtoMessageHookFunc() mapstructure.DecodeHookFunc {
	return func(src reflect.Type, dst reflect.Type, data interface{}) (interface{}, error) {
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
		err = protojson.Unmarshal(dataBytes, dataProto)
		if err != nil {
			return nil, err
		}
		return dataProto, err
	}
}
