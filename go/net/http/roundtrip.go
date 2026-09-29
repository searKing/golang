// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import "net/http"

// RoundTripFunc is an adapter to allow the use of
// ordinary functions as HTTP Transport handlers. If f is a function
// with the appropriate signature, RoundTripFunc(f) is a
// Handler that calls f.
type RoundTripFunc func(req *http.Request) (resp *http.Response, err error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return f(req)
}

// RoundTripDecorator is an interface representing the ability to decorate or wrap a http.RoundTripper.
type RoundTripDecorator interface {
	WrapRoundTrip(rt http.RoundTripper) http.RoundTripper
}

// RoundTripDecoratorFunc is an adapter to allow the use of ordinary functions as RoundTripDecorators.
type RoundTripDecoratorFunc func(rt http.RoundTripper) http.RoundTripper

func (f RoundTripDecoratorFunc) WrapRoundTrip(rt http.RoundTripper) http.RoundTripper {
	return f(rt)
}

// RoundTripDecorators defines a RoundTripDecorator slice.
type RoundTripDecorators []RoundTripDecorator

func (chain RoundTripDecorators) WrapRoundTrip(next http.RoundTripper) http.RoundTripper {
	for i := range chain {
		h := chain[len(chain)-1-i]
		if h != nil {
			next = h.WrapRoundTrip(next)
		}
	}
	return next
}
