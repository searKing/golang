// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpsql

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func paramsAttr(args []driver.Value) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, len(args))
	for i, arg := range args {
		key := "sql.arg" + strconv.Itoa(i)
		attrs = append(attrs, argToAttr(key, arg))
	}
	return attrs
}

func namedParamsAttr(args []driver.NamedValue) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, len(args))
	for _, arg := range args {
		var key string
		if arg.Name != "" {
			key = arg.Name
		} else {
			key = "sql.arg." + strconv.Itoa(arg.Ordinal)
		}
		attrs = append(attrs, argToAttr(key, arg.Value))
	}
	return attrs
}

func argToAttr(key string, val any) attribute.KeyValue {
	switch v := val.(type) {
	case nil:
		return attribute.String(key, "")
	case int64:
		return attribute.Int64(key, v)
	case float64:
		return attribute.Float64(key, v)
	case bool:
		return attribute.Bool(key, v)
	case []byte:
		if len(v) > 256 {
			v = v[0:256]
		}
		return attribute.String(key, fmt.Sprintf("%s", v))
	default:
		s := fmt.Sprintf("%v", v)
		if len(s) > 256 {
			s = s[0:256]
		}
		return attribute.String(key, s)
	}
}

func setSpanStatus(span trace.Span, opts wrapper, err error) {
	switch {
	case err == nil:
		span.SetStatus(codes.Ok, "")
		return
	case errors.Is(err, driver.ErrSkip):
		span.SetStatus(codes.Unset, err.Error())
		if opts.DisableErrSkip {
			// Suppress driver.ErrSkip since at runtime some drivers might not have
			// certain features, and an error would pollute many spans.
			span.SetStatus(codes.Ok, err.Error())
		}
	default:
		span.SetStatus(codes.Error, err.Error())
	}
}
