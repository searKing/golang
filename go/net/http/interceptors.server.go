// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"

	"github.com/searKing/golang/go/net/http/internal"
)

// HandlerInterceptorChain interface that allows for customized handler execution chains.
// Applications can register any number of existing or custom interceptors for certain groups of handlers,
// to add common preprocessing behavior without needing to modify each handler implementation.
// A HandlerInterceptor gets called before the appropriate HandlerAdapter triggers the execution of the handler itself.
// This mechanism can be used for a large field of preprocessing aspects, e.g. for authorization checks, or common
// handler behavior like locale or theme changes.
// Its main purpose is to allow for factoring out repetitive handler code.
// https://docs.spring.io/spring-framework/docs/current/javadoc-api/org/springframework/web/servlet/HandlerInterceptor.html
//go:generate go-option -type "HandlerInterceptorChain"
type HandlerInterceptorChain struct {
	interceptors []internal.HandlerInterceptor
}

func NewHandlerInterceptorChain(opts ...HandlerInterceptorChainOption) *HandlerInterceptorChain {
	chain := &HandlerInterceptorChain{}
	chain.ApplyOptions(opts...)
	return chain
}

// InjectHttpHandler returns a http handler injected by chain
func (chain HandlerInterceptorChain) InjectHttpHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// short circuit
		if len(chain.interceptors) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// record where to reverse
		var it = -1
		defer func() {
			defer func() {
				err := recover()
				// no filter
				if it == -1 {
					panic(err)
				}
				for i := it; i >= 0; i-- {
					chain.interceptors[i].AfterCompletion(w, r, err)
				}
			}()
			if err := recover(); err != nil {
				panic(err)
			}
			for i := it; i >= 0; i-- {
				chain.interceptors[i].PostHandle(w, r)
			}
		}()

		for i, filter := range chain.interceptors {
			err := filter.PreHandle(w, r)
			if err != nil {
				// assumes that this interceptor has already dealt with the response itself
				return
			}
			// the execution chain should proceed with the next interceptor or the handler itself
			it = i
		}

		for i := range chain.interceptors {
			next = chain.interceptors[len(chain.interceptors)-1-i].WrapHandle(next)
		}

		next.ServeHTTP(w, r)
	})

}
