// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"strings"

	"golang.org/x/net/http2/hpack"
)

// HTTP2 parses the frame header of the first frame to detect whether the
// connection is an HTTP2 connection.
func GRPC() MatcherFunc {
	return HTTP2HeaderFieldValue(false, strings.EqualFold, hpack.HeaderField{
		Name:  "Content-Type",
		Value: "application/grpc",
	})
}
