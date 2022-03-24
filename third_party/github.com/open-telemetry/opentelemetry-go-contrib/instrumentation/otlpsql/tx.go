// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpsql

import (
	"context"
	"database/sql/driver"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Compile time validation that our types implement the expected interfaces
var (
	_ driver.Tx = otlpTx{}
)

// otlpTx implements driver.Tx
type otlpTx struct {
	parent  driver.Tx
	ctx     context.Context
	options wrapper
}

func (t otlpTx) Commit() (err error) {
	ctx := t.ctx
	attrs := append([]attribute.KeyValue(nil), t.options.DefaultAttributes...)
	onDeferWithErr := recordCallStats("go.sql.tx.commit", t.options.InstanceName)
	defer func() {
		// Invoking this function in a defer so that we can capture
		// the value of err as set on function exit.
		onDeferWithErr(ctx, err, attrs...)
	}()

	parentSpan := trace.SpanFromContext(ctx)
	if !t.options.AllowRoot && !parentSpan.SpanContext().IsValid() {
		// we already tested driver
		return t.parent.Commit()
	}

	ctx, span := otel.Tracer("").Start(ctx, "sql:commit", trace.WithSpanKind(trace.SpanKindClient))
	defer func() {
		setSpanStatus(span, t.options, err)
		span.SetAttributes(attrs...)
		span.End()
	}()
	span.SetAttributes(attrs...)

	err = t.parent.Commit()
	return
}

func (t otlpTx) Rollback() (err error) {
	ctx := t.ctx
	var attrs []attribute.KeyValue
	onDeferWithErr := recordCallStats("go.sql.tx.rollback", t.options.InstanceName)
	defer func() {
		// Invoking this function in a defer so that we can capture
		// the value of err as set on function exit.
		onDeferWithErr(ctx, err, attrs...)
	}()

	parentSpan := trace.SpanFromContext(ctx)
	if !t.options.AllowRoot && !parentSpan.SpanContext().IsValid() {
		// we already tested driver
		return t.parent.Commit()
	}

	ctx, span := otel.Tracer("").Start(ctx, "sql:rollback", trace.WithSpanKind(trace.SpanKindClient))
	defer func() {
		setSpanStatus(span, t.options, err)
		span.SetAttributes(attrs...)
		span.End()
	}()
	span.SetAttributes(attrs...)
	err = t.parent.Rollback()
	return
}
