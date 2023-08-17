// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

var (
	_ http.ResponseWriter = (*recordResponseWriter)(nil)
	_ http.ResponseWriter = (*recordResponseWriter)(nil)
	_ http.Hijacker       = (*recordResponseWriter)(nil)
	_ http.Flusher        = (*recordResponseWriter)(nil)
	_ http.CloseNotifier  = (*recordResponseWriter)(nil)
)

// NewRecordResponseWriter creates a ResponseWriter that is a wrapper around [http.ResponseWriter] that
// provides extra information about the response.
// It is recommended that middleware handlers use this construct to wrap a [http.ResponseWriter]
// if the functionality calls for it.
func NewRecordResponseWriter(rw http.ResponseWriter) *recordResponseWriter {
	return &recordResponseWriter{
		ResponseWriter: rw,
	}
}

type recordResponseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func (w *recordResponseWriter) Context() context.Context {
	return w.Context()
}

func (w *recordResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func (w *recordResponseWriter) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.size = 0
	w.status = 0
}

func (w *recordResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// ExplicitlyWriteHeader forces to write the http header (status code + headers).
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes or 1xx informational responses.
func (w *recordResponseWriter) ExplicitlyWriteHeader() {
	if !w.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		w.WriteHeader(http.StatusOK)
	}
}

func (w *recordResponseWriter) Write(data []byte) (n int, err error) {
	w.ExplicitlyWriteHeader()
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

// WriteString writes the string into the response body.
func (w *recordResponseWriter) WriteString(s string) (n int, err error) {
	w.ExplicitlyWriteHeader()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += n
	return
}

// Status returns the status code of the response or 0 if the response has
// not been written
func (w *recordResponseWriter) Status() int {
	return w.status
}

// Size returns the number of bytes already written into the response http body.
// See Written()
func (w *recordResponseWriter) Size() int {
	return w.size
}

// Written returns whether or not the ResponseWriter has been written.
func (w *recordResponseWriter) Written() bool {
	return w.status != 0
}

// Hijack implements the http.Hijacker interface.
func (w *recordResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()

	}
	return nil, nil, errors.New("ResponseWriter doesn't support Hijacker interface")
}

// CloseNotify implements the http.CloseNotifier interface.
func (w *recordResponseWriter) CloseNotify() <-chan bool {
	if gone, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return gone.CloseNotify()
	}
	return nil
}

// Flush implements the http.Flusher interface.
func (w *recordResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		w.ExplicitlyWriteHeader()
		flusher.Flush()
	}
}

// Pusher get the http.Pusher for server push
func (w *recordResponseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}

func (w *recordResponseWriter) Push(target string, opts *http.PushOptions) error {
	pusher, ok := w.ResponseWriter.(http.Pusher)
	if ok {
		return pusher.Push(target, opts)
	}
	return fmt.Errorf("the ResponseWriter doesn't support the Pusher interface")
}
