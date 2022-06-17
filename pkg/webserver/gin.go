// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/searKing/golang/pkg/webserver/healthz"
)

type ginMuxer struct {
	muxer gin.IRouter
}

func GinMuxer(muxer gin.IRouter) healthz.Muxer {
	return &ginMuxer{muxer: muxer}
}
func (mux *ginMuxer) Handle(pattern string, handler http.Handler) {
	mux.muxer.GET(pattern, gin.WrapH(handler))
}
