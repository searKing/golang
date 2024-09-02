// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package requestid

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	reflect_ "github.com/searKing/golang/go/reflect"
	strings_ "github.com/searKing/golang/go/strings"
)

const structFieldNameRequestId = "RequestId"

func tryRetrieveRequestId(v any) string {
	if reflect_.IsNil(v) {
		return ""
	}

	if v, ok := v.(interface{ GetRequestId() string }); ok {
		return v.GetRequestId()
	}

	field, has := reflect_.FieldByNames(reflect.ValueOf(v), structFieldNameRequestId)
	if !has {
		return ""
	}
	if v, ok := field.Interface().(string); ok {
		return v
	}
	return ""
}

func trySetRequestId(v any, id string, ignoreEmpty bool) {
	if reflect_.IsNil(v) {
		return
	}
	if ignoreEmpty {
		id := tryRetrieveRequestId(v)
		if id != "" {
			return
		}
	}
	reflect_.SetFieldByNames(reflect.ValueOf(v), []string{structFieldNameRequestId}, reflect.ValueOf(id))
}

func tagLoggingRequestId(ctx context.Context, v any) (context.Context, string) {
	id := tryRetrieveRequestId(v)
	if id == "" {
		id = strings_.ValueOrDefault(extractServerMetadataRequestId(ctx)...)
		if id == "" {
			id = uuid.NewString()
		}
		trySetRequestId(v, id, false)
	}
	return logging.InjectFields(ctx, logging.Fields{"request_id", id}), id
}

func extractServerMetadataRequestId(ctx context.Context) []string {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok || md.HeaderMD == nil {
		return nil
	}
	return md.HeaderMD.Get(requestId)
}
