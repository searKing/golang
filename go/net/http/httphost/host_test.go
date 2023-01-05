// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httphost_test

import (
	"context"
	"testing"

	"github.com/searKing/golang/go/net/http/httphost"
)

func TestWithProxy(t *testing.T) {
	ctx := context.Background()
	want := &httphost.Host{
		HostTarget: "dns:///want.example.com",
	}
	ctx = httphost.WithHost(ctx, want)
	got := httphost.ContextHost(ctx)
	if got == nil || got.HostTarget != want.HostTarget {
		t.Errorf("got %v; want %v", got, want)
	}
}
