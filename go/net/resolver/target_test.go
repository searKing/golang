// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package resolver_test

import (
	"net/url"
	"testing"

	"github.com/searKing/golang/go/net/resolver"
)

func TestParseTarget(t *testing.T) {
	defScheme := resolver.GetDefaultScheme()
	for i, test := range []struct {
		target    string
		badScheme bool
		want      resolver.Target
	}{
		// No scheme is specified.
		{target: "", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: ""}},
		{target: "://", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: "://"}},
		{target: ":///", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: ":///"}},
		{target: "://a/", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: "://a/"}},
		{target: ":///a", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: ":///a"}},
		{target: "://a/b", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: "://a/b"}},
		{target: "/", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: "/"}},
		{target: "a/b", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: "a/b"}},
		{target: "a//b", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: "a//b"}},
		{target: "google.com", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: "google.com"}},
		{target: "google.com/?a=b", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: "google.com/"}},
		{target: "/unix/socket/address", badScheme: true, want: resolver.Target{Scheme: "", Authority: "", Endpoint: "/unix/socket/address"}},

		// A scheme is specified.
		{target: "dns:///google.com", want: resolver.Target{Scheme: "dns", Authority: "", Endpoint: "google.com"}},
		{target: "dns://a.server.com/google.com", want: resolver.Target{Scheme: "dns", Authority: "a.server.com", Endpoint: "google.com"}},
		{target: "dns://a.server.com/google.com/?a=b", want: resolver.Target{Scheme: "dns", Authority: "a.server.com", Endpoint: "google.com/"}},
		{target: "unix:///a/b/c", want: resolver.Target{Scheme: "unix", Authority: "", Endpoint: "a/b/c"}},
		{target: "unix-abstract:a/b/c", want: resolver.Target{Scheme: "unix-abstract", Authority: "", Endpoint: "a/b/c"}},
		{target: "unix-abstract:a b", want: resolver.Target{Scheme: "unix-abstract", Authority: "", Endpoint: "a b"}},
		{target: "unix-abstract:a:b", want: resolver.Target{Scheme: "unix-abstract", Authority: "", Endpoint: "a:b"}},
		{target: "unix-abstract:a-b", want: resolver.Target{Scheme: "unix-abstract", Authority: "", Endpoint: "a-b"}},
		{target: "unix-abstract:/ a///://::!@#$%25^&*()b", want: resolver.Target{Scheme: "unix-abstract", Authority: "", Endpoint: " a///://::!@"}},
		{target: "unix-abstract:passthrough:abc", want: resolver.Target{Scheme: "unix-abstract", Authority: "", Endpoint: "passthrough:abc"}},
		{target: "unix-abstract:unix:///abc", want: resolver.Target{Scheme: "unix-abstract", Authority: "", Endpoint: "unix:///abc"}},
		{target: "unix-abstract:///a/b/c", want: resolver.Target{Scheme: "unix-abstract", Authority: "", Endpoint: "a/b/c"}},
		{target: "unix-abstract:///", want: resolver.Target{Scheme: "unix-abstract", Authority: "", Endpoint: ""}},
		{target: "passthrough:///unix:///a/b/c", want: resolver.Target{Scheme: "passthrough", Authority: "", Endpoint: "unix:///a/b/c"}},

		// Cases for `scheme:absolute-path`.
		{target: "dns:/a/b/c", want: resolver.Target{Scheme: "dns", Authority: "", Endpoint: "a/b/c"}},
	} {
		target := test.target
		if test.badScheme {
			target = defScheme + ":///" + target
		}
		url, err := url.Parse(target)
		if err != nil {
			t.Fatalf("Unexpected error parsing URL: %v", err)
		}
		if test.badScheme {
			url.Scheme = ""
		}
		test.want.URL = *url
		got := resolver.ParseTarget(test.target)
		if got != test.want {
			t.Errorf("#%d: ParseTarget(%q) = %+v, want %+v", i, test.target, got, test.want)
		}
	}
}
