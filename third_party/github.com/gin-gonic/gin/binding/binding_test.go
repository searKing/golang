// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binding_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin/binding"
	binding_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin/binding"
	"github.com/searKing/golang/third_party/github.com/gin-gonic/gin/binding/internal"
)

func TestBindingJSONPB(t *testing.T) {
	testBodyBinding(t,
		binding_.JSONPB, "jsonpb",
		"/", "/",
		`{"foo": "bar"}`, `{"bar": "foo"}`)
}

func testBodyBinding(t *testing.T, b binding.Binding, name, path, badPath, body, badBody string) {
	if name != b.Name() {
		t.Errorf("got %q, want %q", b.Name(), name)
		return
	}

	data := internal.Data{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &data)
	if err != nil {
		t.Errorf("data=%q; %v", data.String(), err)
		return
	}
	if "bar" != data.Foo {
		t.Errorf("got %q, want %q", data.Foo, "bar")
		return
	}

	data = internal.Data{}
	req = requestWithBody("POST", badPath, badBody)
	err = binding_.JSONPB.Bind(req, &data)
	if err != nil {
		t.Errorf("data=%q; %v", data.String(), err)
		return
	}
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}
