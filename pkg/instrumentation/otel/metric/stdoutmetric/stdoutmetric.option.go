// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdoutmetric

import (
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

//go:generate go-option -type=option --trim
type option struct {
	StdoutOptions         []stdoutmetric.Option
	PeriodicReaderOptions []sdkmetric.PeriodicReaderOption

	ExporterWrappers []func(exporter sdkmetric.Exporter) sdkmetric.Exporter
}

func (o *option) SetDefaults() {}
