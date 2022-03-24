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

// WrapHTTPRouterF is a helper function for wrapping httprouter.Handle
// Returns a http middleware
func WrapHTTPRouterF(f httprouter.Handle) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := httprouter.New()
		any(h, "/*path", f)
		h.ServeHTTP(w, r)
	})
}

func any(r *httprouter.Router, path string, handle httprouter.Handle) {
	r.Handle(http.MethodGet, path, handle)
	r.Handle(http.MethodHead, path, handle)
	r.Handle(http.MethodPost, path, handle)
	r.Handle(http.MethodPut, path, handle)
	r.Handle(http.MethodPatch, path, handle)
	r.Handle(http.MethodDelete, path, handle)
	r.Handle(http.MethodConnect, path, handle)
	r.Handle(http.MethodOptions, path, handle)
	r.Handle(http.MethodTrace, path, handle)
}
