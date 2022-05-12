// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package status

import (
	"errors"
	"fmt"

	errors_ "github.com/searKing/golang/go/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Status = status.Status

// Errore returns an error representing c and msg if err is not status.Status.
// If c is OK, returns nil.
// If err is status.Status already, return err.
// else, return an error representing c and msg marked with err.
func Errore(err error, c codes.Code, msg string) error {
	if err == nil || c == codes.OK {
		return nil
	}
	if _, ok := FromError(err); ok {
		return err
	}

	return errors_.Mark(status.New(c, msg).Err(), err)
}

// Errorfe returns Errore(c, fmt.Sprintf(format, a...)).
func Errorfe(err error, c codes.Code, format string, a ...interface{}) error {
	return Errore(err, c, fmt.Sprintf(format, a...))
}

// FromError returns a Status representation of err.
//
// - If err was produced by this package or any err by `errors.UnWrap()` implements the method
//   `GRPCStatus() *Status`, the appropriate Status is returned.
//
// - If err is nil, a Status is returned with codes.OK and no message.
//
// - Otherwise, err is an error not compatible with this package.  In this
//   case, a Status is returned with codes.Unknown and err's Error() message,
//   and ok is false.
func FromError(err error) (s *Status, ok bool) {
	s, ok = status.FromError(err)
	if ok {
		return s, ok
	}
	var gRPCStatus interface {
		GRPCStatus() *status.Status
	}
	if errors.As(err, &gRPCStatus) {
		s = gRPCStatus.GRPCStatus()
		ok = true
	}
	return s, ok
}

// Convert is a convenience function which removes the need to handle the
// boolean return value from FromError.
func Convert(err error) *Status {
	s, _ := FromError(err)
	return s
}
