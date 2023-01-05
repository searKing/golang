// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httphost_test

import (
	"log"
	"net/http"

	http_ "github.com/searKing/golang/go/net/http"
	"github.com/searKing/golang/go/net/http/httphost"
	_ "github.com/searKing/golang/go/net/resolver/passthrough"
)

func Example() {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	proxy := &httphost.Host{
		HostTarget: "127.0.0.1",
	}
	req = req.WithContext(httphost.WithHost(req.Context(), proxy))

	err := http_.HostFuncFromContext(req)
	if err != nil {
		log.Fatal(err)
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
