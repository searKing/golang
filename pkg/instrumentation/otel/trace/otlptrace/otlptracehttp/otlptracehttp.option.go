// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otlptracehttp

import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

//go:generate go-option -type=option --trim
type option struct {
	OtlpOptions []otlptracehttp.Option
}

func (o *option) SetDefaults() {}