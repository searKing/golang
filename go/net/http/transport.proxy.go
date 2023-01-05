// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/searKing/golang/go/net/http/httpproxy"
	"github.com/searKing/golang/go/net/resolver"
)

// RequestWithProxyTarget returns a shallow copy of r with its context changed
// to ctx, TargetUrl and Host inside. The provided ctx must be non-nil.
// proxyUrl is proxy's url, like socks5://127.0.0.1:8080
// proxyTarget is as like gRPC Naming for proxy service discovery, with Host in TargetUrl replaced if not empty.
func RequestWithProxyTarget(req *http.Request, proxy *httpproxy.Proxy) *http.Request {
	if proxy == nil {
		return req
	}
	return req.WithContext(httpproxy.WithProxy(req.Context(), proxy))
}

// ProxyFuncFromContextOrEnvironment builds a proxy function from the given string, which should
// represent a Target that can be used as a proxy. It performs basic
// sanitization of the Target retrieved in context of Request, and returns any error encountered.
func ProxyFuncFromContextOrEnvironment(req *http.Request) (*url.URL, error) {
	proxy := httpproxy.ContextProxy(req.Context())
	// load proxy from environment if proxy not set
	if proxy == nil || proxy.ProxyUrl == "" {
		return http.ProxyFromEnvironment(req)
	}

	proxyUrl, err := httpproxy.ParseProxyUrl(proxy.ProxyUrl)
	if err != nil {
		return nil, err
	}
	if proxyUrl == nil {
		return nil, nil
	}

	if proxy.ProxyTarget == "" {
		return proxyUrl, nil
	}

	// replace host of proxy if target of proxy if resolved
	address, err := resolver.ResolveOneAddr(req.Context(), proxy.ProxyTarget)
	if err != nil {
		return nil, err
	}
	if address.Addr != "" {
		proxyUrl.Host = address.Addr
	}
	proxy.ProxyAddrResolved = address
	return proxyUrl, nil
}

// DefaultTransportWithDynamicProxy is the default implementation of Transport and is
// used by DefaultClientWithDynamicProxy. It establishes network connections as needed
// and caches them for reuse by subsequent calls. It uses HTTP proxies
// as directed by the ProxyFuncFromContextOrEnvironment, $HTTP_PROXY and $NO_PROXY (or $http_proxy and
// $no_proxy) environment variables.
var DefaultTransportWithDynamicProxy http.RoundTripper = &http.Transport{
	Proxy: ProxyFuncFromContextOrEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// DefaultClientWithDynamicProxy is the default Client with DefaultTransportWithDynamicProxy.
var DefaultClientWithDynamicProxy = &http.Client{
	Transport: DefaultTransportWithDynamicProxy,
}

// ProxyFuncWithTargetOrDefault builds a proxy function from the given string, which should
// represent a Target that can be used as a proxy. It performs basic
// sanitization of the Target and returns any error encountered.
// fixedProxyUrl is proxy's url, like socks5://127.0.0.1:8080
// fixedProxyTarget is as like gRPC Naming for proxy service discovery, with Host in TargetUrl replaced if not empty.
func ProxyFuncWithTargetOrDefault(fixedProxyUrl string, fixedProxyTarget string, def func(req *http.Request) (*url.URL, error)) func(req *http.Request) (*url.URL, error) {
	if fixedProxyUrl == "" {
		return def
	}
	proxy, err := httpproxy.ParseProxyUrl(fixedProxyUrl)
	if err != nil {
		return func(req *http.Request) (*url.URL, error) {
			return nil, err
		}
	}
	if proxy == nil || fixedProxyTarget == "" {
		return func(req *http.Request) (*url.URL, error) {
			return proxy, nil
		}
	}
	return func(req *http.Request) (*url.URL, error) {
		req2 := RequestWithProxyTarget(req, &httpproxy.Proxy{
			ProxyUrl:    fixedProxyUrl,
			ProxyTarget: fixedProxyTarget,
		})
		return ProxyFuncFromContextOrEnvironment(req2)
	}
}

// TransportWithProxyTarget wraps http.RoundTripper with request url replaced by Target resolved by resolver.
// fixedProxyUrl is proxy's url, like socks5://127.0.0.1:8080
// fixedProxyTarget is as like gRPC Naming for proxy service discovery, with Host in TargetUrl replaced if not empty.
func TransportWithProxyTarget(t *http.Transport, fixedProxyUrl string, fixedProxyTarget string) *http.Transport {
	t.Proxy = ProxyFuncWithTargetOrDefault(fixedProxyUrl, fixedProxyTarget, t.Proxy)
	return t
}
