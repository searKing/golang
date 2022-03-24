// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type serverHandler struct {
	gateway *Gateway

	once sync.Once
}

func (s *serverHandler) refreshHandler(httpHandler http.Handler) {
	oldOtherErrorHandler := runtime.OtherErrorHandler

	runtime.OtherErrorHandler = func(w http.ResponseWriter, r *http.Request, msg string, code int) {
		if code == http.StatusNotFound || code == http.StatusMethodNotAllowed {
			httpHandler.ServeHTTP(w, r)
			return
		}
		oldOtherErrorHandler(w, r, msg, code)
		return
	}
}

func (s *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	oldOtherErrorHandler := runtime.OtherErrorHandler
	defer func() { runtime.OtherErrorHandler = oldOtherErrorHandler }()

	httpHandler := s.gateway.Handler
	if httpHandler == nil {
		httpHandler = http.DefaultServeMux
	}
	if s.gateway.opt.fastMode {
		s.once.Do(func() {
			s.refreshHandler(httpHandler)
		})
	} else {
		s.refreshHandler(httpHandler)
	}
	s.gateway.opt.interceptors.InjectHttpHandler(s.gateway.httpMuxToGrpc).ServeHTTP(w, r)
}
