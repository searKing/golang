// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gin

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func RedirectTrim(code int, prefix string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u := *ctx.Request.URL
		u.Path = strings.TrimPrefix(u.Path, prefix)
		ctx.Redirect(code, u.String())
	}
}

func Redirect(code int, path string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u := *ctx.Request.URL
		u.Path = path
		ctx.Redirect(code, u.String())
	}
}
