// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	http_ "github.com/searKing/golang/go/net/http"
)

// ExtractServerOptions extracts all [grpc.ServerOption] from the given options.
// The returned slice is the fully assembled server options, including chained interceptors,
// identical to what [NewGatewayTLS] would use.
func ExtractServerOptions(options ...GatewayOption) []grpc.ServerOption {
	var g Gateway
	g.ApplyOptions(options...)
	return g.opt.ServerOptions()
}

// ExtractDialOptions extracts all [grpc.DialOption] from the given options.
// The returned slice is the fully assembled dial options, including chained interceptors,
// identical to what [NewGatewayTLS] would use for the internal gRPC client.
func ExtractDialOptions(options ...GatewayOption) []grpc.DialOption {
	var g Gateway
	g.ApplyOptions(options...)
	return g.opt.ClientDialOpts()
}

// ExtractServeMuxOptions extracts all [runtime.ServeMuxOption] from the given options.
func ExtractServeMuxOptions(options ...GatewayOption) []runtime.ServeMuxOption {
	var g Gateway
	g.ApplyOptions(options...)
	return g.opt.srvMuxOpts
}

// ExtractHttpInterceptorChain extracts the [http_.HandlerInterceptorChain] from the given options.
func ExtractHttpInterceptorChain(options ...GatewayOption) http_.HandlerInterceptorChain {
	var g Gateway
	g.ApplyOptions(options...)
	return g.opt.httpInterceptors
}

// ExtractRoundTripDecorators extracts the [http_.RoundTripDecorators] from the given options.
func ExtractRoundTripDecorators(options ...GatewayOption) http_.RoundTripDecorators {
	var g Gateway
	g.ApplyOptions(options...)
	return g.opt.roundTripDecorators
}
