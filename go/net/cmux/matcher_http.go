package cmux

import (
	"github.com/searKing/golang/go/net/cmux/internal/http"
	"io"
)

// PRI * HTTP/2.0\r\n\r\n
// HTTP parses the first line or upto 4096 bytes of the request to see if
// the conection contains an HTTP request.
func HTTP() MatcherFunc {
	return func(w io.Writer, r io.Reader) bool {
		req := http.ReadRequestLine(r)
		if req == nil {
			return false
		}
		return true
	}
}
