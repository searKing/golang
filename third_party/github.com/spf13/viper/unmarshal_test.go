// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-viper/mapstructure/v2"
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/searKing/golang/third_party/github.com/spf13/viper/testdata"
	"github.com/spf13/viper"
)

func TestDecodeProtoJsonHook(t *testing.T) {
	v := viper.New()
	v.Set("credentials", map[int]string{1: "foo"})

	var got = testdata.Config{
		Credentials: map[int64]string{2: "bar"},
	}

	err := v.Unmarshal(&got, viper_.DecodeProtoJsonHook(&got))
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
	v := viper.New()
	v.Set("credentials", map[int]string{1: "foo"})

	var got = &testdata.Config{
		Credentials: map[int64]string{2: "bar"},
	}
	err := viper_.UnmarshalViper(v, &got)
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

func TestUnmarshalWithOptions(t *testing.T) {
	type Server struct {
		Port int `cfg:"listen_port"`
	}
	viper.Set("server", map[string]any{"listen_port": 8080})
	t.Cleanup(viper.Reset)
	withTag := func(c *mapstructure.DecoderConfig) { c.TagName = "cfg" }

	var gotKeys Server
	if err := viper_.UnmarshalKeys([]string{"server"}, &gotKeys, withTag); err != nil {
		t.Fatalf("UnmarshalKeys: %v", err)
	}
	if gotKeys.Port != 8080 {
		t.Errorf("UnmarshalKeys ignores opts: got port %d, want %d", gotKeys.Port, 8080)
	}

	var got struct {
		Server Server `cfg:"server"`
	}
	if err := viper_.Unmarshal(&got, withTag); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Server.Port != 8080 {
		t.Errorf("Unmarshal ignores opts: got port %d, want %d", got.Server.Port, 8080)
	}
}

func TestUnmarshalKeysViperDecodeHook(t *testing.T) {
	type Server struct {
		Name string
	}
	upper := func(f reflect.Type, _ reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		return strings.ToUpper(data.(string)), nil
	}
	v := viper.NewWithOptions(viper.WithDecodeHook(upper))
	v.Set("server", map[string]any{"name": "foo"})

	for name, fn := range map[string]func(*viper.Viper, []string, any, ...viper.DecoderConfigOption) error{
		"UnmarshalKeysViper":      viper_.UnmarshalKeysViper,
		"UnmarshalKeysExactViper": viper_.UnmarshalKeysExactViper,
	} {
		var got Server
		if err := fn(v, []string{"server"}, &got); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if got.Name != "FOO" {
			t.Errorf("%s ignores decode hook of viper: got name %q, want %q", name, got.Name, "FOO")
		}
	}
}

func TestUnmarshalKeysViperDefaultDecodeHook(t *testing.T) {
	type Server struct {
		Timeout time.Duration
		Hosts   []string
	}
	v := viper.New()
	v.Set("server", map[string]any{"timeout": "1m", "hosts": "a,b"})

	var got Server
	if err := viper_.UnmarshalKeysViper(v, []string{"server"}, &got); err != nil {
		t.Fatalf("UnmarshalKeysViper: %v", err)
	}
	if want := (Server{Timeout: time.Minute, Hosts: []string{"a", "b"}}); !reflect.DeepEqual(got, want) {
		t.Errorf("UnmarshalKeysViper = %+v, want %+v", got, want)
	}
}

func TestUnmarshalKeyViper(t *testing.T) {
	v := viper.New()
	v.Set("credentials", map[int]string{1: "foo"})

	var got = &testdata.Config{
		Credentials: map[int64]string{2: "bar"},
	}
	err := viper_.UnmarshalKeysViper(v, []string{"credentials"}, &got.Credentials)
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
