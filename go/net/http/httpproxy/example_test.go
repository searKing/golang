// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpproxy_test

import (
	"log"
	"net/http"

	http_ "github.com/searKing/golang/go/net/http"
	"github.com/searKing/golang/go/net/http/httpproxy"
	_ "github.com/searKing/golang/go/net/resolver/passthrough"
)

func Example() {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	proxy := &httpproxy.Proxy{
		ProxyUrl:    "socks5://proxy.example.com:8080",
		ProxyTarget: "127.0.0.1",
	}
	req = req.WithContext(httpproxy.WithProxy(req.Context(), proxy))

	_, err := http_.DefaultTransportWithDynamicProxy.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}
}
