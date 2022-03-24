package status

import (
	"github.com/golang/protobuf/proto"
	"github.com/searKing/golang/go/error/cause"
)

// WithDetails returns a new status with the provided details messages appended to the status.
// If any errors are encountered, it returns original status.
// WithDetails does not change original code always.
func WithDetails(s *Status, details ...proto.Message) *Status {
	stat, err := s.WithDetails(details...)
	if err == nil {
		return stat
	}
	if s.Err() == nil {
		return s
	}

	return Newe(s.Code(), cause.WithError(s.Err(), err))
}
