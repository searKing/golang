// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"fmt"
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
	host := req.URL.Host
	req.URL = u2
	if replaceHostInRequest {
		req.Host = u2.Host
	} else {
		req.Host = host
	}
	return nil
}

// ProxyFuncWithTargetOrDefault builds a proxy function from the given string, which should
// represent a target that can be used as a proxy. It performs basic
// sanitization of the Target and returns any error encountered.
// proxyUrl is proxy's url, like sock5://127.0.0.1:8080
// proxyTarget is proxy's addr, replace the HOST in proxyUrl if not empty
func ProxyFuncWithTargetOrDefault(proxyUrl string, proxyTarget string, def func(req *http.Request) (*url.URL, error)) func(req *http.Request) (*url.URL, error) {
	if proxyTarget == "" {
		return def
	}
	return func(req *http.Request) (*url.URL, error) {
		proxy, err := parseProxy(proxyUrl)
		if err != nil {
			return nil, err
		}
		if proxy == nil || proxyTarget == "" {
			return proxy, nil
		}

		address, err := resolver.ResolveOneAddr(req.Context(), proxyTarget)
		if err != nil {
			return nil, err
		}
		if address.Addr != "" {
			proxy.Host = address.Addr
		}
		return proxy, nil
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
// proxyUrl is proxy's url, like socks5://127.0.0.1:8080
// proxyTarget is proxy's addr, replace the HOST in proxyUrl if not empty
func TransportWithProxyTarget(t *http.Transport, proxyUrl string, proxyTarget string) *http.Transport {
	t.Proxy = ProxyFuncWithTargetOrDefault(proxyUrl, proxyTarget, t.Proxy)
	return t
}

func parseProxy(proxy string) (*url.URL, error) {
	if proxy == "" {
		return nil, nil
	}

	proxyURL, err := url.Parse(proxy)
	if err != nil ||
		(proxyURL.Scheme != "http" &&
			proxyURL.Scheme != "https" &&
			proxyURL.Scheme != "socks5") {
		// proxy was bogus. Try prepending "http://" to it and
		// see if that parses correctly. If not, we fall
		// through and complain about the original one.
		if proxyURL, err := url.Parse("http://" + proxy); err == nil {
			return proxyURL, nil
		}
	}
	if err != nil {
		return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
	}
	return proxyURL, nil
}
