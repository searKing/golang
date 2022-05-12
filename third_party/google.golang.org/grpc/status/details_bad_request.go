package status

import "google.golang.org/genproto/googleapis/rpc/errdetails"

// WithBadRequestDetails returns a new status with the provided bad requests messages appended to the status.
func WithBadRequestDetails(s *Status, descByField map[string]string) (*Status, error) {
	if s == nil {
		return nil, nil
	}
	if len(descByField) == 0 {
		return s, nil
	}
	var badRequest errdetails.BadRequest
	for f, desc := range descByField {
		badRequest.FieldViolations = append(badRequest.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       f,
			Description: desc,
		})
	}
	return s.WithDetails(&badRequest)
}
