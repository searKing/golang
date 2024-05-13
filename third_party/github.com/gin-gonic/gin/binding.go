// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	binding_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin/binding"
)

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
func DefaultBinding(ctx *gin.Context) binding.Binding {
	if ctx.Request.Method == http.MethodGet {
		return binding.Form
	}
	switch ctx.ContentType() {
	case binding.MIMEJSON:
		return binding_.JSONPB // support proto3 if enabled
	//return binding.JSON
	default:
		return binding.Default(ctx.Request.Method, ctx.ContentType())
	}
}
func Bind(ctx *gin.Context, obj any) error {
	log.Println(`BindWith(\"interface{}, binding.Binding\") error is going to
	be deprecated, please check issue #662 and either use MustBindWith() if you
	want HTTP 400 to be automatically returned if any error occur, or use
	ShouldBindWith() if you need to manage the error.`)
	return ctx.MustBindWith(obj, DefaultBinding(ctx))
}

// ShouldBind binds the passed struct pointer using the specified binding engine.
// See the binding package.
func ShouldBind(ctx *gin.Context, obj any) error {
	return ctx.ShouldBindWith(obj, DefaultBinding(ctx))
}

// MustBindWith binds the passed struct pointer using the specified binding engine.
// It will abort the request with HTTP 400 if any error occurs.
// See the binding package.
func MustBind(ctx *gin.Context, obj any) error {
	return ctx.MustBindWith(obj, DefaultBinding(ctx))
}
