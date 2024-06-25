// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	time_ "github.com/searKing/golang/go/time"
)

// DoWithBackoff will retry by exponential backoff if failed.
// If request is not rewindable, retry wil be skipped.
func DoWithBackoff(httpReq *http.Request, opts ...DoWithBackoffOption) (resp *http.Response, err error) {
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
	var errs []error
	defer func() {
		if resp != nil {
			if err := errors.Join(errs...); err != nil {
				if resp.Header == nil {
					resp.Header = make(http.Header)
				}
				resp.Header.Add("Warning", Warn{
					Warn:     err.Error(),
					WarnCode: WarnMiscellaneousWarning,
				}.String())
			}
		}
	}()
	for {
		if retries > 0 && httpReq.GetBody != nil {
			newBody, err := httpReq.GetBody()
			if err != nil {
				errs = append(errs, err)
				return nil, errors.Join(errs...)
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
		resp, err = httpDo(httpReq, retries)
		errs = append(errs, err)

		wait, ok := backoff.NextBackOff()
		if !ok {
			if err != nil {
				return nil, fmt.Errorf("http do reach backoff limit after retries %d: %w", retries, errors.Join(errs...))
			} else {
				return resp, nil
			}
		}

		wait, retry := opt.RetryAfter(resp, err, wait)
		if !retry {
			if err != nil {
				return nil, fmt.Errorf("http do reach server limit after retries %d: %w", retries, errors.Join(errs...))
			} else {
				return resp, nil
			}
		}

		if rewindableErr != nil {
			if err != nil {
				return nil, fmt.Errorf("http do cannot rewindbody after retries %d: %w", retries, errors.Join(errs...))
			} else {
				resp.Header.Add("Warning", Warn{
					Warn:     errors.Join(errs...).Error(),
					WarnCode: WarnMiscellaneousWarning,
				}.String())
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
				return nil, fmt.Errorf("http do canceled after retries %d: %w", retries, errors.Join(errs...))
			} else {
				return resp, nil
			}
		}
	}
}

func HeadWithBackoff(ctx context.Context, url string, opts ...DoWithBackoffOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return DoWithBackoff(req, opts...)
}

func GetWithBackoff(ctx context.Context, url string, opts ...DoWithBackoffOption) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return DoWithBackoff(req, opts...)
}

func PostWithBackoff(ctx context.Context, url, contentType string, body io.Reader, opts ...DoWithBackoffOption) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return DoWithBackoff(req, opts...)
}

func PostFormWithBackoff(ctx context.Context, url string, data url.Values, opts ...DoWithBackoffOption) (resp *http.Response, err error) {
	return PostWithBackoff(ctx, url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()), opts...)
}

func PutWithBackoff(ctx context.Context, url, contentType string, body io.Reader, opts ...DoWithBackoffOption) (resp *http.Response, err error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return DoWithBackoff(req, opts...)
}

// DoJson the same as HttpDo, but bind with json
func DoJson(httpReq *http.Request, req, resp any) error {
	if req != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return err
		}
		reqBody := bytes.NewReader(data)
		httpReq.Header.Set("Content-Type", "application/json")
		ReplaceHttpRequestBody(httpReq, reqBody)
	}

	httpResp, err := DefaultClientDoRetryHandler(httpReq, 0)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	if resp == nil {
		return nil
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, resp)
}

// DoJsonWithBackoff the same as DoWithBackoff, but bind with json
func DoJsonWithBackoff(httpReq *http.Request, req, resp any, opts ...DoWithBackoffOption) error {
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

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, resp)
}
