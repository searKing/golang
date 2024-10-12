// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"
	"testing"
	"time"
)

type fakeResponseWriter struct {
	flushErrorCalled       bool
	setWriteDeadlineCalled time.Time
	setReadDeadlineCalled  time.Time
}

func (rw *fakeResponseWriter) Header() http.Header {
	return nil
}

func (rw *fakeResponseWriter) Write(p []byte) (int, error) {
	return 0, nil
}

func (rw *fakeResponseWriter) WriteHeader(statusCode int) {
}

func (rw *fakeResponseWriter) FlushError() error {
	rw.flushErrorCalled = true
	return nil
}

func (rw *fakeResponseWriter) SetWriteDeadline(deadline time.Time) error {
	rw.setWriteDeadlineCalled = deadline
	return nil
}

func (rw *fakeResponseWriter) SetReadDeadline(deadline time.Time) error {
	rw.setReadDeadlineCalled = deadline
	return nil
}

func TestResponseWriterDelegatorUnwrap(t *testing.T) {
	w := &fakeResponseWriter{}
	rwd := &responseWriterDelegator{ResponseWriter: w}

	if rwd.Unwrap() != w {
		t.Error("unwrapped responsewriter must equal to the original responsewriter")
	}

	controller := http.NewResponseController(rwd)
	if err := controller.Flush(); err != nil || !w.flushErrorCalled {
		t.Error("FlushError must be propagated to the original responsewriter")
	}

	timeNow := time.Now()
	if err := controller.SetWriteDeadline(timeNow); err != nil || w.setWriteDeadlineCalled != timeNow {
		t.Error("SetWriteDeadline must be propagated to the original responsewriter")
	}

	if err := controller.SetReadDeadline(timeNow); err != nil || w.setReadDeadlineCalled != timeNow {
		t.Error("SetReadDeadline must be propagated to the original responsewriter")
	}
}
