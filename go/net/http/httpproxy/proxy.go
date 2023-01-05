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

// Proxy specifies ProxyUrl and ProxyTarget to return a dynamic proxy.
type Proxy struct {
	// ProxyUrl is proxy's url, like socks5://127.0.0.1:8080
	ProxyUrl string
	// ProxyTarget is as like gRPC Naming for proxy service discovery, with Host in ProxyUrl replaced if not empty.
	ProxyTarget string

	// ProxyAddrResolved is the proxy's addr resolved and picked from resolver.
	ProxyAddrResolved resolver.Address
}

// ContextProxy returns the Proxy associated with the
// provided context. If none, it returns nil.
func ContextProxy(ctx context.Context) *Proxy {
	proxy, _ := ctx.Value(proxyContextKey{}).(*Proxy)
	return proxy
}

// WithProxy returns a new context based on the provided parent
// ctx. HTTP client requests made with the returned context will use
// the provided proxy hooks
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
