// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	io_ "github.com/searKing/golang/go/io"
)

var (
	ErrBodyNotRewindable = errors.New("body not rewindable")
)
var nopCloserType = reflect.TypeOf(io.NopCloser(nil))

func RequestRewindableWithFileName(name string) (body io.ReadCloser, getBody func() (io.ReadCloser, error), err error) {
	return BodyRewindableWithFilePosition(name, 0, io.SeekCurrent)
}

func BodyRewindableWithFilePosition(name string, offset int64, whence int) (body io.ReadCloser, getBody func() (io.ReadCloser, error), err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}

	if _, err = file.Seek(offset, whence); err != nil {
		return nil, nil, err
	}

	return BodyRewindableWithFile(file)
}

// BodyRewindableWithFile returns a Request suitable for use with Redirect, like 307 redirect for PUT or POST.
// Only a nil GetBody in Request may be replace with a rewindable GetBody, which is a Body *os.File.
// See: https://github.com/golang/go/issues/7912
// See also: https://go-review.googlesource.com/c/go/+/29852/13/src/net/http/client.go#391
func BodyRewindableWithFile(file *os.File) (body io.ReadCloser, getBody func() (io.ReadCloser, error), err error) {
	offset, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, nil, err
	}

	return file, func() (io.ReadCloser, error) {
		file_, err_ := os.Open(file.Name())
		if err_ != nil {
			return nil, err
		}

		if _, err_ = file_.Seek(offset, io.SeekStart); err != nil {
			return nil, err_
		}
		return file_, err_
	}, nil
}

// RequestWithBodyRewindable returns a Request suitable for use with Redirect, like 307 redirect for PUT or POST.
// Only a nil GetBody in Request may be replaced with a rewindable GetBody, which is a Body replayer.
// A body with a type not ioutil.NopCloser(nil) may return error as the Body in Request will be closed before redirect automatically.
// So you can close body by yourself to ensure rewindable always:
// Examples:
// 	body := req.Body
// 	defer body.Close() // body will not be closed inside
// 	req.Body = ioutil.NopCloser(body)
// 	_ = RequestWithBodyRewindable(req)
// // do http requests...
//
// See: https://github.com/golang/go/issues/7912
// See also: https://go-review.googlesource.com/c/go/+/29852/13/src/net/http/client.go#391
func RequestWithBodyRewindable(req *http.Request) error {
	if req.Body == nil || req.Body == http.NoBody {
		// No copying needed.
		return nil
	}

	// If the request body can be reset back to its original
	// state via the optional req.GetBody, do that.
	if req.GetBody != nil {
		return nil
	}

	body, neverClose := isNeverCloseReader(req.Body)
	if !neverClose {
		// Body in Request will be closed before redirect automatically, so io.Seeker can not be used.

		// take care of *os.File, which can be reopen
		switch body_ := body.(type) {
		case *os.File:
			_, getBody, err := BodyRewindableWithFile(body_)
			if err != nil {
				return err
			}
			req.GetBody = getBody
			return nil
		}
		return ErrBodyNotRewindable
	}

	// handle never closed body

	// NewRequest and NewRequestWithContext in net/http will handle
	// See: https://github.com/golang/go/blob/2117ea9737bc9cb2e30cb087b76a283f68768819/src/net/http/request.go#L873
	switch v := body.(type) {
	case *bytes.Buffer:
		req.ContentLength = int64(v.Len())
		buf := v.Bytes()
		req.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(buf)
			return io.NopCloser(r), nil
		}
		return nil
	case *bytes.Reader:
		req.ContentLength = int64(v.Len())
		snapshot := *v
		req.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
		return nil
	case *strings.Reader:
		req.ContentLength = int64(v.Len())
		snapshot := *v
		req.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
		return nil
	case *os.File:
		_, getBody, err := BodyRewindableWithFile(v)
		if err != nil {
			return err
		}
		req.GetBody = getBody
		return nil
	}

	// Handle unknown types

	// Use io.Seeker to rewind if Seek succeed
	{
		if seeker, ok := body.(io.Seeker); ok {
			offset, err := seeker.Seek(0, io.SeekCurrent)
			if err == nil {
				req.GetBody = func() (io.ReadCloser, error) {
					_, err := seeker.Seek(offset, io.SeekStart)
					if err != nil {
						return nil, err
					}
					return req.Body, nil
				}
			}
		}
	}

	// Use a replay reader to capture any body sent in case we have to replay it again
	// All data will be buffered in memory in case of reread.
	{
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
			return ioutil.NopCloser(replayR), nil
		}
	}
	return nil
}

type replayReadCloser struct {
	io.Reader
	io.Closer
}

// isNeverCloseReader reports whether r is a type known to not closed.
// Its caller uses this as an optional optimization to
// send fewer TCP packets.
func isNeverCloseReader(r io.Reader) (rr io.Reader, nopClose bool) {
	return _isNeverCloseReader(r, false)
}

func _isNeverCloseReader(r io.Reader, nopclose bool) (rr io.Reader, nopClose bool) {
	switch r.(type) {
	case *bytes.Reader, *bytes.Buffer, *strings.Reader:
		return r, true
	}
	if reflect.TypeOf(r) == nopCloserType {
		return _isNeverCloseReader(reflect.ValueOf(r).Field(0).Interface().(io.Reader), true)
	}
	return r, nopclose
}
