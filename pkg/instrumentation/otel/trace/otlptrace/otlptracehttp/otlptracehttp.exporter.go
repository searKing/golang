// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlptracehttp

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	_ "google.golang.org/grpc/encoding/gzip" // open gzip
)

// OpenExporter opens a trace exporter specified by its trace exporter name and a
// exporter-specific data source name, usually consisting of at least a
// trace exporter name and connection information.
func OpenExporter(ctx context.Context, opts ...Option) (*otlptrace.Exporter, error) {
	var opt option
	opt.SetDefaults()
	opt.ApplyOptions(opts...)
	return otlptracehttp.New(ctx, opt.OtlpOptions...)
}
