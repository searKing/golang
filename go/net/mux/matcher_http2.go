// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"io"
	"strings"

	"golang.org/x/net/http2/hpack"

	http2_ "github.com/searKing/golang/go/net/mux/internal/http2"
)

// HTTP2 parses the frame header of the first frame to detect whether the
// connection is an HTTP2 connection.
func HTTP2() MatcherFunc {
	return func(_ io.Writer, r io.Reader) bool {
		return http2_.HasClientPreface(r)
	}
}

// HTTP2HeaderField returns a matcher matching the header fields of the first
// headers frame.
// writes the settings to the server if sendSetting is true.
// Prefer HTTP2HeaderField over this one, if the client does not block on receiving a SETTING frame.
func HTTP2HeaderField(sendSetting bool,
	match func(actual, expect map[string]hpack.HeaderField) bool,
	expects ...hpack.HeaderField) MatcherFunc {
	return func(w io.Writer, r io.Reader) bool {
		if !sendSetting {
			w = io.Discard
		}
		return http2_.MatchHTTP2Header(w, r, nil, func(parsedHeader map[string]hpack.HeaderField) bool {
			var expectMap = map[string]hpack.HeaderField{}
			for _, expect := range expects {
				expectMap[expect.Name] = expect
			}
			return match(parsedHeader, expectMap)
		})
	}
}

// helper functions

// HTTP2HeaderFieldValue returns a matcher matching the header fields, registered with the match handler.
func HTTP2HeaderFieldValue(sendSetting bool, match func(actualVal, expectVal string) bool, expects ...hpack.HeaderField) MatcherFunc {
	return HTTP2HeaderField(sendSetting, func(actualHeaderByName, expectHeaderByName map[string]hpack.HeaderField) bool {
		for name := range expectHeaderByName {
			if match(actualHeaderByName[name].Value, expectHeaderByName[name].Value) {
				return false
			}
		}
		return true
	}, expects...)
}

// HTTP2HeaderFieldEqual returns a matcher matching the header fields.
func HTTP2HeaderFieldEqual(sendSetting bool, headers ...hpack.HeaderField) MatcherFunc {
	return HTTP2HeaderFieldValue(sendSetting, func(actual string, expect string) bool {
		return actual == expect
	}, headers...)
}

// HTTP2HeaderFieldPrefix returns a matcher matching the header fields.
// If the header with key name has a
// value prefixed with valuePrefix, this will match.
func HTTP2HeaderFieldPrefix(sendSetting bool, headers ...hpack.HeaderField) MatcherFunc {
	return HTTP2HeaderFieldValue(sendSetting, strings.HasPrefix, headers...)
}
