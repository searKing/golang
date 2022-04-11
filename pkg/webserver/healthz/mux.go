// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package healthz

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginMuxer struct {
	muxer gin.IRouter
}

func GinMuxer(muxer gin.IRouter) mux {
	return &ginMuxer{muxer: muxer}
}
func (mux *ginMuxer) Handle(pattern string, handler http.Handler) {
	mux.muxer.GET(pattern, gin.WrapH(handler))
}
