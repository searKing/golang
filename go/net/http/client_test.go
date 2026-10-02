// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	http_ "github.com/searKing/golang/go/net/http"
	_ "github.com/searKing/golang/go/net/resolver/passthrough"
)

func TestNewClientWithTarget(t *testing.T) {
	var gotHost string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHost = r.Host
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client := http_.NewClientWithTarget(ts.Listener.Addr().String())
	resp, err := client.Get("http://logical.invalid/healthz")
	if err != nil {
		t.Fatalf("GET logical.invalid: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", resp.StatusCode)
	}
	if gotHost != "logical.invalid" {
		t.Fatalf("host = %q, want logical.invalid", gotHost)
	}
}
