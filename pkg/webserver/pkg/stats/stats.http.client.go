// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stats

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptrace"

	"github.com/searKing/golang/pkg/webserver/pkg/logging"
)

func NewHttpClientTrace(ctx context.Context) *httptrace.ClientTrace {
	logger := slog.With(logging.Attrs[any](ctx)...)
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			logger.Info(fmt.Sprintf("HTTP Client Get Conn to %s", hostPort))
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			logger.Info(fmt.Sprintf("HTTP Client Got Conn from %s to %s", connInfo.Conn.LocalAddr(), connInfo.Conn.RemoteAddr()))
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			logger.Info(fmt.Sprintf("HTTP Client Dns Done: %+v", dnsInfo))
		},
	}
}

func NewHttpRequestWithTrace(req *http.Request) *http.Request {
	return req.WithContext(httptrace.WithClientTrace(req.Context(), NewHttpClientTrace(req.Context())))
}
