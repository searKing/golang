// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper_test

import (
	"fmt"
	"testing"

	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/searKing/golang/third_party/github.com/spf13/viper/testdata"
	"github.com/spf13/viper"
)

func TestDecodeProtoJsonHook(t *testing.T) {
	viper.Set("credentials", map[int]string{1: "foo"})

	var got = testdata.Config{
		Credentials: map[int64]string{2: "bar"},
	}

	err := viper.Unmarshal(&got, viper_.DecodeProtoJsonHook(&got))
	if err != nil {
		t.Fatalf("unable to decode into struct, %v", err)
	}
	want := &testdata.Config{
		Credentials: map[int64]string{1: "foo", 2: "bar"},
	}
	if fmt.Sprintf("%v", want.Credentials) != fmt.Sprintf("%v", got.Credentials) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestUnmarshalViper(t *testing.T) {
	viper.Set("credentials", map[int]string{1: "foo"})

	var got = testdata.Config{
		Credentials: map[int64]string{2: "bar"},
	}
	err := viper_.Unmarshal(&got)
	if err != nil {
		t.Fatalf("unable to decode into struct, %v", err)
	}
	want := &testdata.Config{
		Credentials: map[int64]string{1: "foo", 2: "bar"},
	}
	if fmt.Sprintf("%v", want.Credentials) != fmt.Sprintf("%v", got.Credentials) {
		t.Errorf("got %v want %v", got, want)
	}
}
