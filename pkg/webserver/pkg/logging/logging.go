// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logging

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func Attrs[T slog.Attr | any](ctx context.Context) []T {
	return fieldsToAttrSlice[T](logging.ExtractFields(ctx))
}

func fieldsToAttrSlice[T slog.Attr | any](fields logging.Fields) []T {
	var attrs []T
	i := fields.Iterator()
	for i.Next() {
		k, v := i.At()
		attrs = append(attrs, (any(slog.Any(k, v))).(T))
	}
	return attrs
}
