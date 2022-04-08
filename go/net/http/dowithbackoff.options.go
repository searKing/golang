// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"

	time_ "github.com/searKing/golang/go/time"
)

// WithDoWithBackoffOptionChainUnaryInterceptor returns a DoWithBackoffOption that specifies the chained
// interceptor for http clients. The first interceptor will be the outer most,
// while the last interceptor will be the inner most wrapper around the real call.
// All interceptors added by this method will be chained, and the interceptor
// defined by WithClientInterceptor will always be prepended to the chain.
func WithDoWithBackoffOptionChainUnaryInterceptor(interceptors ...ClientInterceptor) DoWithBackoffOption {
	return DoWithBackoffOptionFunc(func(o *doWithBackoff) {
		o.ChainClientInterceptors = append(o.ChainClientInterceptors, interceptors...)
	})
}

func WithDoWithBackoffOptionRetryAfter(f RetryAfterHandler) DoWithBackoffOption {
	return DoWithBackoffOptionFunc(func(o *doWithBackoff) {
		o.RetryAfter = f
	})
}

func WithDoWithBackoffOptionDoRetryHandler(f DoRetryHandler) DoWithBackoffOption {
	return DoWithBackoffOptionFunc(func(o *doWithBackoff) {
		o.DoRetryHandler = f
	})
}

func WithDoWithBackoffOptionExponentialBackOffOption(opts ...time_.ExponentialBackOffOption) DoWithBackoffOption {
	return DoWithBackoffOptionFunc(func(o *doWithBackoff) {
		o.ExponentialBackOffOption = append(o.ExponentialBackOffOption, opts...)
	})
}

func WithDoWithBackoffOptionMaxRetries(retries int) DoWithBackoffOption {
	return DoWithBackoffOptionFunc(func(o *doWithBackoff) {
		o.ExponentialBackOffOption = append(o.ExponentialBackOffOption, time_.WithExponentialBackOffOptionMaxElapsedCount(retries))
	})
}

func WithDoWithBackoffOptionGrpcBackOff(retries int) DoWithBackoffOption {
	return DoWithBackoffOptionFunc(func(o *doWithBackoff) {
		o.ExponentialBackOffOption = append(o.ExponentialBackOffOption, time_.WithExponentialBackOffOptionGRPC())
	})
}

// WithDoWithBackoffOptionRoundTripper returns a DoWithBackoffOption.
func WithDoWithBackoffOptionRoundTripper(rt http.RoundTripper) DoWithBackoffOption {
	return WithDoWithBackoffOptionDoRetryHandler(
		func(req *http.Request, retry int) (*http.Response, error) {
			return rt.RoundTrip(req)
		})
}

// WithDoWithBackoffOptionProxy returns a DoWithBackoffOption.
func WithDoWithBackoffOptionProxy(proxyUsed string, targetAsProxy bool) DoWithBackoffOption {
	if proxyUsed == "" {
		return EmptyDoWithBackoffOption{}
	}
	cli := NewClientWithTarget(proxyUsed, targetAsProxy)
	return WithDoWithBackoffOptionDoRetryHandler(
		func(req *http.Request, retry int) (*http.Response, error) {
			return cli.Do(req)
		})
}
