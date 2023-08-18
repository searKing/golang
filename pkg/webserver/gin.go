// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/searKing/golang/pkg/webserver/healthz"
	gin_ "github.com/searKing/golang/third_party/github.com/gin-gonic/gin"
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

// GinLogFormatter is the log format function [gin.Logger] middleware uses.
func GinLogFormatter(layout string) func(param gin.LogFormatterParams) string {
	return gin_.LogFormatterWithExtra(layout, func(param gin.LogFormatterParams) string {
		if param.Request != nil {
			fields := logging.ExtractFields(param.Request.Context())
			if len(fields) > 0 {
				return fmt.Sprintf(" | %v", fields)
			}
		}
		return ""
	})
}
