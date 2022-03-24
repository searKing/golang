package sql

import "github.com/google/uuid"

// WithDistributedTracing will make it so that a wrapped driver is used that supports the opentracing API.
func WithDistributedTracing() DBOption {
	return DBOptionFunc(func(o *DB) {
		o.opts.UseTracedDriver = true
	})
}

// WithOmitArgsFromTraceSpans will make it so that query arguments are omitted from tracing spans.
func WithOmitArgsFromTraceSpans() DBOption {
	return DBOptionFunc(func(o *DB) {
		o.opts.OmitArgs = true
	})
}

// WithTraceOrphans will make it so that root spans will be created if a trace could not be found using
// opentracing's SpanFromContext method.
func WithTraceOrphans() DBOption {
	return DBOptionFunc(func(o *DB) {
		o.opts.TraceOrphans = true
	})
}

// WithRandomDriverName is specifically for writing tests as you can't register a driver with the same name more than once.
func WithRandomDriverName() DBOption {
	return DBOptionFunc(func(o *DB) {
		o.opts.ForcedDriverName = uuid.New().String()
	})
}
