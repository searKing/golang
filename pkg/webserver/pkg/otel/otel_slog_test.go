// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otel_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	otel_ "github.com/searKing/golang/pkg/webserver/pkg/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestOtelSlogHandler_AddsTraceID(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := otel_.NewSlogHandler(baseHandler)
	logger := slog.New(handler)

	// 创建一个 tracer 和 span
	tracer := trace.NewTracerProvider().Tracer("test")
	ctx, span := tracer.Start(t.Context(), "test-operation")
	defer span.End()

	logger.InfoContext(ctx, "test message", "key", "value")

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatal(err)
	}

	if _, ok := logEntry["trace_id"]; !ok {
		t.Error("Expected trace_id in log output")
	}
	if _, ok := logEntry["span_id"]; !ok {
		t.Error("Expected span_id in log output")
	}
}

func TestOtelSlogHandler_NoTraceIDWithoutSpan(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := otel_.NewSlogHandler(baseHandler)
	logger := slog.New(handler)

	// 没有 span 的 context
	logger.InfoContext(t.Context(), "test message")

	var logEntry map[string]any
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := logEntry["trace_id"]; ok {
		t.Error("Should not have trace_id without span")
	}
}

func TestOtelSlogHandler_NoTraceIDWithNoopSpan(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := otel_.NewSlogHandler(baseHandler)
	logger := slog.New(handler)

	logger.InfoContext(t.Context(), "test message")

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatal(err)
	}

	// noop span 是无效的，不应该添加 trace_id
	if _, ok := logEntry["trace_id"]; ok {
		t.Error("Should not have trace_id with noop span")
	}
}

func TestOtelSlogHandler_PreservesExistingTraceID(t *testing.T) {
	var buf bytes.Buffer
	baseHandler := slog.NewJSONHandler(&buf, nil)
	handler := otel_.NewSlogHandler(baseHandler)
	logger := slog.New(handler)

	tp := trace.NewTracerProvider()
	tracer := tp.Tracer("test")
	ctx, span := tracer.Start(t.Context(), "test-operation")
	defer span.End()

	logger.InfoContext(ctx, "test", slog.String("trace_id", "custom-trace-id"))

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatal(err)
	}

	if logEntry["trace_id"] != "custom-trace-id" {
		t.Errorf("Should preserve manually provided trace_id, got: %v", logEntry["trace_id"])
	}

	if _, ok := logEntry["span_id"]; ok {
		t.Error("Should not add span_id when trace_id already exists")
	}
}

func TestSetDefault(t *testing.T) {
	// Verify that setting the default to itself does not result in deadlock.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	defer func(w io.Writer) { log.SetOutput(w) }(log.Writer())
	log.SetOutput(io.Discard)
	go func() {
		slog.Info("A")
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
		// deadlock if replace as follows:
		//slog.SetDefault(slog.New(otel_.NewSlogHandler(slog.Default().Handler())))
		slog.Info("B")
		cancel()
	}()
	<-ctx.Done()
	if err := ctx.Err(); !errors.Is(err, context.Canceled) {
		t.Errorf("wanted canceled, got %v", err)
	}
}
