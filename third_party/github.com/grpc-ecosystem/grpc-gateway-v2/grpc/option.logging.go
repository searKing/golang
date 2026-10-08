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

// WithHttpLoggingOption customizes logs of http calls, as WithLoggingOption does for gRPC calls.
func WithHttpLoggingOption(opts ...HttpLoggingOption) GatewayOption {
	return GatewayOptionFunc(func(gateway *Gateway) {
		gateway.opt.httpLoggingOpts = append(gateway.opt.httpLoggingOpts, opts...)
	})
}

// ExtractHttpLoggingOptions extract all [HttpLoggingOption] from the given options.
func ExtractHttpLoggingOptions(options ...GatewayOption) []HttpLoggingOption {
	var g Gateway
	g.ApplyOptions(options...)
	return g.opt.httpLoggingOpts
}
