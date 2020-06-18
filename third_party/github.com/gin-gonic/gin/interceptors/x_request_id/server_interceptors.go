// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x_request_id

import "github.com/gin-gonic/gin"

// ServerInterceptor returns a new server interceptors with x-request-id in context.
func ServerInterceptor(keys ...interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newContextForHandleRequestID(ctx, keys...)
		ctx.Next() // execute all the handlers
	}
}
