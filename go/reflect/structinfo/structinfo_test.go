// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structinfo

import (
	"reflect"
	"testing"
)

type inlined struct {
	B int `bson:"b"`
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
