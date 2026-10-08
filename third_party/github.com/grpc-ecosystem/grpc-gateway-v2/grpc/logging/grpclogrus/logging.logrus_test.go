// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpclogrus

import (
	"context"
	"testing"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestInterceptorLogrusLogger_Fields(t *testing.T) {
	logger, hook := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)
	interceptorLogrusLogger(logger).Log(context.Background(), logging.LevelInfo, "finished call",
		"protocol", "grpc", "grpc.method", "Get", "grpc.code", "OK")

	e := hook.LastEntry()
	if e == nil {
		t.Fatal("nothing logged")
	}
	want := logrus.Fields{"protocol": "grpc", "grpc.method": "Get", "grpc.code": "OK"}
	if len(e.Data) != len(want) {
		t.Errorf("fields = %v, want %v", e.Data, want)
	}
	for k, v := range want {
		if e.Data[k] != v {
			t.Errorf("field %s = %v, want %v", k, e.Data[k], v)
		}
	}
	if e.Level != logrus.InfoLevel || e.Message != "finished call" {
		t.Errorf("log = %v %q, want %v %q", e.Level, e.Message, logrus.InfoLevel, "finished call")
	}
}
