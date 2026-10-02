// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/searKing/golang/go/net/http/httphost"
	"github.com/searKing/golang/go/net/http/httpproxy"
)

// parseURL is just url.Parse. It exists only so that url.Parse can be called
// in places where url is shadowed for godoc. See https://golang.org/cl/49930.
var parseURL = url.Parse

// NewClient returns an http.Client that sends requests to hostTarget
// through proxyUrl.
//
// u is the original url to send HTTP request, empty usually.
// hostTarget is the resolver target for the destination host. It replaces
// the host in the request URL, not the HTTP Host header.
// proxyUrl is the proxy URL, like socks5://127.0.0.1:8080.
// proxyTarget is a resolver target for the proxy, like gRPC naming. When
// non-empty, it replaces the host in proxyUrl.
//
// Host and proxy are attached by round-trip decorators, so every method on
// the returned http.Client, including Get and Do, applies them.
func NewClient(u, hostTarget string, proxyUrl string, proxyTarget string) (*http.Client, error) {
	base := DefaultTransportWithDynamicHostAndProxy
	if len(u) > 0 {
		urlParsed, err := parseURL(u)
		if err != nil {
			return nil, err
		}
		hostname := urlParsed.Hostname()
		if strings.Index(hostname, "unix:") == 0 {
			base = &http.Transport{
				DisableCompression: true,
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
					return net.Dial("unix", urlParsed.Host)
				},
			}
		}
	}

	var rts RoundTripDecorators
	if hostTarget != "" {
		rts = append(rts, roundTripDecoratorWithHostTarget(hostTarget))
	}
	if proxyUrl != "" {
		rts = append(rts, roundTripDecoratorWithProxyTarget(proxyUrl, proxyTarget))
	}
	return &http.Client{Transport: rts.WrapRoundTrip(base)}, nil
}

// NewClientWithTarget returns an http.Client whose destination host is
// resolved from target. target replaces the host in the request URL, not
// the HTTP Host header.
func NewClientWithTarget(target string) *http.Client {
	cli, _ := NewClient("", target, "", "")
	return cli
}

// NewClientWithProxy returns an http.Client that sends through proxyUrl.
// proxyUrl is the proxy URL, like socks5://127.0.0.1:8080.
// proxyTarget replaces the host in proxyUrl when non-empty.
func NewClientWithProxy(proxyUrl string, proxyTarget string) *http.Client {
	cli, _ := NewClient("", "", proxyUrl, proxyTarget)
	return cli
}

func NewClientWithUnixDisableCompression(u string) (*http.Client, error) {
	return NewClient(u, "", "", "")
}

func roundTripDecoratorWithHostTarget(target string) RoundTripDecorator {
	return RoundTripDecoratorFunc(func(rt http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			host := &httphost.Host{HostTarget: target}
			return rt.RoundTrip(RequestWithHostTarget(req, host))
		})
	})
}

func roundTripDecoratorWithProxyTarget(proxyUrl string, proxyTarget string) RoundTripDecorator {
	return RoundTripDecoratorFunc(func(rt http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			proxy := &httpproxy.Proxy{
				ProxyUrl:    proxyUrl,
				ProxyTarget: proxyTarget,
			}
			return rt.RoundTrip(RequestWithProxyTarget(req, proxy))
		})
	})
}

func Head(url string) (resp *http.Response, err error) {
	client, err := NewClientWithUnixDisableCompression(url)
	if err != nil {
		return nil, err
	}
	return client.Head(url)
}

func Get(url string) (resp *http.Response, err error) {
	client, err := NewClientWithUnixDisableCompression(url)
	if err != nil {
		return nil, err
	}
	return client.Get(url)
}

func Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	client, err := NewClientWithUnixDisableCompression(url)
	if err != nil {
		return nil, err
	}
	return client.Post(url, contentType, body)
}

func PostForm(url string, data url.Values) (resp *http.Response, err error) {
	client, err := NewClientWithUnixDisableCompression(url)
	if err != nil {
		return nil, err
	}
	return client.PostForm(url, data)
}
