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
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	signal_ "github.com/searKing/golang/go/os/signal"
)

func TestSetSigInvokeChainContext(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	signal_.DumpSignalTo(int(w.Fd()))
	defer signal_.DumpSignalTo(int(syscall.Stderr))

	// SIGTTOU is never raised, so the chain of SIGWINCH waits until ctx is done.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signal_.SetSigInvokeChainContext(ctx, syscall.SIGWINCH, syscall.SIGTTOU, 0, syscall.SIGWINCH)
	c := make(chan os.Signal, 1)
	signal_.Notify(c, syscall.SIGWINCH)
	defer signal.Stop(c)

	const cancelAfter = 500 * time.Millisecond
	relay := func() time.Duration {
		start := time.Now()
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGWINCH); err != nil {
			t.Fatal(err)
		}
		select {
		case <-c:
		case <-time.After(10 * time.Second):
			t.Fatal("SIGWINCH not relayed, the chain waits forever")
		}
		return time.Since(start)
	}

	time.AfterFunc(cancelAfter, cancel)
	if elapsed := relay(); elapsed < cancelAfter {
		t.Errorf("SIGWINCH relayed after %v, want waiting SIGTTOU until ctx is done after %v", elapsed, cancelAfter)
	}
	// The chain is removed once ctx is done.
	if elapsed := relay(); elapsed >= cancelAfter {
		t.Errorf("SIGWINCH relayed after %v, want no chain after ctx is done", elapsed)
	}

	signal_.DumpSignalTo(int(syscall.Stderr))
	_ = w.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("wait Signal(%d) for Signal(%d) canceled.", int(syscall.SIGTTOU), int(syscall.SIGWINCH))
	if n := strings.Count(string(b), want); n != 1 {
		t.Errorf("signal log = %q, want containing %q once, got %d", b, want, n)
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
