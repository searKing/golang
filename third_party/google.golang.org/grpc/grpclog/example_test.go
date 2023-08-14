// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpclog_test

import (
	grpclog_ "github.com/searKing/golang/third_party/google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/grpclog"
)

func ExampleDefaultSlogLogger() {
	grpclog.SetLoggerV2(grpclog_.DefaultSlogLogger())
	grpclog.Info("info")
	grpclog.Warning("warning")
	grpclog.Error("error")
	// 2023/08/14 15:39:54 INFO info
	// 2023/08/14 15:39:54 WARN warning
	// 2023/08/14 15:39:54 ERROR error
}
