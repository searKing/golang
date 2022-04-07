// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import "net/http"

// RoundTripFunc is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type RoundTripFunc func(req *http.Request) (resp *http.Response, err error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return f(req)
}

type RoundTripHandler = http.RoundTripper

// handlersChain defines a HandlerFunc array.
type handlersChain []RoundTripHandler

// Last returns the last handler in the chain. ie. the last handler is the main own.
func (c handlersChain) Last() RoundTripHandler {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
