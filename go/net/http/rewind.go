package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	io_ "github.com/searKing/golang/go/io"
)

// RequestWithBodyRewindable returns a Request suitable for use with Redirect, like 307 redirect for PUT or POST.
// Only a nil GetBody in Request may be replace with a rewindable GetBody, which is a Body replayer.
// See: https://github.com/golang/go/issues/7912
// See also: https://go-review.googlesource.com/c/go/+/29852/13/src/net/http/client.go#391
func RequestWithBodyRewindable(req *http.Request) *http.Request {
	if req.Body == nil || req.Body == http.NoBody {
		// No copying needed.
		return req
	}

	// If the request body can be reset back to its original
	// state via the optional req.GetBody, do that.
	if req.GetBody != nil {
		return req
	}

	var body io.Reader = req.Body

	// NewRequest and NewRequestWithContext in net/http will handle
	// See: https://github.com/golang/go/blob/2117ea9737bc9cb2e30cb087b76a283f68768819/src/net/http/request.go#L873
	switch body.(type) {
	case *bytes.Buffer:
	case *bytes.Reader:
	case *strings.Reader:
		return req
	}

	// Body in Request will be closed before redirect automaticly, so io.Seeker can not be used.

	var replay bytes.Buffer
	// Use a replay reader to capture any body sent in case we have to replay it again
	replayR := io_.ReplayReader(req.Body)
	replayRC := replayReadCloser{Reader: replayR, Closer: req.Body}
	req.Body = replayRC
	req.GetBody = func() (io.ReadCloser, error) {
		replayR.Replay()

		// Refresh the body reader so the body can be sent again
		// take care of req.Body set to nil by caller outside
		if req.Body == nil {
			return nil, nil
		}
		return ioutil.NopCloser(&replay), nil
	}
	return req
}

type replayReadCloser struct {
	io.Reader
	io.Closer
}
