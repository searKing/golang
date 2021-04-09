// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

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

// Do sends an HTTP request and returns an HTTP response, following
// policy (such as redirects, cookies, auth) as configured on the
// client.
func Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

// DoWithBackoff will retry by exponential backoff if failed.
// If request is not rewindable, retry wil be skipped.
func DoWithBackoff(httpReq *http.Request, opts ...time_.ExponentialBackOffOption) (*http.Response, error) {
	var option []time_.ExponentialBackOffOption
	option = append(option, time_.WithExponentialBackOffOptionMaxElapsedCount(3))
	option = append(option, opts...)
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
		resp, err := Do(httpReq)
		if err == nil {
			return resp, nil
		}
		retries++
		if rewindableErr != nil {
			return nil, fmt.Errorf("http do cannot backoff: %w", retries, rewindableErr)
		}

		wait, ok := backoff.NextBackOff()
		if !ok {
			return nil, fmt.Errorf("http do reach backoff limit after retries %d", retries)
		}

		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
			continue
		case <-httpReq.Context().Done():
			timer.Stop()
			return nil, fmt.Errorf("http do canceled after retries %d: %w", retries, httpReq.Context().Err())
		}
	}
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

	httpResp, err := Do(httpReq)
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
func DoJsonWithBackoff(httpReq *http.Request, req, resp interface{}, opts ...time_.ExponentialBackOffOption) error {
	var option []time_.ExponentialBackOffOption
	option = append(option, time_.WithExponentialBackOffOptionMaxElapsedCount(3))
	option = append(option, opts...)
	backoff := time_.NewDefaultExponentialBackOff(option...)
	var retries int
	for {
		err := DoJson(httpReq, req, resp)
		if err == nil {
			return nil
		}
		retries++

		wait, ok := backoff.NextBackOff()
		if !ok {
			return fmt.Errorf("http do reach backoff limit after retries %d", retries)
		}

		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
			continue
		case <-httpReq.Context().Done():
			timer.Stop()
			return fmt.Errorf("http do canceled after retries %d: %w", retries, httpReq.Context().Err())
		}
	}
}
