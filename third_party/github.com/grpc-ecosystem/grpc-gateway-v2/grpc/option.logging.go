// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

func WithLoggingOption(opts ...logging.Option) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.loggingOpts = append(gateway.opt.loggingOpts, opts...)
	})
}

// ExtractLoggingOptions extract all [logging.Option] from the given options.
func ExtractLoggingOptions(options ...GatewayOption) []logging.Option {
	var g Gateway
	g.ApplyOptions(options...)
	return g.opt.loggingOpts
}
