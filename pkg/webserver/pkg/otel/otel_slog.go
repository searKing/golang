// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otel

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// NewSlogHandler wraps an existing slog.Handler to automatically add
// OpenTelemetry trace_id and span_id to log records when a valid span
// is present in the context.
//
// Pass a handler that writes to its own io.Writer, such as
// slog.NewJSONHandler or slog.NewTextHandler. Do not pass
// slog.Default().Handler(), or a wrapper that delegates to it, when the
// result will be installed with slog.SetDefault. Until the process replaces
// it, that handler is slog's builtin handler: it writes through log.Default,
// and SetDefault points log.Default back at the wrapper, which deadlocks.
// A caller-defined Handler type does not avoid this when its inner handler
// is still the builtin one.
// See https://github.com/golang/go/issues/61892 and
// https://github.com/golang/go/issues/77716.
//
// Example usage:
//
//	baseHandler := slog.NewJSONHandler(os.Stdout, nil)
//	handler := otel.NewSlogHandler(baseHandler)
//	logger := slog.New(handler)
//
//	// Logs will automatically include trace_id and span_id as attrs if context has a span
//	logger.InfoContext(ctx, "processing request", "user_id", 123)
func NewSlogHandler(handler slog.Handler) slog.Handler {
	if handler == nil {
		return nil
	}
	return &otelSlogHandler{
		handler: handler,
	}
}

const (
	traceIDKey = "trace_id"
	spanIDKey  = "span_id"
)

var _ slog.Handler = &otelSlogHandler{}

// otelSlogHandler is a slog.Handler that adds OpenTelemetry trace IDs to log records.
type otelSlogHandler struct {
	handler slog.Handler
}

func (osh *otelSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return osh.handler.Enabled(ctx, level)
}

func (osh *otelSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		var hasTraceIds bool
		record.Attrs(func(attr slog.Attr) bool {
			if attr.Key == traceIDKey || attr.Key == spanIDKey {
				hasTraceIds = true
				return false
			}
			return true
		})
		if !hasTraceIds {
			record.AddAttrs(slog.String(traceIDKey, span.SpanContext().TraceID().String()),
				slog.String(spanIDKey, span.SpanContext().SpanID().String()))
		}
	}

	return osh.handler.Handle(ctx, record)
}

func (osh *otelSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewSlogHandler(osh.handler.WithAttrs(attrs))
}

func (osh *otelSlogHandler) WithGroup(name string) slog.Handler {
	return NewSlogHandler(osh.handler.WithGroup(name))
}
