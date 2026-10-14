// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix && cgo

package signal_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	signal_ "github.com/searKing/golang/go/os/signal"
)

func TestDeprecatedSigInvokeChainDoesNotBlock(t *testing.T) {
	if mode := os.Getenv("SEARKING_SIGNAL_CHAIN_TEST"); mode != "" {
		switch mode {
		case "plain":
			signal_.SetSigInvokeChain(syscall.SIGUSR1, syscall.SIGTTOU, 10, syscall.SIGWINCH)
		case "context":
			signal_.SetSigInvokeChainContext(context.Background(), syscall.SIGUSR1, syscall.SIGTTOU, 10*time.Second, syscall.SIGWINCH)
		default:
			t.Fatalf("unknown child mode %q", mode)
		}
		c := make(chan os.Signal, 1)
		signal_.Notify(c, syscall.SIGWINCH)
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGWINCH); err != nil {
			t.Fatal(err)
		}
		select {
		case <-c:
		case <-time.After(500 * time.Millisecond):
			t.Fatal("deprecated signal chain blocked delivery")
		}
		return
	}

	for _, mode := range []string{"plain", "context"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestDeprecatedSigInvokeChainDoesNotBlock$")
			cmd.Env = append(os.Environ(), "SEARKING_SIGNAL_CHAIN_TEST="+mode)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("child failed: %v\n%s", err, output)
			}
		})
	}
}

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
