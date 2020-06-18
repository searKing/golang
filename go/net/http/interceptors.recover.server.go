// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"io"
	"net/http"
)

// RecoveryServerInterceptor returns a new server interceptors with recovery from panic.
// affect as recover{f()}; next()
func RecoveryServerInterceptor(next http.Handler, out io.Writer, f func(w http.ResponseWriter, r *http.Request, err interface{})) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			Recover(out, r, func(err interface{}) interface{} {
				if f == nil {
					return nil
				}
				f(w, r, err)
				return nil
			})
		}()
		next.ServeHTTP(w, r)
	})
}
