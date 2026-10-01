// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpproxy

import (
	"context"
	"fmt"
	"net/url"

	"github.com/searKing/golang/go/net/resolver"
)

// unique type to prevent assignment.
type proxyContextKey struct{}

// Proxy describes a dynamically resolved proxy server used by an HTTP transport.
//
// ProxyUrl defines the proxy endpoint and its scheme, such as:
//
//	http://proxy.example.com:8080
//	https://proxy.example.com:8443
//	socks5://127.0.0.1:1080
//
// If ProxyTarget is set, the resolver resolves it to a concrete address and
// replaces the host portion of ProxyUrl before the request is sent. The scheme
// remains unchanged, so the transport continues to use the same proxy protocol
// (HTTP, HTTPS, or SOCKS5).
//
// Unlike httphost.Host, Proxy is for actual proxy configuration. The transport
// will use the proxy URL to initiate an HTTP proxy, CONNECT tunnel, or SOCKS5
// connection to the configured proxy server. It does not rewrite the final
// destination service address of the user request itself.
type Proxy struct {
	// ProxyUrl is the proxy endpoint URL, including scheme and host.
	//
	// Examples:
	//   http://127.0.0.1:8080
	//   https://proxy.example.com:8443
	//   socks5://127.0.0.1:1080
	ProxyUrl string

	// ProxyTarget is a resolver target used to locate the proxy server
	// dynamically. It is not the final business target of the HTTP request.
	//
	// Example:
	//   ProxyTarget: "proxy-service" -> resolves to "10.0.0.20:3128"
	ProxyTarget string

	// ProxyAddrResolved stores the address selected by the resolver for
	// ProxyTarget. It is used for reporting and error wrapping.
	ProxyAddrResolved resolver.Address
}

// ContextProxy returns the dynamic proxy configuration associated with ctx.
//
// It returns nil when ctx does not contain a Proxy value.
func ContextProxy(ctx context.Context) *Proxy {
	proxy, _ := ctx.Value(proxyContextKey{}).(*Proxy)
	return proxy
}

// WithProxy returns a child context carrying the dynamic proxy configuration.
//
// The configuration is consumed by http.Transport's Proxy function when the
// request's proxy URL is selected. If proxy is nil, WithProxy panics.
func WithProxy(ctx context.Context, proxy *Proxy) context.Context {
	if proxy == nil {
		panic("nil proxy")
	}
	return context.WithValue(ctx, proxyContextKey{}, proxy)
}

func ParseProxyUrl(proxy string) (*url.URL, error) {
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
