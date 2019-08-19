package status

import (
	"github.com/golang/protobuf/proto"
	"github.com/searKing/golang/go/error/cause"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Newe returns a Status representing c and error, that is .
// if err == nil, return codes.OK
// if err != nil but c == codes.OK, return codes.Internal
// otherwise, return c
func Newe(c codes.Code, err error, details ...proto.Message) *status.Status {
	if err == nil {
		return status.New(codes.OK, "")
	}
	if c == codes.OK {
		// no error details for status with code OK
		c = codes.Internal
	}
	return WithDetails(status.New(c, err.Error()), details...)
}

// Errore returns an error representing c and error.  If err is nil, returns nil.
func Errore(c codes.Code, err error, details ...proto.Message) error {
	return Newe(c, err, details...).Err()
}

// FromError returns a Status representing err if it was produced from this
// package or has a method `GRPCStatus() *Status`. Otherwise, ok is false and a
// Status is returned with code.Code and the original error message.
func FromError(c codes.Code, err error) (s *status.Status, ok bool) {
	stat, ok := status.FromError(err)
	if ok {
		return stat, ok
	}
	return status.New(c, err.Error()), false
}

// Convert is a convenience function which removes the need to handle the
// boolean return value from FromError.
func Convert(c codes.Code, err error) *status.Status {
	s, _ := FromError(c, err)
	return s
}

// WithBadRequestDetails returns a new status with the provided bad requests messages appended to the status.
func WithBadRequestDetails(s *status.Status, fields ...errdetails.BadRequest_FieldViolation) *status.Status {
	var badRequest errdetails.BadRequest
	for _, f := range fields {
		badRequest.FieldViolations = append(badRequest.FieldViolations, &f)
	}
	return WithDetails(s, &badRequest)
}

// WithDetails returns a new status with the provided details messages appended to the status.
// If any errors are encountered, it returns original status.
// WithDetails does not change original code always.
func WithDetails(s *status.Status, details ...proto.Message) *status.Status {
	stat, err := s.WithDetails(details...)
	if err == nil {
		return stat
	}
	if s.Err() == nil {
		return s
	}

	Newe(s.Code(), cause.WithError(s.Err(), err))
	return s
}
