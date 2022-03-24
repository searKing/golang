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
	_ driver.Result = otlpResult{}
)

// otlpResult implements driver.Result
type otlpResult struct {
	parent  driver.Result
	ctx     context.Context
	options wrapper
}

func (r otlpResult) LastInsertId() (id int64, err error) {
	if r.options.LastInsertID {
		attrs := append([]attribute.KeyValue(nil), r.options.DefaultAttributes...)
		ctx := r.ctx
		onDeferWithErr := recordCallStats("go.sql.result.last_insert_id", r.options.InstanceName)
		defer func() {
			// Invoking this function in a defer so that we can capture
			// the value of err as set on function exit.
			onDeferWithErr(ctx, err, attrs...)
		}()

		parentSpan := trace.SpanFromContext(ctx)
		if !r.options.AllowRoot && !parentSpan.SpanContext().IsValid() {
			// we already tested driver
			return r.parent.LastInsertId()
		}

		ctx, span := otel.Tracer("").Start(ctx, "sql:exec", trace.WithSpanKind(trace.SpanKindClient))
		defer func() {
			setSpanStatus(span, r.options, err)
			span.SetAttributes(attrs...)
			span.End()
		}()
		span.SetAttributes(attrs...)
	}

	id, err = r.parent.LastInsertId()
	return
}

func (r otlpResult) RowsAffected() (cnt int64, err error) {
	if r.options.RowsAffected {
		attrs := append([]attribute.KeyValue(nil), r.options.DefaultAttributes...)
		ctx := r.ctx
		onDeferWithErr := recordCallStats("go.sql.result.rows_affected", r.options.InstanceName)
		defer func() {
			// Invoking this function in a defer so that we can capture
			// the value of err as set on function exit.
			onDeferWithErr(ctx, err, attrs...)
		}()

		parentSpan := trace.SpanFromContext(ctx)
		if !r.options.AllowRoot && !parentSpan.SpanContext().IsValid() {
			// we already tested driver
			return r.parent.LastInsertId()
		}

		ctx, span := otel.Tracer("").Start(ctx, "sql:exec", trace.WithSpanKind(trace.SpanKindClient))
		defer func() {
			setSpanStatus(span, r.options, err)
			span.SetAttributes(attrs...)
			span.End()
		}()
		span.SetAttributes(attrs...)
	}

	cnt, err = r.parent.RowsAffected()
	return
}
