// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httprouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// WrapF is a helper function for wrapping http.HandlerFunc
// Returns a httprouter middleware
func WrapF(f http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		f.ServeHTTP(w, r)
	}
}

// WrapH is a helper function for wrapping http.Handler
// Returns a httprouter middleware
func WrapH(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}
