// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"crypto/tls"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
)

// CloneURLValues returns a deep copy of v, or nil if v is nil.
func CloneURLValues(v url.Values) url.Values {
	return v.Clone()
}

// CloneURL returns a deep copy of u, or nil if u is nil.
func CloneURL(u *url.URL) *url.URL {
	return u.Clone()
}

// CloneMultipartForm returns a deep copy of f, or nil if f is nil.
func CloneMultipartForm(f *multipart.Form) *multipart.Form {
	if f == nil {
		return nil
	}
	f2 := &multipart.Form{
		Value: (map[string][]string)(http.Header(f.Value).Clone()),
	}
	if f.File != nil {
		m := make(map[string][]*multipart.FileHeader, len(f.File))
		for k, vv := range f.File {
			vv2 := make([]*multipart.FileHeader, len(vv))
			for i, v := range vv {
				vv2[i] = CloneMultipartFileHeader(v)
			}
			m[k] = vv2
		}
		f2.File = m
	}
	return f2
}

// CloneMultipartFileHeader returns a copy of fh with its Header deep copied, or nil if fh is nil.
func CloneMultipartFileHeader(fh *multipart.FileHeader) *multipart.FileHeader {
	if fh == nil {
		return nil
	}
	fh2 := new(*fh)
	fh2.Header = textproto.MIMEHeader(http.Header(fh.Header).Clone())
	return fh2
}

// CloneOrMakeHeader invokes Header.Clone but if the
// result is nil, it'll instead make and return a non-nil Header.
func CloneOrMakeHeader(hdr http.Header) http.Header {
	clone := hdr.Clone()
	if clone == nil {
		clone = make(http.Header)
	}
	return clone
}

// CloneTLSConfig returns a shallow clone of cfg, or a new zero tls.Config if
// cfg is nil. This is safe to call even if cfg is in active use by a TLS
// client or server.
func CloneTLSConfig(cfg *tls.Config) *tls.Config {
	if cfg == nil {
		return &tls.Config{}
	}
	return cfg.Clone()
}
