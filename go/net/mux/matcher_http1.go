// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"io"
	"net/http"
	"strings"

	http_ "github.com/searKing/golang/go/net/http"
	httputil "github.com/searKing/golang/go/net/mux/internal/http"
)

// HTTP1Fast only matches the methods in the HTTP request.
//
// This matcher is very optimistic: if it returns true, it does not mean that
// the request is a valid HTTP response. If you want a correct but slower HTTP1
// matcher, use HTTP1 instead.
func HTTP1Fast(extMethods ...string) MatcherFunc {
	return AnyPrefixMatcher(append(http_.Methods, extMethods...)...)
}

// HTTP1 parses the first line or upto 4096 bytes of the request to see if
// the conection contains an HTTP request.
func HTTP1() MatcherFunc {
	return func(_ io.Writer, r io.Reader) bool {
		req := httputil.ReadRequestLine(r)
		if req == nil {
			return false
		}
		return req.ProtoMajor == 1
	}
}

// HTTP1Header returns true if all headers are expected
func HTTP1Header(match func(actual, expect http.Header) bool, expect http.Header) MatcherFunc {
	return func(_ io.Writer, r io.Reader) bool {
		return httputil.MatchHTTPHeader(r, func(parsedHeader http.Header) bool {
			return match(parsedHeader, expect)
		})
	}
}

// helper functions

// HTTP1HeaderValue returns true if all headers are expected, shorthand for HTTP1Header
// strings.Match for all value in expects
func HTTP1HeaderValue(match func(actualVal, expectVal string) bool, expect http.Header) MatcherFunc {
	return HTTP1Header(func(actual, expect http.Header) bool {
		for name := range expect {
			if match(actual.Get(name), expect.Get(name)) {
				return false
			}
		}
		return true
	}, expect)
}

// HTTP1HeaderEqual returns a matcher matching the header fields of the first
// request of an HTTP 1 connection.
// strings.Equal for all value in expects
func HTTP1HeaderEqual(header http.Header) MatcherFunc {
	return HTTP1HeaderValue(func(actual string, expect string) bool {
		return actual == expect
	}, header)
}

// HTTP1HeaderPrefix returns a matcher matching the header fields of the
// first request of an HTTP 1 connection. If the header with key name has a
// value prefixed with valuePrefix, this will match.
// strings.HasPrefix for all value in expects
func HTTP1HeaderPrefix(header http.Header) MatcherFunc {
	return HTTP1HeaderValue(strings.HasPrefix, header)
}
