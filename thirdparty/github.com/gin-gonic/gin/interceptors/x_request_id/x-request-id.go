// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x_request_id

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/searKing/golang/go/net/http/interceptors/x_request_id"
)

// key is RequestID within Context if have
func newContextForHandleRequestID(ctx *gin.Context, keys ...interface{}) {
	requestIDs, ok := fromGinContext(ctx)
	if !ok || len(requestIDs) == 0 {
		appendInOutMetadata(ctx, newRequestID(ctx, keys...)...)
		return
	}
	appendInOutMetadata(ctx, requestIDs...)
	return
}

// to chain multiple request ids by generating new request id for each request and concatenating it to original request ids.
// key is RequestID within Context if have
func newContextForHandleRequestIDChain(ctx *gin.Context, keys ...interface{}) {
	requestIDs, ok := fromGinContext(ctx)
	if !ok || len(requestIDs) == 0 {
		appendInOutMetadata(ctx, newRequestIDChain(ctx, keys...)...)
		return
	}
	requestIDs = append(requestIDs, uuid.New().String())
	appendInOutMetadata(ctx, requestIDs...)
	return
}

func appendInOutMetadata(ctx *gin.Context, requestIDs ...string) {
	ctx.Set(x_request_id.DefaultXRequestIDKey, requestIDs)
	for _, id := range requestIDs {
		ctx.Writer.Header().Add(x_request_id.DefaultXRequestIDKey, id)
	}
}

func newRequestID(ctx context.Context, keys ...interface{}) []string {
	for _, key := range keys {
		val := ctx.Value(key)
		switch val := val.(type) {
		case string:
			return []string{val}
		case []string:
			return val
		}
	}
	return []string{uuid.New().String()}
}

func newRequestIDChain(ctx context.Context, keys ...interface{}) []string {
	for _, key := range keys {
		val := ctx.Value(key)
		switch val := val.(type) {
		case string:
			return []string{val}
		case []string:
			return append(val, uuid.New().String())
		}
	}
	return []string{uuid.New().String()}
}

// parse request id from gin.Context
// query | header | post form | context
func fromGinContext(ctx *gin.Context) ([]string, bool) {
	key := x_request_id.DefaultXRequestIDKey
	if requestID := ctx.GetHeader(key); requestID != "" {
		return []string{requestID}, true
	}
	if requestIDs, ok := ctx.GetQueryArray(key); ok {
		return requestIDs, true
	}
	if requestIDs, ok := ctx.GetPostFormArray(key); ok {
		return requestIDs, true
	}

	requestIDs, has := ctx.Get(key)
	if !has {
		return nil, false
	}
	switch requestIDs := requestIDs.(type) {
	case string:
		return []string{requestIDs}, true
	case []string:
		return requestIDs, true
	default:
		return nil, false
	}
}

func RequestIDFromGinContext(ctx *gin.Context) []string {
	requestIDs, has := ctx.Get(x_request_id.DefaultXRequestIDKey)
	if !has {
		return nil
	}
	switch requestIDs := requestIDs.(type) {
	case string:
		return []string{requestIDs}
	case []string:
		return requestIDs
	default:
		return nil
	}

}
