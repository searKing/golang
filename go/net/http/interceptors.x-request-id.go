// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"context"
	"net/http"
)

// DefaultXRequestIDKey is metadata key name for request ID
var DefaultXRequestIDKey = "X-Request-ID"

// setInOutMetadata injects requestIDs in req|resp's Header and context
// 将request-id追加注入请求|响应头及context中
func setInOutMetadata(ctx context.Context, w http.ResponseWriter, r *http.Request, requestID string) context.Context {
	if r != nil {
		r.Header.Set(DefaultXRequestIDKey, requestID)
	}
	if w != nil {
		w.Header().Set(DefaultXRequestIDKey, requestID)
	}
	return context.WithValue(ctx, DefaultXRequestIDKey, requestID)
}

// parse request id from gin.Context
// query | header | post form | context
// 从请求中提取request-id
func fromHTTPContext(r *http.Request, keys ...interface{}) string {
	key := DefaultXRequestIDKey
	if requestID := r.Header.Get(key); requestID != "" {
		return requestID
	}
	if requestID := r.URL.Query().Get(key); requestID != "" {
		return requestID
	}
	if requestID := r.FormValue(key); requestID != "" {
		return requestID
	}
	if requestID := r.PostFormValue(key); requestID != "" {
		return requestID
	}

	switch requestIDs := r.Context().Value(key).(type) {
	case string:
		if requestIDs != "" {
			return requestIDs
		}
	case []string:
		if len(requestIDs) > 0 {
			return requestIDs[0]
		}
	}

	return fromContext(r.Context(), keys...)
}

// fromContext takes out first value from context by keys
func fromContext(ctx context.Context, keys ...interface{}) string {
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

func RequestIDFromContext(ctx context.Context) string {
	switch requestIDs := ctx.Value(DefaultXRequestIDKey).(type) {
	case string:
		return requestIDs
	case []string:
		if len(requestIDs) > 0 {
			return requestIDs[0]
		}
	}
	return ""
}
