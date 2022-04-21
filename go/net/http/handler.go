// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import "net/http"

// HandlerDecorator is an interface representing the ability to decorate or wrap a http.Handler.
type HandlerDecorator interface {
	WrapHandler(rt http.Handler) http.Handler
}

// HandlerDecorators defines a HandlerDecorator slice.
type HandlerDecorators []HandlerDecorator

func (chain HandlerDecorators) WrapHandler(next http.Handler) http.Handler {
	for i := range chain {
		h := chain[len(chain)-1-i]
		if h != nil {
			next = h.WrapHandler(next)
		}
	}
	return next
}
