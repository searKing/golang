// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

// DefaultLoad load config from file and protos into v, and save to a using file
// load sequence: protos..., file, env, replace if member has been set
// that is, later cfg appeared, higher priority cfg has
func DefaultLoad(viper_ *viper.Viper, v proto.Message, cfgFile string, envPrefix string, protos ...proto.Message) error {
	if err := MergeAll(viper_, cfgFile, envPrefix, protos...); err != nil {
		return err
	}

	return UnmarshalProtoMessageByJsonpb(viper_, v)
}

// LoadGlobalConfig load config from the global Viper instance.
func LoadGlobalConfig(v proto.Message, cfgFile string, envPrefix string, protos ...proto.Message) error {
	return DefaultLoad(viper.GetViper(), v, cfgFile, envPrefix, protos...)
}
