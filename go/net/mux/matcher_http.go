// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"io"

	"github.com/searKing/golang/go/net/mux/internal/http"
)

// PRI * HTTP/2.0\r\n\r\n
// HTTP parses the first line or upto 4096 bytes of the request to see if
// the connection contains an HTTP request.
func HTTP() MatcherFunc {
	return func(_ io.Writer, r io.Reader) bool {
		req := http.ReadRequestLine(r)
		if req == nil {
			return false
		}
		return true
	}
}
