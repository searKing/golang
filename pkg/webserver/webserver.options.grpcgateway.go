// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webserver

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/searKing/golang/pkg/webserver/pkg/otel"
)

func (f *Factory) ServeMuxOptions(opts ...runtime.ServeMuxOption) []runtime.ServeMuxOption {
	if f.fc.OtelHandling {
		opts = append(opts, otel.ServeMuxOptions()...)
	}
	return opts
}
