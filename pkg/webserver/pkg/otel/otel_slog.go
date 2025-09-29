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
