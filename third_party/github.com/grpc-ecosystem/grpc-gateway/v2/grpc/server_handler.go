// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"net/http"
)

type serverHandler struct {
	gateway *Gateway
}

func (s *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.gateway.opt.interceptors.InjectHttpHandler(s.gateway.httpMuxToGrpc).ServeHTTP(w, r)
}
