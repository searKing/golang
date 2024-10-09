// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdouttrace

import (
	"context"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

// OpenExporter opens a trace exporter specified by its trace exporter name and a
// exporter-specific data source name, usually consisting of at least a
// trace exporter name and connection information.
func OpenExporter(ctx context.Context, opts ...Option) (*stdouttrace.Exporter, error) {
	var opt option
	opt.SetDefaults()
	opt.ApplyOptions(opts...)
	return stdouttrace.New(opt.StdoutOptions...)
}
