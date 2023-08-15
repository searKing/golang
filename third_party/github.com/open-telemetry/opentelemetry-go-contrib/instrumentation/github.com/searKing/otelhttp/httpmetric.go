// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelgrpc

import (
	otelcontrib "go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel/attribute"
)

var (
	// InstrumentationName is the name of this instrumentation package.
	InstrumentationName = "github.com/searKing/golang/third_party/github.com/open-telemetry/opentelemetry-go-contrib/instrumentation/github.com/searKing/octelhttp"
	// InstrumentationVersion is the version of this instrumentation package.
	InstrumentationVersion = otelcontrib.Version()

	// AttrsFilter is a filter before Report
	AttrsFilter = func(attrs ...attribute.KeyValue) []attribute.KeyValue {
		return attrs
	}
)
