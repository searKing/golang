// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httphost

import (
	"context"
	"fmt"
	"net/url"

	"github.com/searKing/golang/go/net/resolver"
)

// unique type to prevent assignment.
type hostContextKey struct{}

// Host describes how the destination host of an HTTP request is resolved and
// optionally rewritten before the request is sent.
//
// HostTarget is resolved by the configured resolver before the request is sent.
// The resolved address replaces req.URL.Host, which determines the network
// address the underlying transport connects to.
//
// In many service-discovery scenarios, the original URL host is a logical name
// or placeholder, not the real backend address. For example:
//
//	original URL: http://api.example.com/v1
//	HostTarget:   "backend-service" -> resolves to "10.0.0.8:8080"
//
// The transport will connect to "10.0.0.8:8080", while the original logical
// host may still be used for routing or identity checks.
//
// ReplaceHostInRequest controls whether the HTTP Host header is also rewritten
// to match the resolved destination address. When false, only req.URL.Host is
// changed; when true, req.Host is also set to req.URL.Host.
//
// This type is for destination resolution and request rewriting. It is not an
// HTTP forward proxy configuration. It does not implement proxy protocols such
// as HTTP CONNECT, HTTPS proxy tunneling, or SOCKS5.
type Host struct {
	// HostTarget is the resolver target for dynamic destination discovery.
	//
	// It may resolve to a concrete network address such as "10.0.0.8:8080"
	// or another logical service name depending on the underlying resolver.
	HostTarget string

	// ReplaceHostInRequest controls whether the HTTP Host field should be
	// rewritten after HostTarget resolution.
	//
	// When false, only req.URL.Host is updated and the client connects to the
	// resolved backend while the HTTP Host header keeps its original logical
	// value. This is useful for virtual-host routing and tenant-aware
	// applications.
	//
	// When true, req.Host is also set to req.URL.Host after resolution so the
	// HTTP Host header matches the resolved network address.
	ReplaceHostInRequest bool

	// HostTargetAddrResolved stores the address selected by the resolver for
	// HostTarget. It is used for metrics, reporting, and error wrapping.
	HostTargetAddrResolved resolver.Address
}

// ContextHost returns the dynamic destination-host configuration associated
// with ctx.
//
// It returns nil when ctx does not contain a Host value.
func ContextHost(ctx context.Context) *Host {
	host, _ := ctx.Value(hostContextKey{}).(*Host)
	return host
}

// WithHost returns a child context carrying the dynamic destination-host
// configuration.
//
// The configuration is consumed by the HTTP transport before the request is
// sent. If host is nil, WithHost panics.
func WithHost(ctx context.Context, host *Host) context.Context {
	if host == nil {
		panic("nil host")
	}
	return context.WithValue(ctx, hostContextKey{}, host)
}

func ParseTargetUrl(host string) (*url.URL, error) {
	if host == "" {
		return nil, nil
	}

	hostURL, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("invalid host address %q: %v", host, err)
	}
	return hostURL, nil
}
