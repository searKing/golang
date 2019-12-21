package cmux

import (
	"io"
	"net/http"
	"strings"

	http_util "github.com/searKing/golang/go/net/cmux/internal/http"
	http_ "github.com/searKing/golang/go/net/http"
)

// HTTP1Fast only matches the methods in the HTTP request.
//
// This matcher is very optimistic: if it returns true, it does not mean that
// the request is a valid HTTP response. If you want a correct but slower HTTP1
// matcher, use HTTP1 instead.
func HTTP1Fast(extMethods ...string) MatcherFunc {
	return PrefixMatcher(append(http_.Methods, extMethods...)...)
}

// HTTP1 parses the first line or upto 4096 bytes of the request to see if
// the conection contains an HTTP request.
func HTTP1() MatcherFunc {
	return func(w io.Writer, r io.Reader) bool {
		req := http_util.ReadRequestLine(r)
		if req == nil {
			return false
		}
		return req.ProtoMajor == 1
	}
}

func HTTP1Header(match func(actual, expect http.Header) bool, expect http.Header) MatcherFunc {
	return func(w io.Writer, r io.Reader) bool {
		return http_util.MatchHTTPHeader(r, func(parsedHeader http.Header) bool {
			return match(parsedHeader, expect)
		})
	}
}

// helper functions
func HTTP1HeaderValue(match func(actual, expect string) bool, expect http.Header) MatcherFunc {
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
func HTTP1HeaderEqual(header http.Header) MatcherFunc {
	return HTTP1HeaderValue(func(actual string, expect string) bool {
		return actual == expect
	}, header)
}

// HTTP1HeaderPrefix returns a matcher matching the header fields of the
// first request of an HTTP 1 connection. If the header with key name has a
// value prefixed with valuePrefix, this will match.
func HTTP1HeaderPrefix(header http.Header) MatcherFunc {
	return HTTP1HeaderValue(strings.HasPrefix, header)
}
