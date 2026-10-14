// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix && cgo

package signal_test

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	signal_ "github.com/searKing/golang/go/os/signal"
)

func TestNotifySkipsRuntimeSignals(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	signal_.DumpSignalTo(int(w.Fd()))
	defer signal_.DumpSignalTo(int(syscall.Stderr))

	c := make(chan os.Signal, 8)
	signal_.Notify(c, syscall.SIGURG, syscall.SIGUSR2)
	defer signal.Stop(c)

	pid := syscall.Getpid()
	if err := syscall.Kill(pid, syscall.SIGURG); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Kill(pid, syscall.SIGUSR2); err != nil {
		t.Fatal(err)
	}
	// The cgo signal handler logs before it invokes the Go signal handler, which relays the signal to c.
	timeout := time.After(10 * time.Second)
	for got := false; !got; {
		select {
		case s := <-c:
			got = s == syscall.SIGUSR2
		case <-timeout:
			t.Fatal("SIGUSR2 not relayed")
		}
	}

	signal_.DumpSignalTo(int(syscall.Stderr))
	_ = w.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	log := string(b)
	if want := fmt.Sprintf("Signal received(%d)", int(syscall.SIGUSR2)); !strings.Contains(log, want) {
		t.Errorf("signal log = %q, want containing %q", log, want)
	}
	if unwanted := fmt.Sprintf("Signal received(%d)", int(syscall.SIGURG)); strings.Contains(log, unwanted) {
		t.Errorf("signal log = %q, want SIGURG handled by the Go runtime only, without %q", log, unwanted)
	}
}
