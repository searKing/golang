package http2

import (
	"golang.org/x/net/http2"
	"io"
	"io/ioutil"
)

func HasClientPreface(w io.Writer, r io.Reader) bool {
	r = io.LimitReader(r, int64(len(http2.ClientPreface)))

	clientPrefaceLine, err := ioutil.ReadAll(r)
	if err != nil {
		return false
	}
	return string(clientPrefaceLine) == http2.ClientPreface
}
