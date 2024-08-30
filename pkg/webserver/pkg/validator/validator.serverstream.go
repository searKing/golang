// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package validator

import (
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ grpc.ServerStream = (*validatorServerStream)(nil)

// validatorServerStream wraps grpc.ServerStream allowing each Sent/Recv of message.
type validatorServerStream struct {
	// validator contains the validator settings and cache
	validator *validator.Validate

	grpc.ServerStream
}

func (s *validatorServerStream) RecvMsg(req interface{}) error {
	err := s.ServerStream.RecvMsg(req)
	if err != nil {
		return err
	}
	if v := s.validator; v != nil {
		if err = v.StructCtx(s.Context(), req); err != nil {
			return status.Errorf(codes.InvalidArgument, err.Error())
		}
	}
	return nil
}

// DON'T CHECK RESPONSE
//func (s *validatorServerStream) SendMsg(resp interface{}) error {
//	if s.validator != nil {
//		if err := s.validator.Struct(resp); err != nil {
//			return status.Errorf(codes.Internal, err.Error())
//		}
//	}
//	return s.ServerStream.SendMsg(resp)
//}
