// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build !go1.13

package otelgrpc

import (
	"google.golang.org/grpc/status"
)

func unwrapNativeWrappedGRPCStatus(err error) (*status.Status, bool) {
	return nil, false
}
