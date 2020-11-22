// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package status

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Status = status.Status

// Newe returns a Status representing c and error, that is .
// if err == nil, return codes.OK
// if err != nil but c == codes.OK, return codes.Internal
// otherwise, return c
func Newe(c codes.Code, err error, details ...proto.Message) *Status {
	if err == nil {
		return status.New(codes.OK, "")
	}

	if c == codes.OK {
		// no error details for status with code OK
		c = codes.Internal
	}

	return Convert(c, err, details...)
}

// Errore returns an error representing c and error.  If err is nil, returns nil.
func Errore(c codes.Code, err error, details ...proto.Message) error {
	return Newe(c, err, details...).Err()
}

// FromError returns a Status representing err if it was produced from this
// package or has a method `GRPCStatus() *Status`. Otherwise, ok is false and a
// Status is returned with code.Code and the original error message.
// code is set only if err has not implemented interface {
//		GRPCStatus() *Status
//	}
func FromError(c codes.Code, err error, details ...proto.Message) (s *Status, ok bool) {
	stat, ok := status.FromError(err)
	if ok {
		return WithDetails(stat, details...), ok
	}
	return WithDetails(status.New(c, err.Error()), details...), false
}

// Convert is a convenience function which removes the need to handle the
// boolean return value from FromError.
func Convert(c codes.Code, err error, details ...proto.Message) *Status {
	s, _ := FromError(c, err, details...)
	return s
}
