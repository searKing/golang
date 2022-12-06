// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"

	url_ "github.com/searKing/golang/go/net/url"
)

// RequestWithTarget replace Host in url.Url by resolver.Target
// replace Host in req if replaceHostInRequest is true
func RequestWithTarget(req *http.Request, target string, replaceHostInRequest bool) error {
	if target == "" {
		return nil
	}
	u2, err := url_.ResolveWithTarget(req.Context(), req.URL, target)
	if err != nil {
		return err
	}
	host := req.URL.Host
	req.URL = u2
	if replaceHostInRequest {
		req.Host = u2.Host
	} else {
		req.Host = host
	}
	return nil
}
