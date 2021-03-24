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
	"time"

	time_ "github.com/searKing/golang/go/time"
)

func ReplaceHttpRequestBody(req *http.Request, body io.Reader) {
	if req.Body != nil {
		req.Body.Close()
	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}
	req.Body = rc
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
	backoff := time_.NewExponentialBackOff(option...)
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
	backoff := time_.NewExponentialBackOff(option...)
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
