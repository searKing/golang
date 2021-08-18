// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"
	"net/url"

	"github.com/searKing/golang/go/net/resolver"
	url_ "github.com/searKing/golang/go/net/url"
)

// RequestWithTarget reset Host in url.Url by resolver.Target
// reset Host in req if replaceHostInRequest is true
func RequestWithTarget(req *http.Request, target string, replaceHostInRequest bool) error {
	u2, err := url_.ResolveWithTarget(req.Context(), req.URL, target)
	if err != nil {
		return err
	}
	req.URL = u2
	if replaceHostInRequest {
		req.Host = u2.Host
	}
	return nil
}

// ProxyFuncWithTargetOrDefault builds a proxy function from the given string, which should
// represent a target that can be used as a proxy. It performs basic
// sanitization of the Target and returns any error encountered.
func ProxyFuncWithTargetOrDefault(target string, def func(req *http.Request) (*url.URL, error)) func(req *http.Request) (*url.URL, error) {
	if target == "" {
		return def
	}
	return func(req *http.Request) (*url.URL, error) {
		reqURL := req.URL
		if target == "" {
			return nil, nil
		}
		address, err := resolver.ResolveOneAddr(req.Context(), target)
		if err != nil {
			return nil, err
		}
		tgt := resolver.ParseTarget(target, false)

		var proxy url.URL
		if tgt.Scheme == "http" || tgt.Scheme == "https" {
			proxy.Scheme = tgt.Scheme
		} else {
			proxy.Scheme = reqURL.Scheme
		}
		proxy.Host = address.Addr
		return &proxy, nil
	}
}
