// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package validator_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	validator_ "github.com/searKing/golang/pkg/webserver/pkg/validator"
)

type request struct {
	Rate int `json:"rate%d" validate:"required"`
}

// newValidate returns a validator naming fields by json tags, so that messages of errors contain "%d".
func newValidate() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(f reflect.StructField) string { return f.Tag.Get("json") })
	return v
}

// checkInvalidArgument checks err is InvalidArgument with the message of the validation error kept as is.
func checkInvalidArgument(t *testing.T, err error) {
	t.Helper()
	want := newValidate().Struct(&request{}).Error()
	s, _ := status.FromError(err)
	if s.Code() != codes.InvalidArgument || s.Message() != want {
		t.Errorf("error = %v %q, want %v %q", s.Code(), s.Message(), codes.InvalidArgument, want)
	}
}

func TestUnaryServerInterceptor(t *testing.T) {
	interceptor := validator_.UnaryServerInterceptor(newValidate())
	_, err := interceptor(context.Background(), &request{}, &grpc.UnaryServerInfo{},
		func(ctx context.Context, req any) (any, error) {
			t.Fatal("handler called with invalid request")
			return nil, nil
		})
	checkInvalidArgument(t, err)

	resp, err := interceptor(context.Background(), &request{Rate: 1}, &grpc.UnaryServerInfo{},
		func(ctx context.Context, req any) (any, error) { return "ok", nil })
	if err != nil || resp != "ok" {
		t.Errorf("valid request = %v, %v, want ok, nil", resp, err)
	}
}

type serverStream struct {
	grpc.ServerStream
	req request
}

func (s *serverStream) Context() context.Context { return context.Background() }

func (s *serverStream) RecvMsg(m any) error {
	*m.(*request) = s.req
	return nil
}

func TestStreamServerInterceptor(t *testing.T) {
	interceptor := validator_.StreamServerInterceptor(newValidate())
	recv := func(ss *serverStream) error {
		return interceptor(nil, ss, &grpc.StreamServerInfo{}, func(srv any, stream grpc.ServerStream) error {
			return stream.RecvMsg(&request{})
		})
	}
	checkInvalidArgument(t, recv(&serverStream{}))
	if err := recv(&serverStream{req: request{Rate: 1}}); err != nil {
		t.Errorf("valid request = %v, want nil", err)
	}
}
