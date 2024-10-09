// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlpmetricgrpc

import (
	"context"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// OpenReader opens a metric reader specified by its metric exporter name and a
// exporter-specific data source name, usually consisting of at least a
// metric exporter name and connection information.
func OpenReader(ctx context.Context, opts ...Option) (sdkmetric.Reader, error) {
	var opt option
	opt.SetDefaults()
	opt.ApplyOptions(opts...)
	exporter, err := OpenExporter(ctx, opts...)
	if err != nil {
		return nil, err
	}
	// handle exporter, as periodic pusher
	return sdkmetric.NewPeriodicReader(exporter, opt.PeriodicReaderOptions...), nil
}
