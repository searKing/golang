// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import "log"

// RejectInsecureWithErrorLog specifies an optional logger for errors
func RejectInsecureWithErrorLog(l *log.Logger) RejectInsecureOption {
	return RejectInsecureOptionFunc(func(r *rejectInsecure) {
		r.ErrorLog = l
	})
}

// RejectInsecureWithForceHttp specifies whether to allow any request, as a shortcut circuit
func RejectInsecureWithForceHttp(forceHttp bool) RejectInsecureOption {
	return RejectInsecureOptionFunc(func(r *rejectInsecure) {
		r.ForceHttp = forceHttp
	})
}

// RejectInsecureWithAllowedTlsCidrs specifies whether to allow any request which client or proxy's ip included
// a cidr is a CIDR notation IP address and prefix length,
// like "192.0.2.0/24" or "2001:db8::/32", as defined in
// RFC 4632 and RFC 4291.
func RejectInsecureWithAllowedTlsCidrs(allowedTLSCIDRs []string) RejectInsecureOption {
	return RejectInsecureOptionFunc(func(r *rejectInsecure) {
		r.AllowedTlsCidrs = allowedTLSCIDRs
	})
}

// RejectInsecureWithWhitelistedPaths specifies whether to allow any request which http path matches
func RejectInsecureWithWhitelistedPaths(whitelistedPaths []string) RejectInsecureOption {
	return RejectInsecureOptionFunc(func(r *rejectInsecure) {
		r.WhitelistedPaths = whitelistedPaths
	})
}
