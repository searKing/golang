// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpmetrichttp

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	_ "google.golang.org/grpc/encoding/gzip" // open gzip
)

// OpenExporter opens a metric exporter specified by its metric exporter name and a
// exporter-specific data source name, usually consisting of at least a
// metric exporter name and connection information.
func OpenExporter(ctx context.Context, opts ...Option) (sdkmetric.Exporter, error) {
	var opt option
	opt.SetDefaults()
	opt.ApplyOptions(opts...)
	exporter, err := otlpmetrichttp.New(ctx, opt.OtlpOptions...)
	if err != nil {
		return nil, err
	}
	var exp sdkmetric.Exporter = exporter
	for _, wrapper := range opt.ExporterWrappers {
		exp = wrapper(exp)
	}
	return exp, nil
}
