// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

import "github.com/google/uuid"

// WithDistributedTracing will make it so that a wrapped driver is used that supports the opentracing API.
// Deprecated: remove trace options.
func WithDistributedTracing() DBOption {
	return DBOptionFunc(func(o *DB) {
		o.opts.UseTracedDriver = true
	})
}

// WithOmitArgsFromTraceSpans will make it so that query arguments are omitted from tracing spans.
// Deprecated: remove trace options.
func WithOmitArgsFromTraceSpans() DBOption {
	return DBOptionFunc(func(o *DB) {
		o.opts.OmitArgs = true
	})
}

// WithTraceOrphans will make it so that root spans will be created if a trace could not be found using
// opentracing's SpanFromContext method.
// Deprecated: remove trace options.
func WithTraceOrphans() DBOption {
	return DBOptionFunc(func(o *DB) {
		o.opts.TraceOrphans = true
	})
}

// WithRandomDriverName is specifically for writing tests as you can't register a driver with the same name more than once.
// Deprecated: remove trace options.
func WithRandomDriverName() DBOption {
	return DBOptionFunc(func(o *DB) {
		o.opts.ForcedDriverName = uuid.New().String()
	})
}
