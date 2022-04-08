// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"fmt"
	"net/http"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

// RoundTripperWithBackoff wraps http.RoundTripper retryable by backoff.
func RoundTripperWithBackoff(rt http.RoundTripper, opts ...DoWithBackoffOption) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
		var opt doWithBackoff
		opt.SetDefault()
		if rt == nil {
			rt = http.DefaultTransport
		}
		opt.DoRetryHandler = func(req *http.Request, retry int) (*http.Response, error) { return rt.RoundTrip(req) }

		opt.ApplyOptions(opts...)
		if opt.RetryAfter == nil {
			opt.RetryAfter = RetryAfter
		}
		opt.Complete()

		var option []time_.ExponentialBackOffOption
		option = append(option, time_.WithExponentialBackOffOptionMaxElapsedCount(3))
		option = append(option, opt.ExponentialBackOffOption...)
		backoff := time_.NewDefaultExponentialBackOff(option...)
		rewindableErr := RequestWithBodyRewindable(req)
		var retries int
		for {
			if retries > 0 && req.GetBody != nil {
				newBody, err := req.GetBody()
				if err != nil {
					return nil, err
				}
				req.Body = newBody
			}
			var do = opt.DoRetryHandler
			httpDo := do
			if opt.clientInterceptor != nil {
				httpDo = func(req *http.Request, retry int) (*http.Response, error) {
					return opt.clientInterceptor(req, retry, do, opts...)
				}
			}
			resp, err := httpDo(req, retries)

			wait, ok := backoff.NextBackOff()
			if !ok {
				if err != nil {
					return nil, fmt.Errorf("http do reach backoff limit after retries %d: %w", retries, err)
				} else {
					return resp, nil
				}
			}

			wait, retry := opt.RetryAfter(resp, err, wait)
			if !retry {
				if err != nil {
					return nil, fmt.Errorf("http do reach server limit after retries %d: %w", retries, err)
				} else {
					return resp, nil
				}
			}

			if rewindableErr != nil {
				if err != nil {
					return nil, fmt.Errorf("http do cannot rewindbody after retries %d: %w", retries, err)
				} else {
					return resp, nil
				}
			}

			timer := time.NewTimer(wait)
			select {
			case <-timer.C:
				retries++
				continue
			case <-req.Context().Done():
				timer.Stop()
				if err != nil {
					return nil, fmt.Errorf("http do canceled after retries %d: %w", retries, err)
				} else {
					return resp, nil
				}
			}
		}
	})
}
