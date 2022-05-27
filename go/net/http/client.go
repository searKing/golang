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
	"time"
)

type Client struct {
	http.Client
	target               string // resolver.Target, will replace Host in url.Url
	replaceHostInRequest bool   // resolve target to proxy and replace host if target resolved
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

// NewClient returns a http client wrapper behaves like http.Client
// u is the original url to send HTTP request
// target is the resolver to resolve Host to send HTTP request,
// that is replacing host in url(NOT HOST in http header) by address resolved by target
// proxyUrl is proxy's url, like sock5://127.0.0.1:8080
// proxyTarget is proxy's addr, replace the HOST in proxyUrl if not empty
func NewClient(u, target string, proxyUrl string, proxyTarget string) (*Client, error) {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
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
	if proxyUrl != "" && proxyTarget != "" {
		tr.Proxy = ProxyFuncWithTargetOrDefault(proxyUrl, proxyTarget, tr.Proxy)
	}
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
	return &Client{
		Client:               client,
		target:               target,
		replaceHostInRequest: false,
	}, nil
}

// NewClientWithTarget returns a Client with http.Client and host replaced by resolver.Target
// target is the resolver to resolve Host to send HTTP request,
// that is replacing host in url(NOT HOST in http header) by address resolved by target
func NewClientWithTarget(target string) *Client {
	cli, _ := NewClient("", target, "", "")
	return cli
}

// NewClientWithProxy returns a Client with http.Client with proxy set by resolver.Target
// proxyUrl is proxy's url, like sock5://127.0.0.1:8080
// proxyTarget is proxy's addr, replace the HOST in proxyUrl if not empty
func NewClientWithProxy(proxyUrl string, proxyTarget string) *Client {
	cli, _ := NewClient("", "", proxyUrl, proxyTarget)
	return cli
}

func NewClientWithUnixDisableCompression(u string) (*Client, error) {
	return NewClient(u, "", "", "")
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	err := RequestWithTarget(req, c.target, c.replaceHostInRequest)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
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
