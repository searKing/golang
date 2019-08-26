package http2

import (
	"bytes"
	"golang.org/x/net/http2"
	"io"
)

var (
	clientPreface = []byte(http2.ClientPreface)
)

func HasClientPreface(w io.Writer, r io.Reader) bool {
	// Check the validity of client preface.
	preface := make([]byte, len(clientPreface))
	if _, err := io.ReadFull(r, preface); err != nil {
		return false
	}
	return bytes.Equal(preface, clientPreface)
}
