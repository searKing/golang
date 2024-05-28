// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package restful

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

const timeout = 4 * time.Second

var defaultClient = &http.Client{
	Timeout: timeout,
}

func httpRequest(method string, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

func httpMethod(req *http.Request) (string, error) {
	resp, err := defaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		return "", errors.New(string(body))
	}
	return string(body), nil
}

func Get(url string) (string, error) {
	req, err := httpRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	return httpMethod(req)
}

func Post(url string, body []byte) (string, error) {
	req, err := httpRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	return httpMethod(req)
}

func Put(url string, body []byte) (string, error) {
	req, err := httpRequest(http.MethodPut, url, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	return httpMethod(req)
}
