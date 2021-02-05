// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"reflect"

	"github.com/searKing/golang/third_party/github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	json_ "github.com/searKing/golang/go/encoding/json"
)

// Unmarshal returns the latest config viper proto
func UnmarshalProtoMessageByJsonpb(viper_ *viper.Viper, v proto.Message, opts ...viper.DecoderConfigOption) error {
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
func UnmarshalProtoMessageByJsonpbHookFunc(v proto.Message) mapstructure.DecodeHookFunc {
	return func(src reflect.Type, dst reflect.Type, data interface{}) (interface{}, error) {
		// Convert it by parsing
		dataBytes, err := json_.Marshal(data)
		if err != nil {
			return nil, err
		}

		// apply protobuf check
		err = jsonpb.Unmarshal(dataBytes, v, jsonpb.WithUnmarshalAllowUnknownFields(true))
		if err != nil {
			return data, err
		}
		return v, nil
	}
}
