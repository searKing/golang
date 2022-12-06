// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpproxy_test

import (
	"context"
	"testing"

	"github.com/searKing/golang/go/net/http/httpproxy"
)

func TestWithProxy(t *testing.T) {
	ctx := context.Background()
	want := &httpproxy.Proxy{
		ProxyUrl:    "socks5://127.0.0.1:8080",
		ProxyTarget: "dns:///want.example.com",
	}
	ctx = httpproxy.WithProxy(ctx, want)
	got := httpproxy.ContextProxy(ctx)
	if got == nil || got.ProxyUrl != want.ProxyUrl || got.ProxyTarget != want.ProxyTarget {
		t.Errorf("got %v; want %v", got, want)
	}
}
