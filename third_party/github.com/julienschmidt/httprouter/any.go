// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httprouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Any(r *httprouter.Router, path string, handle httprouter.Handle) {
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
