// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metric

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
)

//go:generate go-option -type=option --trim
type option struct {
	// ExporterEndpoints is the target endpoint URL (scheme, host, port, path) the Exporter will connect to.
	ExporterEndpoints []string
	Readers           []metric.Reader
	ResourceAttrs     []attribute.KeyValue

	MetricOptions               []metric.Option
	MetricPeriodicReaderOptions []metric.PeriodicReaderOption
}

func (o *option) SetDefaults() {}

// Syntax sugar.

// WithOptionMetricCollectPeriod  configures the intervening time between exports for a
// PeriodicReader.
// MetricCollectPeriod is the interval between calls to Collect a
// checkpoint.
// When pulling metrics and not exporting, this is the minimum
// time between calls to Collect.  In a pull-only
// configuration, collection is performed on demand; set
// CollectPeriod to 0 always recompute the export record set.
//
// When exporting metrics, this must be > 0.
//
// Default value is 60s.
func WithOptionMetricCollectPeriod(v time.Duration) Option {
	return OptionFunc(func(o *option) {
		o.MetricPeriodicReaderOptions = append(o.MetricPeriodicReaderOptions, metric.WithInterval(v))
	})
}
