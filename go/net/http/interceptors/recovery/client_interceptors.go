// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recovery

import (
	"io"
	"net/http"

	http_ "github.com/searKing/golang/go/net/http"
)

// ClientInterceptor returns a new client interceptors with recovery from panic.
// affect as recover{f()}; next()
func ClientInterceptor(next http_.RoundTripHandler, out io.Writer, f func(resp *http.Response, req *http.Request, err interface{})) http_.RoundTripHandler {
	return http_.RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		defer func() {
			Recover(out, req, func(err interface{}) interface{} {
				if f == nil {
					return nil
				}
				f(resp, req, err)
				return nil
			})
		}()
		resp, err = next.RoundTrip(req)
		return
	})
}
