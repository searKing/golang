// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http_test

import (
	"crypto/tls"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"
	"testing"

	http_ "github.com/searKing/golang/go/net/http"
)

func TestCloneNil(t *testing.T) {
	if got := http_.CloneURLValues(nil); got != nil {
		t.Errorf("CloneURLValues(nil) = %v, want nil", got)
	}
	if got := http_.CloneURL(nil); got != nil {
		t.Errorf("CloneURL(nil) = %v, want nil", got)
	}
	if got := http_.CloneMultipartForm(nil); got != nil {
		t.Errorf("CloneMultipartForm(nil) = %v, want nil", got)
	}
	if got := http_.CloneMultipartFileHeader(nil); got != nil {
		t.Errorf("CloneMultipartFileHeader(nil) = %v, want nil", got)
	}
	if got := http_.CloneOrMakeHeader(nil); got == nil || len(got) != 0 {
		t.Errorf("CloneOrMakeHeader(nil) = %#v, want empty non-nil Header", got)
	}
	if got := http_.CloneTLSConfig(nil); got == nil {
		t.Errorf("CloneTLSConfig(nil) = nil, want zero tls.Config")
	}
}

func TestCloneURLValues(t *testing.T) {
	v := url.Values{"a": {"1", "2"}}
	c := http_.CloneURLValues(v)
	if !reflect.DeepEqual(c, v) {
		t.Fatalf("CloneURLValues = %v, want %v", c, v)
	}
	c["a"][0] = "x"
	c.Add("b", "3")
	if want := (url.Values{"a": {"1", "2"}}); !reflect.DeepEqual(v, want) {
		t.Errorf("original modified via clone: %v, want %v", v, want)
	}
}

func TestCloneURL(t *testing.T) {
	u := &url.URL{Scheme: "https", User: url.UserPassword("user", "pass"), Host: "example.com", Path: "/p"}
	c := http_.CloneURL(u)
	if c == u || c.User == u.User {
		t.Fatalf("CloneURL shares memory with the original")
	}
	if c.String() != u.String() {
		t.Fatalf("CloneURL = %q, want %q", c, u)
	}
	c.Host = "other.com"
	c.User = url.User("other")
	if got, want := u.String(), "https://user:pass@example.com/p"; got != want {
		t.Errorf("original modified via clone: %q, want %q", got, want)
	}
}

func TestCloneMultipartForm(t *testing.T) {
	f := &multipart.Form{
		Value: map[string][]string{"k": {"v"}},
		File: map[string][]*multipart.FileHeader{
			"f": {{Filename: "a.txt", Header: textproto.MIMEHeader{"Content-Type": {"text/plain"}}, Size: 1}},
		},
	}
	c := http_.CloneMultipartForm(f)
	if !reflect.DeepEqual(c, f) {
		t.Fatalf("CloneMultipartForm = %+v, want %+v", c, f)
	}
	c.Value["k"][0] = "x"
	c.File["f"][0].Filename = "b.txt"
	c.File["f"][0].Header.Set("Content-Type", "x")
	if f.Value["k"][0] != "v" || f.File["f"][0].Filename != "a.txt" || f.File["f"][0].Header.Get("Content-Type") != "text/plain" {
		t.Errorf("original modified via clone: %+v", f)
	}
}

func TestCloneOrMakeHeader(t *testing.T) {
	h := http.Header{"A": {"1"}}
	c := http_.CloneOrMakeHeader(h)
	c.Set("A", "2")
	if got := h.Get("A"); got != "1" {
		t.Errorf("original modified via clone: %q, want %q", got, "1")
	}
}

func TestCloneTLSConfig(t *testing.T) {
	cfg := &tls.Config{ServerName: "example.com", MinVersion: tls.VersionTLS12}
	c := http_.CloneTLSConfig(cfg)
	if c == cfg || c.ServerName != cfg.ServerName || c.MinVersion != cfg.MinVersion {
		t.Fatalf("CloneTLSConfig = %+v, want a copy of %+v", c, cfg)
	}
	c.ServerName = "other.com"
	if cfg.ServerName != "example.com" {
		t.Errorf("original modified via clone: %q", cfg.ServerName)
	}
}
