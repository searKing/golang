// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func extractLoggingAttrs(ctx context.Context) []slog.Attr {
	return fieldsToAttrSlice(logging.ExtractFields(ctx))
}

func fieldsToAttrSlice(fields logging.Fields) []slog.Attr {
	var attrs []slog.Attr
	i := fields.Iterator()
	for i.Next() {
		k, v := i.At()
		attrs = append(attrs, slog.Any(k, v))
	}
	return attrs
}
