// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package restful

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func HttpGetHandler(v any, w http.ResponseWriter, r *http.Request) {
	body, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		body = []byte("[]")
	}

	// response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.FormatInt(int64(len(body)), 10))
	w.Header().Set("Connection", "close")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func HttpPostHandler(v any, w http.ResponseWriter, r *http.Request) (finished bool) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		if r.ContentLength == -1 {
			http.NotFound(w, r)
			return true
		} else {
			ctype := strings.ToLower(r.Header.Get("Content-Type"))
			if ctype != "application/json" {
				http.NotFound(w, r)
				return true
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return true
			}

			err = json.Unmarshal(body, v)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				http.NotFound(w, r)
				return true
			}
			return false
		}
	} else {
		http.NotFound(w, r)
		return true
	}
}
