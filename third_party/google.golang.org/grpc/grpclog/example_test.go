// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpclog_test

import (
	"fmt"
	"os"

	grpclog_ "github.com/searKing/golang/third_party/google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/grpclog"
)

func ExampleDefaultSlogLogger() {
	err := os.Setenv("GRPC_GO_LOG_SEVERITY_LEVEL", "INFO")
	if err != nil {
		fmt.Printf("set env %q failed: %s", "GRPC_GO_LOG_SEVERITY_LEVEL", err)
	}
	grpclog.SetLoggerV2(grpclog_.DefaultSlogLogger())
	grpclog.Info("info")
	grpclog.Warning("warning")
	grpclog.Error("error")

	// 2023/08/14 15:39:54 INFO info
	// 2023/08/14 15:39:54 WARN warning
	// 2023/08/14 15:39:54 ERROR error
}
