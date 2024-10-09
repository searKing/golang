// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheusmetric

import (
	"context"

	prometheusmetric "go.opentelemetry.io/otel/exporters/prometheus"
)

// OpenReader opens a metric reader specified by its metric exporter name and a
// exporter-specific data source name, usually consisting of at least a
// metric exporter name and connection information.
func OpenReader(ctx context.Context, opts ...Option) (*prometheusmetric.Exporter, error) {
	var opt option
	opt.SetDefaults()
	opt.ApplyOptions(opts...)
	return prometheusmetric.New(opt.PrometheusOptions...) // such as prometheus, that's a manual puller actually
}
