// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"fmt"
	"net/http"

	"github.com/searKing/golang/go/net/http/httphost"
	"github.com/searKing/golang/go/net/http/httpproxy"
	"github.com/searKing/golang/go/net/resolver"
	time_ "github.com/searKing/golang/go/time"
)

var _ RoundTripDecorator = RoundTripDecoratorFunc(RoundTripperWithTarget)

// RequestWithHostTarget attaches a dynamic destination host to the request
// context.
//
// The target is resolved later by HostFuncFromContext before the request is
// sent. This allows dynamic service discovery and destination rewriting.
func RequestWithHostTarget(req *http.Request, target *httphost.Host) *http.Request {
	if target == nil {
		return req
	}
	return req.WithContext(httphost.WithHost(req.Context(), target))
}

// HostFuncFromContext resolves the dynamic destination host configured in the
// request context and updates the request before it is sent to the underlying
// transport.
//
// It resolves HostTarget through the configured resolver. The resolved address
// replaces req.URL.Host, which determines the network address to which the
// transport connects.
//
// If ReplaceHostInRequest is true, req.Host is also set to the resolved
// address. This affects the HTTP Host header sent to the server.
//
// Example:
//
//	Original request URL: http://api.example.com/v1
//	HostTarget:          "backend-service"
//	Resolver output:     "10.0.0.8:8080"
//
//	After resolution:
//	  req.URL.Host = "10.0.0.8:8080"
//	  req.Host = "api.example.com" (unchanged if ReplaceHostInRequest is false)
//
// This is not an HTTP forward proxy. It is for service discovery and
// destination rewriting only.
func HostFuncFromContext(req *http.Request) error {
	host := httphost.ContextHost(req.Context())
	// load host from environment if host not set
	if host == nil || host.HostTarget == "" {
		return nil
	}
	if req.URL == nil {
		return nil
	}

	if host.HostTarget == "" {
		return nil
	}

	// replace host of host if target of host if resolved
	address, err := resolver.ResolveOneAddr(req.Context(), host.HostTarget)
	if err != nil {
		return err
	}
	if address.Addr != "" {
		req.URL.Host = address.Addr
	}
	host.HostTargetAddrResolved = address
	if host.ReplaceHostInRequest {
		req.Host = req.URL.Host
	}
	return nil
}

// DefaultTransportWithDynamicHost is an http.RoundTripper that supports
// dynamic destination host resolution.
//
// It wraps http.DefaultTransport to add service discovery and host rewriting
// capability through HostFuncFromContext. This is useful for scenarios where
// the destination address is resolved at request time based on a logical
// service name.
var DefaultTransportWithDynamicHost = RoundTripperWithTarget(http.DefaultTransport)

// DefaultClientWithDynamicHost is an http.Client that uses
// DefaultTransportWithDynamicHost.
var DefaultClientWithDynamicHost = &http.Client{
	Transport: DefaultTransportWithDynamicHost,
}

// RoundTripperWithTarget wraps an http.RoundTripper to support dynamic
// destination host resolution and optional HTTP Host header rewriting.
//
// For each request, it:
//  1. Resolves HostTarget through the configured resolver
//  2. Updates req.URL.Host with the resolved address
//  3. Optionally updates req.Host based on ReplaceHostInRequest
//  4. Reports resolver results and request metrics through resolver.ResolveDone
//  5. Wraps errors with resolver and destination information
//
// The underlying RoundTripper (typically http.DefaultTransport) connects to
// the resolved destination address.
//
// This wrapper also handles proxy resolution if httpproxy.Proxy is configured
// in the request context, and reports both host and proxy metrics.
func RoundTripperWithTarget(rt http.RoundTripper) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		err = HostFuncFromContext(req)
		if err != nil {
			return nil, err
		}
		var cost time_.Cost
		cost.Start()
		defer func() {
			if host := httphost.ContextHost(req.Context()); host != nil {
				_ = resolver.ResolveDone(req.Context(), host.HostTarget, resolver.DoneInfo{
					Err:      err,
					Addr:     host.HostTargetAddrResolved,
					Duration: cost.Elapse(),
				})
				if err != nil && host.HostTargetAddrResolved.Addr != "" {
					var s string
					if host.HostTarget != "" {
						s = fmt.Sprintf(" in target(%s)", host.HostTarget)
					}
					err = fmt.Errorf("->http_host(%s)%s: %w", host.HostTargetAddrResolved.Addr, s, err)
				}
			}

			if proxy := httpproxy.ContextProxy(req.Context()); proxy != nil {
				_ = resolver.ResolveDone(req.Context(), proxy.ProxyTarget, resolver.DoneInfo{
					Err:      err,
					Addr:     proxy.ProxyAddrResolved,
					Duration: cost.Elapse(),
				})
				if err != nil && proxy.ProxyAddrResolved.Addr != "" {
					var s string
					if proxy.ProxyTarget != "" {
						s = fmt.Sprintf(" in target(%s)", proxy.ProxyTarget)
					}
					err = fmt.Errorf("->http_proxy(%s)%s: %w", proxy.ProxyAddrResolved.Addr, s, err)
				}
			}
		}()
		return rt.RoundTrip(req)
	})
}

// DefaultTransportWithDynamicHostAndProxy is an http.RoundTripper that
// supports both dynamic destination host resolution and dynamic proxy resolution.
//
// It combines the functionality of DefaultTransportWithDynamicHost and
// DefaultTransportWithDynamicProxy, allowing independent control over the
// final business destination and the proxy server used to reach it.
var DefaultTransportWithDynamicHostAndProxy = RoundTripperWithTarget(DefaultTransportWithDynamicProxy)

// DefaultClientWithDynamicHostAndProxy is an http.Client that uses DefaultTransportWithDynamicHostAndProxy.
var DefaultClientWithDynamicHostAndProxy = &http.Client{
	Transport: DefaultTransportWithDynamicHostAndProxy,
}
