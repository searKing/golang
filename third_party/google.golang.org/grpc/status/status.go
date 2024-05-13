// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package status

import (
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
	if _, ok := status.FromError(err); ok {
		return err
	}

	return errors_.Mark(status.New(c, msg).Err(), err)
}

// Errorfe returns Errore(c, fmt.Sprintf(format, a...)).
func Errorfe(err error, c codes.Code, format string, a ...any) error {
	return Errore(err, c, fmt.Sprintf(format, a...))
}

// FromError returns a Status representation of err.
//
//   - If err was produced by this package or implements the method `GRPCStatus()
//     *Status` and `GRPCStatus()` does not return nil, or if err wraps a type
//     satisfying this, the Status from `GRPCStatus()` is returned.  For wrapped
//     errors, the message returned contains the entire err.Error() text and not
//     just the wrapped status. In that case, ok is true.
//
//   - If err is nil, a Status is returned with codes.OK and no message, and ok
//     is true.
//
//   - If err implements the method `GRPCStatus() *Status` and `GRPCStatus()`
//     returns nil (which maps to Codes.OK), or if err wraps a type
//     satisfying this, a Status is returned with codes.Unknown and err's
//     Error() message, and ok is false.
//
//   - Otherwise, err is an error not compatible with this package.  In this
//     case, a Status is returned with codes.Unknown and err's Error() message,
//     and ok is false.
//
// Deprecated: As of grpc 1.58.0, Use status.FromError instead.
func FromError(err error) (s *Status, ok bool) {
	return status.FromError(err)
}

// Convert is a convenience function which removes the need to handle the
// boolean return value from FromError.
// Deprecated: As of grpc 1.58.0, Use status.Convert instead.
func Convert(err error) *Status {
	return status.Convert(err)
}
