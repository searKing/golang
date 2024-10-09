// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdoutmetric

import (
	"context"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// OpenExporter opens a metric exporter specified by its metric exporter name and a
// exporter-specific data source name, usually consisting of at least a
// metric exporter name and connection information.
func OpenExporter(ctx context.Context, opts ...Option) (sdkmetric.Exporter, error) {
	var opt option
	opt.SetDefaults()
	opt.ApplyOptions(opts...)
	exporter, err := stdoutmetric.New(opt.StdoutOptions...)
	if err != nil {
		return nil, err
	}
	for _, wrapper := range opt.ExporterWrappers {
		exporter = wrapper(exporter)
	}
	return exporter, nil
}
