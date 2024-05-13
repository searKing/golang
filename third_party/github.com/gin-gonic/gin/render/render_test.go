// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package render_test

import (
	"net/http/httptest"
	"testing"

	"github.com/searKing/golang/third_party/github.com/gin-gonic/gin/render"
	"github.com/searKing/golang/third_party/github.com/gin-gonic/gin/render/internal"
)

func TestRenderJSONPB(t *testing.T) {
	w := httptest.NewRecorder()

	data := internal.Data{
		Foo:  "bar",
		Html: "<b>",
	}

	(render.JSONPB{Data: &data}).WriteContentType(w)
	if "application/json; charset=utf-8" != w.Header().Get("Content-Type") {
		t.Errorf("got %q, want %q", w.Body.String(), `{"foo":"bar","html":"\u003cb\u003e"}`)
	}

	err := (render.JSONPB{Data: &data}).Render(w)
	if err != nil {
		t.Errorf("data=%q; %v", data.String(), err)
		return
	}
	want := `{"foo":"bar","html":"\u003cb\u003e"}`
	if want != w.Body.String() {
		t.Errorf("got %s, want %s", w.Body.String(), want)
	}
	if "application/json; charset=utf-8" != w.Header().Get("Content-Type") {
		t.Errorf("got %q, want %q", w.Header().Get("Content-Type"), "application/json; charset=utf-8")
	}
}
