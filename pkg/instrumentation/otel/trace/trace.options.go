// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

//go:generate go-option -type=option --trim
type option struct {
	// ExporterEndpoints is the target endpoint URL (scheme, host, port, path) the Exporter will connect to.
	ExporterEndpoints []string
	Exporters         []sdktrace.SpanExporter
	ResourceAttrs     []attribute.KeyValue
	SpanAttrs         []attribute.KeyValue

	TracerProviderOptions []sdktrace.TracerProviderOption
}

func (o *option) SetDefaults() {}
