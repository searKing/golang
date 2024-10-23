// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package driver

import (
	"context"
	"net/url"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// ExporterURLOpener is the interface that must be implemented by a trace exporter
// driver.
// ExporterURLOpener represents types that can open metric exporters based on a URL.
// The opener must not modify the URL argument. OpenExporterURL must be safe to
// call from multiple goroutines.
//
// This interface is generally implemented by types in driver packages.
type ExporterURLOpener interface {
	// OpenExporterURL creates a new exporter for the given target.
	OpenExporterURL(ctx context.Context, u *url.URL) (sdktrace.SpanExporter, error)

	// Scheme returns the scheme supported by this exporter.
	// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
	Scheme() string
}
