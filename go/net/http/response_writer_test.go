// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseWriterUnwrap(t *testing.T) {
	testWriter := httptest.NewRecorder()
	writer := &recordResponseWriter{ResponseWriter: testWriter}
	if testWriter != writer.Unwrap() {
		t.Errorf("response writer not the same")
	}
}

func TestResponseWriterReset(t *testing.T) {
	testWriter := httptest.NewRecorder()
	w := &recordResponseWriter{}

	w.reset(testWriter)
	if w.size != 0 {
		t.Errorf("writer.size got (%v), want (%v)", w.size, 0)
	}
	if w.status != 0 {
		t.Errorf("writer.status got (%v), want (%v)", w.status, 0)
	}
	if w.ResponseWriter != testWriter {
		t.Errorf("response writer not the same")
	}
}

func TestResponseWriterWriteHeader(t *testing.T) {
	testWriter := httptest.NewRecorder()
	w := &recordResponseWriter{}
	w.reset(testWriter)

	w.WriteHeader(http.StatusMultipleChoices)
	if !w.Written() {
		t.Errorf("w.Written() got (%v), want (%t)", w.Written(), true)
	}
	if w.Status() != http.StatusMultipleChoices {
		t.Errorf("w.Status() got (%v), want (%v)", w.Status(), http.StatusMultipleChoices)
	}
	if testWriter.Code != http.StatusMultipleChoices {
		t.Errorf("testWriter.Code got (%v), want (%v)", testWriter.Code, http.StatusMultipleChoices)
	}
	w.WriteHeader(-1)
	if w.Status() != -1 {
		t.Errorf("w.Status() got (%v), want (%v)", w.Status(), -1)
	}
}

func TestResponseWriterWriteHeadersNow(t *testing.T) {
	testWriter := httptest.NewRecorder()
	w := &recordResponseWriter{}
	w.reset(testWriter)

	w.WriteHeader(http.StatusMultipleChoices)

	if !w.Written() {
		t.Errorf("w.Written() got (%v), want (%t)", w.Written(), true)
	}
	if w.Size() != 0 {
		t.Errorf("w.Size() got (%v), want (%v)", w.Size(), 0)
	}
	if testWriter.Code != http.StatusMultipleChoices {
		t.Errorf("testWriter.Code got (%v), want (%v)", testWriter.Code, http.StatusMultipleChoices)
	}

	w.size = 10
	if w.Size() != 10 {
		t.Errorf("w.Size() got (%v), want (%v)", w.Size(), 10)
	}
}

func TestResponseWriterWrite(t *testing.T) {
	testWriter := httptest.NewRecorder()
	w := &recordResponseWriter{}
	w.reset(testWriter)

	n, err := w.Write([]byte("hola"))
	if n != 4 {
		t.Errorf("n got (%v), want (%v)", n, 4)
	}
	if w.Size() != 4 {
		t.Errorf("w.Size() got (%v), want (%v)", w.Size(), 4)
	}
	if w.Status() != http.StatusOK {
		t.Errorf("w.Status() got (%v), want (%v)", w.Status(), http.StatusOK)
	}
	if testWriter.Code != http.StatusOK {
		t.Errorf("testWriter.Code got (%v), want (%v)", testWriter.Code, http.StatusOK)
	}
	if testWriter.Body.String() != "hola" {
		t.Errorf("testWriter.Body.String() got (%v), want (%v)", testWriter.Body.String(), "hola")
	}
	if err != nil {
		t.Errorf("w.Write() got an error: %s", err.Error())
	}

	n, err = w.Write([]byte(" adios"))
	if n != 6 {
		t.Errorf("n got (%v), want (%v)", n, 6)
	}
	if w.Size() != 10 {
		t.Errorf("w.Size() got (%v), want (%v)", w.Size(), 10)
	}
	if testWriter.Body.String() != "hola adios" {
		t.Errorf("testWriter.Body.String() got (%v), want (%v)", testWriter.Body.String(), "hola adios")
	}
	if err != nil {
		t.Errorf("w.Write() got an error: %s", err.Error())
	}
}

func TestResponseWriterHijack(t *testing.T) {
	testWriter := httptest.NewRecorder()
	w := &recordResponseWriter{}
	w.reset(testWriter)
	_, _, err := w.Hijack()
	if err == nil {
		t.Error("w.Hijack() should got an error")
	}
	if w.Written() {
		t.Errorf("w.Written() got (%v), want (%t)", w.Written(), false)
	}
	w.CloseNotify()

	w.Flush()
}

func TestResponseWriterFlush(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &recordResponseWriter{}
		writer.reset(w)

		writer.WriteHeader(http.StatusInternalServerError)
		writer.Flush()
	}))
	defer testServer.Close()

	// should return 500
	resp, err := http.Get(testServer.URL)
	if err != nil {
		t.Errorf("http.Get() got an error: %s", err.Error())
	}
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("resp.StatusCode got (%v), want (%v)", resp.StatusCode, http.StatusInternalServerError)
	}
}

func TestResponseWriterStatusCode(t *testing.T) {
	testWriter := httptest.NewRecorder()
	w := &recordResponseWriter{}
	w.reset(testWriter)

	w.WriteHeader(http.StatusOK)

	if w.Status() != http.StatusOK {
		t.Errorf("w.Status() got (%v), want (%v)", w.Status(), http.StatusOK)
	}
	if !w.Written() {
		t.Errorf("w.Written() got (%v), want (%t)", w.Written(), true)
	}

	w.WriteHeader(http.StatusUnauthorized)

	// status must be 200 although we tried to change it
	if w.Status() != http.StatusUnauthorized {
		t.Errorf("w.Status() got (%v), want (%v)", w.Status(), http.StatusUnauthorized)
	}
}
