// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdouttrace

import "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

//go:generate go-option -type=option --trim
type option struct {
	StdoutOptions []stdouttrace.Option
}

func (o *option) SetDefaults() {}
