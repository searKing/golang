package http

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

const maxHTTPRead = 4096

// read first line of HTTP request
func ReadRequestLine(r io.Reader) *http.Request {
	br := bufio.NewReader(&io.LimitedReader{R: r, N: maxHTTPRead})
	l, part, err := br.ReadLine()
	if err != nil || part {
		return nil
	}
	// padding with http header tailer bytes \r\n\r\n
	l = append(l, []byte("\r\n\r\n")...)

	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(l)))
	if err != nil {
		return nil
	}
	return req
}

func MatchHTTPHeader(r io.Reader, matches func(parsedHeader http.Header) bool) (matched bool) {
	req, err := http.ReadRequest(bufio.NewReader(r))
	if err != nil {
		return false
	}

	return matches(req.Header)
}
