// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	slices_ "github.com/searKing/golang/go/exp/slices"
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
			attrs := fieldsToAttrSlice(fields)
			if len(attrs) > 0 {
				extra := strings.Join(slices_.MapFunc(attrs, func(e slog.Attr) string { return e.String() }), ", ")
				return fmt.Sprintf(" | %v", extra)
			}
		}
		return ""
	})
}

func fieldsToAttrSlice(fields logging.Fields) []slog.Attr {
	var attrs []slog.Attr
	i := fields.Iterator()
	for i.Next() {
		k, v := i.At()
		attrs = append(attrs, slog.Any(k, v))
	}
	return attrs
}

const badKey = "!BADKEY"

// argsToAttr turns a prefix of the nonempty args slice into an Attr
// and returns the unconsumed portion of the slice.
// If args[0] is an Attr, it returns it.
// If args[0] is a string, it treats the first two elements as
// a key-value pair.
// Otherwise, it treats args[0] as a value with a missing key.
func argsToAttr(args []any) (slog.Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return slog.String(badKey, x), nil
		}
		return slog.Any(x, args[1]), args[2:]

	case slog.Attr:
		return x, args[1:]

	default:
		return slog.Any(badKey, x), args[1:]
	}
}
