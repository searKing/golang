// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"reflect"
	"slices"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	maps_ "github.com/searKing/golang/go/exp/maps"
	json_ "github.com/searKing/golang/third_party/github.com/spf13/viper/json"
	protojson_ "github.com/searKing/golang/third_party/google.golang.org/protobuf/encoding/protojson"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
)

// decode is a wrapper around mapstructure.Decode that mimics the WeakDecode functionality.
//
// Code copied from decode in viper.go of github.com/spf13/viper v1.21.0, keep in sync when upgrading viper.
func decode(input any, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

// defaultDecoderConfig returns the mapstructure.DecoderConfig that v uses to unmarshal,
// with the decode hook of [viper.WithDecodeHook] and opts applied.
//
// v builds the config only when it unmarshals, so the config is captured by an extra
// option applied last, and is then reset to a zero config for v to decode nothing into a dummy.
func defaultDecoderConfig(v *viper.Viper, output any, opts ...viper.DecoderConfigOption) *mapstructure.DecoderConfig {
	var c mapstructure.DecoderConfig
	var dummy struct{}
	_ = v.UnmarshalKey("", &dummy, append(slices.Clip(opts), func(config *mapstructure.DecoderConfig) {
		c = *config
		*config = mapstructure.DecoderConfig{}
	})...)

	// Do not allow overwriting the output
	c.Result = output
	return &c
}

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
	return UnmarshalKeysViper(viper.GetViper(), keys, rawVal, opts...)
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
	return decode(val, defaultDecoderConfig(v, rawVal, opts...))
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
	config := defaultDecoderConfig(v, rawVal, opts...)
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
		err = protojson_.Unmarshal(dataBytes, dataProto)
		if err != nil {
			return nil, err
		}

		if def == nil {
			return dataProto, nil
		}

		allProto := proto.Clone(def)
		proto.Merge(allProto, dataProto)
		return allProto, nil
	}
}
