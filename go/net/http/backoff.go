// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

var (
	// A regular expression to match the error returned by net/http when the
	// configured number of redirects is exhausted. This error isn't typed
	// specifically so we resort to matching on the error string.
	redirectsErrorRe = regexp.MustCompile(`stopped after \d+ redirects\z`)

	// A regular expression to match the error returned by net/http when the
	// scheme specified in the URL is invalid. This error isn't typed
	// specifically so we resort to matching on the error string.
	schemeErrorRe = regexp.MustCompile(`unsupported protocol scheme`)
)

// RetryAfter tries to parse Retry-After response header when a http.StatusTooManyRequests
// (HTTP Code 429) is found in the resp parameter. Hence, it will return the number of
// seconds the server states it may be ready to process more requests from this client.
// Don't retry if the error was due to too many redirects.
// Don't retry if the error was due to an invalid protocol scheme.
// Don't retry if the error was due to TLS cert verification failure.
// Don't retry if the http's StatusCode is http.StatusNotImplemented.
func RetryAfter(resp *http.Response, err error, defaultBackoff time.Duration) (backoff time.Duration, retry bool) {
	backoff = defaultBackoff
	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if sleepInSeconds, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64); err == nil {
				backoff = time.Second * time.Duration(sleepInSeconds)
			}
		}
	}

	if err != nil {
		if v, ok := err.(*url.Error); ok {
			// Don't retry if the error was due to too many redirects.
			if redirectsErrorRe.MatchString(v.Error()) {
				return backoff, false
			}

			// Don't retry if the error was due to an invalid protocol scheme.
			if schemeErrorRe.MatchString(v.Error()) {
				return backoff, false
			}

			// Don't retry if the error was due to TLS cert verification failure.
			if _, ok := v.Err.(x509.UnknownAuthorityError); ok {
				return backoff, false
			}
		}

		// The error is likely recoverable so retry.
		return backoff, true
	}

	if resp != nil {
		// 429 Too Many Requests is recoverable. Sometimes the server puts
		// a Retry-After response header to indicate when the server is
		// available to start processing request from client.
		if resp.StatusCode == http.StatusTooManyRequests {
			return backoff, true
		}

		// Check the response code. We retry on 500-range responses to allow
		// the server time to recover, as 500's are typically not permanent
		// errors and may relate to outages on the server side. This will catch
		// invalid response codes as well, like 0 and 999.
		if resp.StatusCode == 0 || (resp.StatusCode >= http.StatusInternalServerError && resp.StatusCode != http.StatusNotImplemented) {
			return backoff, true
		}
	}
	return backoff, false
}

// ReplaceHttpRequestBody replace Body and recalculate ContentLength
// If ContentLength should not be recalculated, save and restore it after ReplaceHttpRequestBody
func ReplaceHttpRequestBody(req *http.Request, body io.Reader) {
	if req.Body != nil {
		req.Body.Close()
	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	req.Body = rc
	req.ContentLength = 0
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		default:
			// This is where we'd set it to -1 (at least
			// if body != NoBody) to mean unknown, but
			// that broke people during the Go 1.8 testing
			// period. People depend on it being 0 I
			// guess. Maybe retry later. See Issue 18117.
		}
		// For client requests, Request.ContentLength of 0
		// means either actually 0, or unknown. The only way
		// to explicitly say that the ContentLength is zero is
		// to set the Body to nil. But turns out too much code
		// depends on NewRequest returning a non-nil Body,
		// so we use a well-known ReadCloser variable instead
		// and have the http package also treat that sentinel
		// variable to mean explicitly zero.
		if req.GetBody != nil && req.ContentLength == 0 {
			req.Body = http.NoBody
			req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		}
	}
}

// ClientInvoker is called by ClientInterceptor to complete RPCs.
type ClientInvoker func(req *http.Request, retry int) (*http.Response, error)

// ClientInterceptor intercepts the execution of a HTTP on the client.
// interceptors can be specified as a DoWithBackoffOption, using
// WithClientInterceptor() or WithChainClientInterceptor(), when DoWithBackoffOption.
// When a interceptor(s) is set, gRPC delegates all http invocations to the interceptor,
// and it is the responsibility of the interceptor to call invoker to complete the processing
// of the HTTP.
type ClientInterceptor func(req *http.Request, retry int, invoker ClientInvoker, opts ...DoWithBackoffOption) (resp *http.Response, err error)

type RetryAfterHandler func(resp *http.Response, err error, defaultBackoff time.Duration) (backoff time.Duration, retry bool)

// DoRetryHandler send an HTTP request with retry seq and returns an HTTP response, following
// policy (such as redirects, cookies, auth) as configured on the
// client.
type DoRetryHandler = ClientInvoker

var DefaultDoRetryHandler = func(req *http.Request, retry int) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

//go:generate go-option -type "doWithBackoff"
type doWithBackoff struct {
	DoRetryHandler           DoRetryHandler
	clientInterceptor        ClientInterceptor
	ChainClientInterceptors  []ClientInterceptor
	RetryAfter               RetryAfterHandler
	ExponentialBackOffOption []time_.ExponentialBackOffOption
}

func (o *doWithBackoff) SetDefault() {
	o.DoRetryHandler = DefaultDoRetryHandler
	o.RetryAfter = RetryAfter
}

// getClientInvoker recursively generate the chained client invoker.
func getClientInvoker(interceptors []ClientInterceptor, curr int, finalInvoker ClientInvoker, opts ...DoWithBackoffOption) ClientInvoker {
	if curr == len(interceptors)-1 {
		return finalInvoker
	}
	return func(req *http.Request, retry int) (*http.Response, error) {
		return interceptors[curr+1](req, retry, getClientInvoker(interceptors, curr+1, finalInvoker), opts...)
	}
}

func (o *doWithBackoff) Complete() {
	if o.DoRetryHandler == nil {
		o.DoRetryHandler = DefaultDoRetryHandler
	}
	interceptors := o.ChainClientInterceptors
	o.ChainClientInterceptors = nil
	// Prepend o.ClientInterceptor to the chaining interceptors if it exists, since ClientInterceptor will
	// be executed before any other chained interceptors.
	if o.clientInterceptor != nil {
		interceptors = append([]ClientInterceptor{o.clientInterceptor}, interceptors...)
	}
	var chainedInt ClientInterceptor
	if len(interceptors) == 0 {
		chainedInt = nil
	} else if len(interceptors) == 1 {
		chainedInt = interceptors[0]
	} else {
		chainedInt = func(req *http.Request, retry int, invoker ClientInvoker, opts ...DoWithBackoffOption) (resp *http.Response, err error) {
			return interceptors[0](req, retry, getClientInvoker(interceptors, 0, invoker), opts...)
		}
	}
	o.clientInterceptor = chainedInt
}

// DoWithBackoff will retry by exponential backoff if failed.
// If request is not rewindable, retry wil be skipped.
func DoWithBackoff(httpReq *http.Request, opts ...DoWithBackoffOption) (*http.Response, error) {
	var opt doWithBackoff
	opt.SetDefault()
	opt.ApplyOptions(opts...)
	if opt.RetryAfter == nil {
		opt.RetryAfter = RetryAfter
	}
	opt.Complete()

	var option []time_.ExponentialBackOffOption
	option = append(option, time_.WithExponentialBackOffOptionMaxElapsedCount(3))
	option = append(option, opt.ExponentialBackOffOption...)
	backoff := time_.NewDefaultExponentialBackOff(option...)
	rewindableErr := RequestWithBodyRewindable(httpReq)
	var retries int
	for {
		if retries > 0 && httpReq.GetBody != nil {
			newBody, err := httpReq.GetBody()
			if err != nil {
				return nil, err
			}
			httpReq.Body = newBody
		}
		var do = opt.DoRetryHandler
		httpDo := do
		if opt.clientInterceptor != nil {
			httpDo = func(req *http.Request, retry int) (*http.Response, error) {
				return opt.clientInterceptor(req, retry, do, opts...)
			}
		}
		resp, err := httpDo(httpReq, retries)

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
		case <-httpReq.Context().Done():
			timer.Stop()
			if err != nil {
				return nil, fmt.Errorf("http do canceled after retries %d: %w", retries, err)
			} else {
				return resp, nil
			}
		}
	}
}

func HeadWithBackoff(ctx context.Context, url string, opts ...DoWithBackoffOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}
	return DoWithBackoff(req, opts...)
}

func GetWithBackoff(ctx context.Context, url string, opts ...DoWithBackoffOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}
	return DoWithBackoff(req, opts...)
}

func PostWithBackoff(ctx context.Context, url, contentType string, body io.Reader, opts ...DoWithBackoffOption) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return DoWithBackoff(req, opts...)
}

func PostFormWithBackoff(ctx context.Context, url string, data url.Values, opts ...DoWithBackoffOption) (resp *http.Response, err error) {
	return PostWithBackoff(ctx, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), opts...)
}

// DoJson the same as HttpDo, but bind with json
func DoJson(httpReq *http.Request, req, resp interface{}) error {
	if req != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return err
		}
		reqBody := bytes.NewReader(data)
		httpReq.Header.Set("Content-Type", "application/json")
		ReplaceHttpRequestBody(httpReq, reqBody)
	}

	httpResp, err := DefaultDoRetryHandler(httpReq, 0)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	if resp == nil {
		return nil
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, resp)
}

// DoJsonWithBackoff the same as DoWithBackoff, but bind with json
func DoJsonWithBackoff(httpReq *http.Request, req, resp interface{}, opts ...DoWithBackoffOption) error {
	if req != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return err
		}
		reqBody := bytes.NewReader(data)
		httpReq.Header.Set("Content-Type", "application/json")
		ReplaceHttpRequestBody(httpReq, reqBody)
	}
	httpResp, err := DoWithBackoff(httpReq, opts...)

	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	if resp == nil {
		return nil
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, resp)
}
