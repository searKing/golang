// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpsql

import (
	"context"
	"database/sql/driver"
	"io"
	"reflect"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Compile time validation that our types implement the expected interfaces
var (
	_ driver.Rows                           = otlpRows{}
	_ driver.RowsColumnTypeDatabaseTypeName = otlpRows{}
	_ driver.RowsColumnTypeLength           = otlpRows{}
	_ driver.RowsColumnTypeNullable         = otlpRows{}
	_ driver.RowsColumnTypePrecisionScale   = otlpRows{}
	// Currently, the one exception is RowsColumnTypeScanType which does not have a
	// valid zero value. This interface is tested for and only enabled in case the
	// parent implementation supports it.
	//_ driver.RowsColumnTypeScanType         = otlpRows{}
	_ driver.RowsNextResultSet = otlpRows{}
)

// withRowsColumnTypeScanType is the same as the driver.RowsColumnTypeScanType
// interface except it omits the driver.Rows embedded interface.
// If the original driver.Rows implementation wrapped by otlpsql supports
// RowsColumnTypeScanType we enable the original method implementation in the
// returned driver.Rows from wrapRows by doing a composition with otlpRows.
type withRowsColumnTypeScanType interface {
	ColumnTypeScanType(index int) reflect.Type
}

// otlpRows implements driver.Rows and all enhancement interfaces except
// driver.RowsColumnTypeScanType.
type otlpRows struct {
	parent  driver.Rows
	ctx     context.Context
	options wrapper
}

//func (r otlpRows) ColumnTypeScanType(index int) reflect.Type {
//	if v, ok := r.parent.(driver.RowsColumnTypeScanType); ok {
//		return v.ColumnTypeScanType(index)
//	}
//
//	return reflect.TypeOf(new(interface{}))
//}

// HasNextResultSet calls the implements the driver.RowsNextResultSet for otlpRows.
// It returns the the underlying result of HasNextResultSet from the otlpRows.parent
// if the parent implements driver.RowsNextResultSet.
func (r otlpRows) HasNextResultSet() bool {
	if v, ok := r.parent.(driver.RowsNextResultSet); ok {
		return v.HasNextResultSet()
	}

	return false
}

// NextResultSet calls the implements the driver.RowsNextResultSet for otlpRows.
// It returns the the underlying result of NextResultSet from the otlpRows.parent
// if the parent implements driver.RowsNextResultSet.
func (r otlpRows) NextResultSet() error {
	if v, ok := r.parent.(driver.RowsNextResultSet); ok {
		return v.NextResultSet()
	}

	return io.EOF
}

// ColumnTypeDatabaseTypeName calls the implements the driver.RowsColumnTypeDatabaseTypeName for otlpRows.
// It returns the the underlying result of ColumnTypeDatabaseTypeName from the otlpRows.parent
// if the parent implements driver.RowsColumnTypeDatabaseTypeName.
func (r otlpRows) ColumnTypeDatabaseTypeName(index int) string {
	if v, ok := r.parent.(driver.RowsColumnTypeDatabaseTypeName); ok {
		return v.ColumnTypeDatabaseTypeName(index)
	}

	return ""
}

// ColumnTypeLength calls the implements the driver.RowsColumnTypeLength for otlpRows.
// It returns the the underlying result of ColumnTypeLength from the otlpRows.parent
// if the parent implements driver.RowsColumnTypeLength.
func (r otlpRows) ColumnTypeLength(index int) (length int64, ok bool) {
	if v, ok := r.parent.(driver.RowsColumnTypeLength); ok {
		return v.ColumnTypeLength(index)
	}

	return 0, false
}

// ColumnTypeNullable calls the implements the driver.RowsColumnTypeNullable for otlpRows.
// It returns the the underlying result of ColumnTypeNullable from the otlpRows.parent
// if the parent implements driver.RowsColumnTypeNullable.
func (r otlpRows) ColumnTypeNullable(index int) (nullable, ok bool) {
	if v, ok := r.parent.(driver.RowsColumnTypeNullable); ok {
		return v.ColumnTypeNullable(index)
	}

	return false, false
}

// ColumnTypePrecisionScale calls the implements the driver.RowsColumnTypePrecisionScale for otlpRows.
// It returns the the underlying result of ColumnTypePrecisionScale from the otlpRows.parent
// if the parent implements driver.RowsColumnTypePrecisionScale.
func (r otlpRows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	if v, ok := r.parent.(driver.RowsColumnTypePrecisionScale); ok {
		return v.ColumnTypePrecisionScale(index)
	}

	return 0, 0, false
}

func (r otlpRows) Columns() []string {
	return r.parent.Columns()
}

func (r otlpRows) Close() (err error) {
	if r.options.RowsClose {
		attrs := append([]attribute.KeyValue(nil), r.options.DefaultAttributes...)
		ctx := r.ctx
		onDeferWithErr := recordCallStats("go.sql.rows.close", r.options.InstanceName)
		defer func() {
			// Invoking this function in a defer so that we can capture
			// the value of err as set on function exit.
			onDeferWithErr(ctx, err, attrs...)
		}()

		parentSpan := trace.SpanFromContext(ctx)
		if !r.options.AllowRoot && !parentSpan.SpanContext().IsValid() {
			// we already tested driver
			return r.parent.Close()
		}

		ctx, span := otel.Tracer("").Start(ctx, "sql:exec", trace.WithSpanKind(trace.SpanKindClient))
		defer func() {
			setSpanStatus(span, r.options, err)
			span.SetAttributes(attrs...)
			span.End()
		}()
	}

	err = r.parent.Close()
	return
}

func (r otlpRows) Next(dest []driver.Value) (err error) {
	if r.options.RowsNext {
		attrs := append([]attribute.KeyValue(nil), r.options.DefaultAttributes...)
		ctx := r.ctx
		onDeferWithErr := recordCallStats("go.sql.rows.next", r.options.InstanceName)
		defer func() {
			// Invoking this function in a defer so that we can capture
			// the value of err as set on function exit.
			onDeferWithErr(ctx, err, attrs...)
		}()

		parentSpan := trace.SpanFromContext(ctx)
		if !r.options.AllowRoot && !parentSpan.SpanContext().IsValid() {
			// we already tested driver
			return r.parent.Close()
		}

		ctx, span := otel.Tracer("").Start(ctx, "sql:exec", trace.WithSpanKind(trace.SpanKindClient))
		defer func() {
			if err == io.EOF {
				// not an error; expected to happen during iteration
				setSpanStatus(span, r.options, nil)
			} else {
				setSpanStatus(span, r.options, err)
			}
			span.SetAttributes(attrs...)
			span.End()
		}()
	}

	err = r.parent.Next(dest)
	return
}

// wrapRows returns a struct which conforms to the driver.Rows interface.
// otlpRows implements all enhancement interfaces that have no effect on
// sql/database logic in case the underlying parent implementation lacks them.
// Currently, the one exception is RowsColumnTypeScanType which does not have a
// valid zero value. This interface is tested for and only enabled in case the
// parent implementation supports it.
func wrapRows(ctx context.Context, parent driver.Rows, options wrapper) driver.Rows {
	var (
		ts, hasColumnTypeScan = parent.(driver.RowsColumnTypeScanType)
	)

	r := otlpRows{
		parent:  parent,
		ctx:     ctx,
		options: options,
	}

	if hasColumnTypeScan {
		return struct {
			otlpRows
			withRowsColumnTypeScanType
		}{r, ts}
	}

	return r
}
