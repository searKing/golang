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

type Client struct {
	http.Client

	Proxy *httpproxy.Proxy
	Host  *httphost.Host
}

// Use adds middleware handlers to the transport.
func (c *Client) Use(d ...RoundTripDecorator) *Client {
	if len(d) == 0 {
		return c
	}
	var rts RoundTripDecorators
	rts = append(rts, d...)
	c.Transport = rts.WrapRoundTrip(c.Transport)
	// for chained call
	return c
}

// parseURL is just url.Parse. It exists only so that url.Parse can be called
// in places where url is shadowed for godoc. See https://golang.org/cl/49930.
var parseURL = url.Parse

// NewClient returns a http client wrapper behaves like http.Client,
// sends HTTP Request to target by proxy url with Host replaced by proxyTarget
//
// u is the original url to send HTTP request, empty usually.
// target is the resolver to resolve Host to send HTTP request,
// that is replacing host in url(NOT HOST in http header) by address resolved by Host
// fixedProxyUrl is proxy's url, like socks5://127.0.0.1:8080
// fixedProxyTarget is as like gRPC Naming for proxy service discovery, with Host in TargetUrl replaced if not empty.
func NewClient(u, hostTarget string, proxyUrl string, proxyTarget string) (*Client, error) {
	tr := DefaultTransportWithDynamicHostAndProxy
	if len(u) > 0 {
		urlParsed, err := parseURL(u)
		if err != nil {
			return nil, err
		}
		hostname := urlParsed.Hostname()
		if strings.Index(hostname, "unix:") == 0 {
			tr = &http.Transport{
				DisableCompression: true,
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
					return net.Dial("unix", urlParsed.Host)
				},
			}
		}
	}
	client := http.Client{Transport: tr}
	c := &Client{
		Client: client,
	}
	if hostTarget != "" {
		c.Host = &httphost.Host{
			HostTarget:           hostTarget,
			ReplaceHostInRequest: false,
		}
	}
	if proxyUrl != "" {
		c.Proxy = &httpproxy.Proxy{
			ProxyUrl:    proxyUrl,
			ProxyTarget: proxyTarget,
		}
	}
	return c, nil
}

// NewClientWithTarget returns a Client with http.Client and host replaced by resolver.Host
// target is the resolver to resolve Host to send HTTP request,
// that is replacing host in url(NOT HOST in http header) by address resolved by Host
func NewClientWithTarget(target string) *Client {
	cli, _ := NewClient("", target, "", "")
	return cli
}

// NewClientWithProxy returns a Client with http.Client with proxy set by resolver.Host
// TargetUrl is proxy's url, like socks5://127.0.0.1:8080
// Host is proxy's addr, replace the HOST in TargetUrl if not empty
func NewClientWithProxy(proxyUrl string, proxyTarget string) *Client {
	cli, _ := NewClient("", "", proxyUrl, proxyTarget)
	return cli
}

func NewClientWithUnixDisableCompression(u string) (*Client, error) {
	return NewClient(u, "", "", "")
}

func (c *Client) Do(req *http.Request) (_ *http.Response, err error) {
	return c.Client.Do(RequestWithHostTarget(RequestWithProxyTarget(req, c.Proxy), c.Host))
}

func (c *Client) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
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
