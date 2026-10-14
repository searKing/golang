// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !plan9

package signal_test

import (
	"os"
	"syscall"
	"testing"

	signal_ "github.com/searKing/golang/go/os/signal"
)

// TestSignum runs on windows too, so that the cgo sources of windows are built and linked by tests.
func TestSignum(t *testing.T) {
	tests := []struct {
		sig  os.Signal
		want int
	}{
		{syscall.SIGINT, 2},
		{syscall.SIGTERM, int(syscall.SIGTERM)},
		{syscall.Signal(65), -1},
		{os.Interrupt, int(syscall.SIGINT)},
	}
	for _, tt := range tests {
		if got := signal_.Signum(tt.sig); got != tt.want {
			t.Errorf("Signum(%v) = %d, want %d", tt.sig, got, tt.want)
		}
	}
}
