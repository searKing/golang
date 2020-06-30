// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WrapGinF is a helper function for wrapping gin.HandlerFunc
// Returns a http middleware
func WrapGinF(f gin.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		g := gin.New()
		g.Any("/*path", f)
		g.ServeHTTP(w, r)
	}
}
