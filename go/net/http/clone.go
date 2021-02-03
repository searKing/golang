// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	_ "unsafe" // for go:linkname

	"crypto/tls"
	"mime/multipart"
	"net/http"
	"net/url"
)

//go:linkname CloneURLValues net/http.cloneURLValues
func CloneURLValues(v url.Values) url.Values

//go:linkname CloneURL net/http.cloneURL
func CloneURL(u *url.URL) *url.URL

//go:linkname CloneMultipartForm net/http.cloneMultipartForm
func CloneMultipartForm(f *multipart.Form) *multipart.Form

//go:linkname CloneMultipartFileHeader net/http.cloneMultipartFileHeader
func CloneMultipartFileHeader(fh *multipart.FileHeader) *multipart.FileHeader

// CloneOrMakeHeader invokes Header.Clone but if the
// result is nil, it'll instead make and return a non-nil Header.
//go:linkname CloneOrMakeHeader net/http.cloneOrMakeHeader
func CloneOrMakeHeader(hdr http.Header) http.Header

// CloneTLSConfig returns a shallow clone of cfg, or a new zero tls.Config if
// cfg is nil. This is safe to call even if cfg is in active use by a TLS
// client or server.
//go:linkname CloneTLSConfig net/http.cloneTLSConfig
func CloneTLSConfig(cfg *tls.Config) *tls.Config
