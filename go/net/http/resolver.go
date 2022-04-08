// Copyright 2022 The searKing Author. All rights reserved.
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
	if target == "" {
		return nil
	}
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

// RoundTripperWithTarget wraps http.RoundTripper with request url replaced by target resolved by resolver.
// target is as like gRPC Naming for service discovery.
func RoundTripperWithTarget(rt http.RoundTripper, target string, replaceHostInRequest bool) http.RoundTripper {
	if target == "" {
		return rt
	}
	return RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		err = RequestWithTarget(req, target, replaceHostInRequest)
		if err != nil {
			return nil, err
		}
		return rt.RoundTrip(req)
	})
}

// TransportWithProxyTarget wraps http.RoundTripper with request url replaced by target resolved by resolver.
// target is as like gRPC Naming for service discovery.
func TransportWithProxyTarget(t *http.Transport, target string) *http.Transport {
	t.Proxy = ProxyFuncWithTargetOrDefault(target, t.Proxy)
	return t
}
