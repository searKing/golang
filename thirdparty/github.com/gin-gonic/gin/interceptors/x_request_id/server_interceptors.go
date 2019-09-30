package x_request_id

import "github.com/gin-gonic/gin"

// ServerInterceptor returns a new server interceptors with x-request-id in context.
func ServerInterceptor(keys ...interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newContextForHandleRequestID(ctx, keys...)
		ctx.Next() // execute all the handlers
	}
}

// ServerInterceptor returns a new server interceptors with x-request-id in context.
func ServerInterceptorChain(keys ...interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newContextForHandleRequestID(ctx, keys...)
		ctx.Next() // execute all the handlers
	}
}
