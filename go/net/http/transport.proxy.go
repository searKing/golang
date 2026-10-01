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

// ProxyFuncFromContextOrEnvironment resolves the proxy URL for an HTTP
// request, supporting both dynamic proxy discovery and fallback to
// environment variables.
//
// If the request context contains a proxy configuration via
// RequestWithProxyTarget, it is used. Otherwise, the function falls back to
// http.ProxyFromEnvironment, which respects HTTP_PROXY, HTTPS_PROXY, and
// NO_PROXY environment variables.
//
// When ProxyTarget is set in the context, it is resolved through the
// configured resolver and the host portion of the proxy URL is updated with
// the resolved address. This enables dynamic proxy discovery.
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

// DefaultTransportWithDynamicProxy is an http.RoundTripper that supports
// both static and dynamically resolved HTTP/HTTPS/SOCKS5 proxies.
//
// It uses ProxyFuncFromContextOrEnvironment to determine the proxy URL for
// each request. Proxies can be configured via:
//   - RequestWithProxyTarget (dynamic resolution)
//   - HTTP_PROXY, HTTPS_PROXY, NO_PROXY environment variables
//
// The transport maintains connection pooling and configures reasonable
// timeout and TLS handshake defaults suitable for proxied requests.
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

// DefaultClientWithDynamicProxy is an http.Client that uses DefaultTransportWithDynamicProxy.
var DefaultClientWithDynamicProxy = &http.Client{
	Transport: DefaultTransportWithDynamicProxy,
}

// ProxyFuncWithTargetOrDefault returns a proxy function that applies dynamic
// proxy resolution when configured, otherwise falls back to a default proxy
// function.
//
// If fixedProxyUrl is empty, def is returned unchanged.
//
// If fixedProxyUrl is provided but fixedProxyTarget is empty, a function
// returning the static proxy URL is returned (no dynamic resolution).
//
// If both fixedProxyUrl and fixedProxyTarget are provided, a function is
// returned that resolves fixedProxyTarget through the resolver and updates
// the proxy URL's host with the resolved address. This enables dynamic proxy
// discovery.
//
// Parameters:
//   - fixedProxyUrl: The proxy URL, e.g. "http://127.0.0.1:8080" or
//     "socks5://proxy.example.com:1080"
//   - fixedProxyTarget: The resolver target for dynamic proxy discovery, e.g.
//     "proxy-service" or "proxy-lb"
//   - def: The fallback proxy function, used when fixedProxyUrl is empty
//
// Returns:
//   - A proxy function suitable for http.Transport.Proxy
//   - An error function if fixedProxyUrl parsing fails
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

// TransportWithProxyTarget wraps an http.Transport to apply dynamic proxy
// resolution.
//
// It updates the transport's Proxy field with a proxy function that supports
// both dynamic resolution and fallback to the transport's current proxy
// configuration.
//
// Parameters:
//   - t: The http.Transport to wrap
//   - fixedProxyUrl: The proxy URL, e.g. "http://127.0.0.1:8080"
//   - fixedProxyTarget: The resolver target for dynamic discovery, e.g.
//     "proxy-service"
//
// Returns:
//   - The modified transport (pointer to the same object)
func TransportWithProxyTarget(t *http.Transport, fixedProxyUrl string, fixedProxyTarget string) *http.Transport {
	t.Proxy = ProxyFuncWithTargetOrDefault(fixedProxyUrl, fixedProxyTarget, t.Proxy)
	return t
}
