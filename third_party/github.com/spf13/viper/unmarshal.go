// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"bytes"
	"reflect"

	"github.com/mitchellh/mapstructure"
	json_ "github.com/searKing/golang/go/encoding/json"
	"github.com/searKing/golang/third_party/google.golang.org/protobuf/encoding/protojson"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"
)

// UnmarshalProtoMessageByJsonpb returns the latest config viper proto
func UnmarshalProtoMessageByJsonpb(viper_ *viper.Viper, v proto.Message, opts ...viper.DecoderConfigOption) error {
	if viper_ == nil { // nop for nil source
		return nil
	}
	// config file -> ViperProto
	var opt []viper.DecoderConfigOption
	opt = append(opt, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "json" // trick of protobuf, which generates json tag only
		decoderConfig.WeaklyTypedInput = true
		decoderConfig.DecodeHook = UnmarshalProtoMessageByJsonpbHookFunc(v)
	})
	opt = append(opt, opts...)
	return viper_.Unmarshal(v, opt...)
}

// UnmarshalProtoMessageByJsonpbHookFunc returns a DecodeHookFunc that converts
// root struct to config.ViperProto.
// Trick of protobuf, which generates json tag only
func UnmarshalProtoMessageByJsonpbHookFunc(def proto.Message) mapstructure.DecodeHookFunc {
	return func(src reflect.Type, dst reflect.Type, data interface{}) (interface{}, error) {
		protoBytes, err := protojson.Marshal(def,
			protojson.WithMarshalUseProtoNames(true), // compatible with TextName
		)
		if err != nil {
			return nil, err
		}

		// fix(json): error decoding '': json: unsupported type: map[interface {}]interface {}
		dataBytes, err := json_.Marshal(data)
		if err != nil {
			return nil, err
		}

		{ // trick: transfer data to the same format as def, that is proto.Message
			dataProto := proto.Clone(def)
			err = protojson.Unmarshal(dataBytes, dataProto)
			if err != nil {
				return nil, err
			}
			dataBytes, err = protojson.Marshal(dataProto,
				protojson.WithMarshalUseProtoNames(true), // compatible with TextName
			)
			if err != nil {
				return nil, err
			}
		}

		v := viper.New()
		v.SetConfigType("json")
		err = v.MergeConfig(bytes.NewReader(protoBytes))
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

		err = protojson.Unmarshal(allBytes, def)
		if err != nil {
			return nil, err
		}
		return def, nil
	}
}