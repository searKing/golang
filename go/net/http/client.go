// Copyright 2020 The searKing Author. All rights reserved.
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
)

type Client struct {
	http.Client
}

func (c *Client) Use(h ...RoundTripHandler) *Client {
	_, ok := c.Transport.(*Transport)
	if !ok {
		c.Transport = &Transport{Base: c.Transport}
	}

	// above guarantee its type is *Transport
	(c.Transport.(*Transport)).Use(h...)

	// for chained call
	return c
}

// parseURL is just url.Parse. It exists only so that url.Parse can be called
// in places where url is shadowed for godoc. See https://golang.org/cl/49930.
var parseURL = url.Parse

func NewClient(u string) (*http.Client, error) {
	url, err := parseURL(u)
	if err != nil {
		return nil, err
	}
	hostname := url.Hostname()
	if strings.Index(hostname, "unix:") != 0 {
		return http.DefaultClient, nil
	}
	tr := &http.Transport{
		DisableCompression: true,
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			return net.Dial("unix", url.Host)
		},
	}
	client := &http.Client{Transport: tr}
	return client, nil
}

func Head(url string) (resp *http.Response, err error) {
	client, err := NewClient(url)
	if err != nil {
		return nil, err
	}
	return client.Head(url)

}

func Get(url string) (resp *http.Response, err error) {
	client, err := NewClient(url)
	if err != nil {
		return nil, err
	}
	return client.Get(url)
}

func Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	client, err := NewClient(url)
	if err != nil {
		return nil, err
	}
	return client.Post(url, contentType, body)
}

func PostForm(url string, data url.Values) (resp *http.Response, err error) {
	client, err := NewClient(url)
	if err != nil {
		return nil, err
	}
	return client.PostForm(url, data)
}
