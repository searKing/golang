// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package negroni

import (
	"github.com/urfave/negroni"
	"net/http"
)

func NopHandlerFunc(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next.ServeHTTP(rw, r)
	return
}

var NopHandler = negroni.HandlerFunc(NopHandlerFunc)
