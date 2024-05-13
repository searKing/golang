// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package status_test

import (
	"errors"
	"fmt"
	"testing"

	errors_ "github.com/searKing/golang/go/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	status_ "github.com/searKing/golang/third_party/google.golang.org/grpc/status"
)

func TestFromError(t *testing.T) {
	var markErr = errors.New("mark")
	code, message := codes.Internal, "test description"
	err := errors_.Mark(status.Error(code, message), markErr)

	s, ok := status.FromError(err)
	if !ok || s.Code() != code || s.Message() != fmt.Sprintf("rpc error: code = Internal desc = %s", message) || s.Err() == nil {
		t.Fatalf("FromError(%v) = %v, %v; want <Code()=%s, Message()=%q, Err()!=nil>, true", err, s, ok, code, message)
	}

	var normalErr = errors.New("normal")
	err = errors_.Mark(normalErr, status.Error(code, message))
	s, ok = status.FromError(err)
	if !ok || s.Code() != code || s.Message() != "normal" || s.Err() == nil {
		t.Fatalf("FromError(%v) = %v, %v; want <Code()=%s, Message()=%q, Err()!=nil>, true", err, s, ok, code, message)
	}
}

func TestConvertKnownError(t *testing.T) {
	var markErr = errors.New("mark")
	code, message := codes.Internal, "test description"
	err := errors_.Mark(status.Error(code, message), markErr)
	s := status.Convert(err)
	if s.Code() != code || s.Message() != fmt.Sprintf("rpc error: code = Internal desc = %s", message) {
		t.Fatalf("Convert(%v) = %v; want <Code()=%s, Message()=%q>", err, s, code, message)
	}
}

func TestConvertUnknownError(t *testing.T) {
	code, message := codes.Unknown, "unknown error"
	err := errors.New("unknown error")
	s := status.Convert(err)
	if s.Code() != code || s.Message() != message {
		t.Fatalf("Convert(%v) = %v; want <Code()=%s, Message()=%q>", err, s, code, message)
	}
}

func TestError(t *testing.T) {
	var errorStringErr = errors.New("errorString")
	code, message := codes.Internal, "test description"
	err := status_.Errore(errorStringErr, code, message)
	if got, want := err.Error(), fmt.Sprintf("rpc error: code = Internal desc = %s", message); got != want {
		t.Fatalf("err.Error() = %q; want %q", got, want)
	}
	s, _ := status.FromError(err)
	if got, want := s.Code(), code; got != want {
		t.Fatalf("err.Code() = %s; want %s", got, want)
	}
	if got, want := s.Message(), fmt.Sprintf("rpc error: code = Internal desc = %s", message); got != want {
		t.Fatalf("err.Message() = %s; want %s", got, want)
	}
}

func TestErrorOK(t *testing.T) {
	var errorStringErr = errors.New("errorString")
	err := status_.Errore(errorStringErr, codes.OK, "foo")
	if err != nil {
		t.Fatalf("Errore(codes.OK, _) = %p; want nil", err)
	}
}
