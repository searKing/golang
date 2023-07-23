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

// RequestWithHostTarget replace Host in url.Url by resolver.Host
// replace Host in req if replaceHostInRequest is true
func RequestWithHostTarget(req *http.Request, target *httphost.Host) *http.Request {
	if target == nil {
		return req
	}
	return req.WithContext(httphost.WithHost(req.Context(), target))
}

// HostFuncFromContext builds a host function from the given string, which should
// represent a Target that can be used as a host. It performs basic
// sanitization of the Target retrieved in context of Request, and returns any error encountered.
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

// DefaultTransportWithDynamicHost is the default implementation of Transport and is
// used by DefaultClientWithDynamicHost. It establishes network connections as needed
// and caches them for reuse by subsequent calls.
var DefaultTransportWithDynamicHost = RoundTripperWithTarget(http.DefaultTransport)

// DefaultClientWithDynamicHost is the default Client with DefaultTransportWithDynamicHost.
var DefaultClientWithDynamicHost = &http.Client{
	Transport: DefaultTransportWithDynamicHost,
}

// RoundTripperWithTarget wraps http.RoundTripper with request url replaced by Target resolved by resolver.
// Target is as like gRPC Naming for service discovery.
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

// DefaultTransportWithDynamicHostAndProxy is the default implementation of Transport and is
// used by DefaultClientWithDynamicHostAndProxy. It establishes network connections as needed
// and caches them for reuse by subsequent calls.
var DefaultTransportWithDynamicHostAndProxy = RoundTripperWithTarget(DefaultTransportWithDynamicProxy)

// DefaultClientWithDynamicHostAndProxy is the default Client with DefaultTransportWithDynamicHostAndProxy.
var DefaultClientWithDynamicHostAndProxy = &http.Client{
	Transport: DefaultTransportWithDynamicHostAndProxy,
}
