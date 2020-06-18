// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x_request_id

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/searKing/golang/go/net/http/interceptors/x_request_id"
)

// key is RequestID within Context if have
func newContextForHandleRequestID(ctx *gin.Context, keys ...interface{}) {
	requestID := fromGinContext(ctx, keys)
	if requestID == "" {
		requestID = uuid.New().String()
	}
	setInOutMetadata(ctx, requestID)
}

func setInOutMetadata(ctx *gin.Context, requestIDs ...string) {
	for _, id := range requestIDs {
		ctx.Request.Header.Set(x_request_id.DefaultXRequestIDKey, id)
		ctx.Writer.Header().Set(x_request_id.DefaultXRequestIDKey, id)
	}
	ctx.Set(x_request_id.DefaultXRequestIDKey, requestIDs)
}

// parse request id from gin.Context
// query | header | post form | context
func fromGinContext(ctx *gin.Context, keys ...interface{}) string {
	key := x_request_id.DefaultXRequestIDKey
	if requestID := ctx.GetHeader(key); requestID != "" {
		return requestID
	}
	if requestIDs, ok := ctx.GetQueryArray(key); ok {
		return requestIDs[0]
	}
	if requestIDs, ok := ctx.GetPostFormArray(key); ok {
		return requestIDs[0]
	}

	requestIDs, has := ctx.Get(key)
	if !has {
		return ""
	}
	switch requestIDs := requestIDs.(type) {
	case string:
		if requestIDs != "" {
			return requestIDs
		}
	case []string:
		if len(requestIDs) > 0 {
			return requestIDs[0]
		}
	}

	for _, key := range keys {
		val := ctx.Value(key)
		switch val := val.(type) {
		case string:
			if val != "" {
				return val
			}
		case []string:
			if len(val) > 0 {
				if val[0] != "" {
					return val[0]
				}
			}
		}
	}
	return ""
}

func RequestIDFromGinContext(ctx *gin.Context) string {
	requestIDs, has := ctx.Get(x_request_id.DefaultXRequestIDKey)
	if !has {
		return ""
	}
	switch requestIDs := requestIDs.(type) {
	case string:
		return requestIDs
	case []string:
		if len(requestIDs) > 0 {
			return requestIDs[0]
		}
	}
	return ""
}
