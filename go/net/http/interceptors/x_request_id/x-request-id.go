package x_request_id

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

// DefaultXRequestIDKey is metadata key name for request ID
var DefaultXRequestIDKey = "X-Request-ID"

func appendInOutMetadata(ctx context.Context, r interface{}, requestIDs ...string) context.Context {
	switch rr := r.(type) {
	case *http.Request:
		for _, id := range requestIDs {
			rr.Header.Add(DefaultXRequestIDKey, id)
		}
	case http.ResponseWriter:
		for _, id := range requestIDs {
			rr.Header().Add(DefaultXRequestIDKey, id)
		}
	}
	return context.WithValue(ctx, DefaultXRequestIDKey, requestIDs)
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
func fromHTTPContext(r *http.Request) ([]string, bool) {
	key := DefaultXRequestIDKey
	if requestID := r.Header.Get(key); requestID != "" {
		return []string{requestID}, true
	}
	if requestID := r.URL.Query().Get(key); requestID != "" {
		return []string{requestID}, true
	}
	if requestID := r.FormValue(key); requestID != "" {
		return []string{requestID}, true
	}
	if requestID := r.PostFormValue(key); requestID != "" {
		return []string{requestID}, true
	}

	switch requestIDs := r.Context().Value(key).(type) {
	case string:
		return []string{requestIDs}, true
	case []string:
		return requestIDs, true
	default:
		return nil, false
	}
}

func RequestIDFromHTTPContext(ctx context.Context) []string {
	switch requestIDs := ctx.Value(DefaultXRequestIDKey).(type) {
	case string:
		return []string{requestIDs}
	case []string:
		return requestIDs
	default:
		return nil
	}
}
