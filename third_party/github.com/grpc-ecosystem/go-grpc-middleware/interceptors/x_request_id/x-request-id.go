// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x_request_id

import (
	"context"
	"strings"

	"github.com/google/uuid"
	metadata_ "github.com/searKing/golang/third_party/google.golang.org/grpc/metadata"
	"google.golang.org/grpc/metadata"
)

// DefaultXRequestIDKey is metadata key name for request ID
var DefaultXRequestIDKey = "X-Request-ID"

// ExtractOrCreateRequestId extracts or create request_id from context by mdIncomingKey and key
func ExtractOrCreateRequestId(ctx context.Context, key interface{}) []string {
	requestIDs := ExtractRequestId(ctx, key)
	if len(requestIDs) == 0 {
		return []string{newRequestID(ctx, key)}
	}
	return requestIDs
}

// ExtractRequestId extracts request_id from context by mdIncomingKey and key
func ExtractRequestId(ctx context.Context, key interface{}) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return md.Get(DefaultXRequestIDKey)
	}

	requestIDs := md.Get(DefaultXRequestIDKey)
	if len(requestIDs) == 0 {
		return []string{newRequestID(ctx, key)}
	}
	return requestIDs
}

// InjectOrCreateRequestId injects or create request_id to context by mdIncomingKey and key
func InjectOrCreateRequestId(ctx context.Context, key interface{}, requestId string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		if requestId == "" {
			requestId = newRequestID(ctx, key)
		}
		return appendInOutMetadata(ctx, metadata.New(map[string]string{DefaultXRequestIDKey: requestId}))
	}
	requestIDs := md.Get(DefaultXRequestIDKey)
	if len(requestIDs) == 0 {
		if requestId == "" {
			requestId = newRequestID(ctx, key)
		}
		return appendInOutMetadata(ctx, metadata.New(map[string]string{DefaultXRequestIDKey: requestId}))
	}
	if requestId != "" {
		requestIDs[len(requestIDs)-1] = requestId
		md.Set(DefaultXRequestIDKey, requestIDs...)
	}

	return ctx
}

// InjectRequestId injects request_id to context by mdIncomingKey
func InjectRequestId(ctx context.Context, requestId string) context.Context {
	if requestId == "" {
		return ctx
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return appendInOutMetadata(ctx, metadata.New(map[string]string{DefaultXRequestIDKey: requestId}))
	}
	requestIDs := md.Get(DefaultXRequestIDKey)
	if len(requestIDs) == 0 {
		return appendInOutMetadata(ctx, metadata.New(map[string]string{DefaultXRequestIDKey: requestId}))
	}

	requestIDs[len(requestIDs)-1] = requestId
	md.Set(DefaultXRequestIDKey, requestIDs...)
	return ctx
}

// key is RequestID within Context if have
func newContextForHandleRequestID(ctx context.Context, key interface{}) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return appendInOutMetadata(ctx, metadata.New(map[string]string{DefaultXRequestIDKey: newRequestID(ctx, key)}))
	}
	requestIDs := md.Get(DefaultXRequestIDKey)
	if len(requestIDs) == 0 {
		return appendInOutMetadata(ctx, metadata.New(map[string]string{DefaultXRequestIDKey: newRequestID(ctx, key)}))
	}
	return ctx
}

// to chain multiple request ids by generating new request id for each request and concatenating it to original request ids.
// key is RequestID within Context if have
func newContextForHandleRequestIDChain(ctx context.Context, key interface{}) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return appendInOutMetadata(ctx, metadata.New(map[string]string{DefaultXRequestIDKey: newRequestIDChain(ctx, key)}))
	}
	requestIDs := md.Get(DefaultXRequestIDKey)
	if len(requestIDs) == 0 {
		return appendInOutMetadata(ctx, metadata.New(map[string]string{DefaultXRequestIDKey: newRequestIDChain(ctx, key)}))
	}
	requestIDs = append(requestIDs, uuid.New().String())

	return appendInOutMetadata(ctx, metadata_.New(DefaultXRequestIDKey, requestIDs...))
}

func appendInOutMetadata(ctx context.Context, mds ...metadata.MD) context.Context {
	return appendOutgoingMetadata(appendIncomingMetadata(ctx, mds...), mds...)
}

func appendIncomingMetadata(ctx context.Context, mds ...metadata.MD) context.Context {
	out, _ := metadata.FromIncomingContext(ctx)
	out = out.Copy()
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return metadata.NewIncomingContext(ctx, out)
}

func appendOutgoingMetadata(ctx context.Context, mds ...metadata.MD) context.Context {
	out, _ := metadata.FromOutgoingContext(ctx)
	out = out.Copy()
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return metadata.NewOutgoingContext(ctx, out)
}

func newRequestID(ctx context.Context, key interface{}) string {
	if key == nil {
		return uuid.New().String()
	}
	val := ctx.Value(key)
	if valStr, ok := val.(string); ok {
		return valStr
	}
	return uuid.New().String()
}

func newRequestIDChain(ctx context.Context, key interface{}) string {
	if key == nil {
		return uuid.New().String()
	}
	val := ctx.Value(key)
	if valStr, ok := val.(string); ok {
		return strings.Join([]string{valStr, uuid.New().String()}, ",")
	}
	return uuid.New().String()
}
