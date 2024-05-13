// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"io"
	"net/http"
)

// RecoveryClientInterceptor returns a new client interceptors with recovery from panic.
// affect as recover{f()}; next()
func RecoveryClientInterceptor(next http.RoundTripper, out io.Writer, f func(resp *http.Response, req *http.Request, err any)) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		defer func() {
			Recover(out, req, func(err any) any {
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
