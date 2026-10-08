// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structinfo

import (
	"maps"
	"reflect"
	"testing"
)

type inlined struct {
	B int `bson:"b"`
}

func TestGetStructInfo(t *testing.T) {
	type T struct {
		A int
		B int `bson:"b,omitempty"`
	}
	sinfo, err := GetStructInfo(reflect.TypeFor[T]())
	if err != nil {
		t.Fatalf("GetStructInfo() error = %v", err)
	}
	tests := []fieldInfo{
		{Key: "a", Num: 0},
		{Key: "b", Num: 1, Tags: map[string]string{"OmitEmpty": "omitempty"}},
	}
	for _, want := range tests {
		v, ok := sinfo.fieldsLRU.Find(want.Key)
		got, _ := v.(fieldInfo)
		if !ok || got.Key != want.Key || got.Num != want.Num || !maps.Equal(got.Tags, want.Tags) {
			t.Errorf("field %q = %v, %t; want %v", want.Key, v, ok, want)
		}
	}
}

func TestGetStructInfoInline(t *testing.T) {
	type T struct {
		A       int
		Inlined inlined `bson:",inline"`
	}
	sinfo, err := GetStructInfo(reflect.TypeFor[T]())
	if err != nil {
		t.Fatalf("GetStructInfo() error = %v", err)
	}
	v, ok := sinfo.fieldsLRU.Find("b")
	if got, _ := v.(fieldInfo); !ok || got.Key != "b" {
		t.Errorf("inline field b = %v, %t; want fieldInfo of b", v, ok)
	}
}
