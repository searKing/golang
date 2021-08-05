// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver_test

import (
	"context"
	"testing"
	"time"

	"github.com/searKing/golang/go/net/resolver"
	_ "github.com/searKing/golang/go/net/resolver/dns"
	_ "github.com/searKing/golang/go/net/resolver/passthrough"
	_ "github.com/searKing/golang/go/net/resolver/unix"
)

func TestResolveAddr(t *testing.T) {

	testCases := []struct {
		target string
		expect string
	}{
		{
			target: "passthrough://a.server.com/google.com",
			expect: "google.com",
		},
		//{
		//	target: "dns://www/google.com",
		//	expect: "google.com",
		//},
		{
			target: "unix:///a/b/c",
			expect: "/a/b/c",
		},
	}

	for i, test := range testCases {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			var target = test.target
			addr, err := resolver.ResolveOneAddr(ctx, target)
			if err != nil {
				t.Fatalf("#%d: ResolveOneAddr failed: %s", i, err)
			}
			if addr.Addr != "google.com" {
				t.Fatalf("#%d: expected %s got %s", i, test.expect, addr.Addr)
			}
			t.Logf("#%d: addr: %s", i, addr.Addr)
		}()
	}
}
