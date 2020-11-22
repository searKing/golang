package status

import "google.golang.org/genproto/googleapis/rpc/errdetails"

// WithBadRequestDetails returns a new status with the provided bad requests messages appended to the status.
func WithBadRequestDetails(s *Status, descByField map[string]string) *Status {
	if len(descByField) == 0 {
		return s
	}
	var badRequest errdetails.BadRequest
	for f, desc := range descByField {
		badRequest.FieldViolations = append(badRequest.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       f,
			Description: desc,
		})
	}
	return WithDetails(s, &badRequest)
}
