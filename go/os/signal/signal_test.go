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
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
	"unsafe"

	signal_ "github.com/searKing/golang/go/os/signal"
)

const fatalForwardingSentinel = "reached after synchronous SIGSEGV"

func TestNotifyForwardsSynchronousFatalSignal(t *testing.T) {
	if os.Getenv("SEARKING_SIGNAL_FATAL_FORWARDING_TEST") == "1" {
		if got := os.Getenv("GOTRACEBACK"); got != "system" {
			fmt.Fprintf(os.Stderr, "GOTRACEBACK=%q, want system\n", got)
			os.Exit(2)
		}
		signal_.DumpSignalTo(-1)
		signal_.DumpStacktraceTo(os.Getenv("SEARKING_SIGNAL_FATAL_DUMP_PATH"))
		c := make(chan os.Signal, 1)
		signal_.Notify(c, syscall.SIGSEGV)
		triggerSynchronousSIGSEGV()
		fmt.Print(fatalForwardingSentinel)
		os.Exit(0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dumpPath := filepath.Join(t.TempDir(), "fatal.dump")
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestNotifyForwardsSynchronousFatalSignal$")
	cmd.Env = replaceEnv(os.Environ(), "GOTRACEBACK", "system")
	cmd.Env = replaceEnv(cmd.Env, "SEARKING_SIGNAL_FATAL_FORWARDING_TEST", "1")
	cmd.Env = replaceEnv(cmd.Env, "SEARKING_SIGNAL_FATAL_DUMP_PATH", dumpPath)
	output, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		t.Fatalf("fatal-forwarding child timed out: %v", ctx.Err())
	}
	if err == nil {
		t.Fatalf("fatal-forwarding child exited successfully\n%s", output)
	}
	if strings.Contains(string(output), fatalForwardingSentinel) {
		t.Fatalf("fatal-forwarding child reached sentinel\n%s", output)
	}
	dumpInfo, err := os.Stat(dumpPath)
	if err != nil {
		t.Fatalf("fatal-forwarding child did not write native stack dump: %v\n%s", err, output)
	}
	if dumpInfo.Size() == 0 {
		t.Fatalf("fatal-forwarding child wrote an empty native stack dump\n%s", output)
	}
	if !strings.Contains(string(output), "[signal SIGSEGV:") {
		t.Fatalf("fatal-forwarding child output does not show Go runtime fault handling\n%s", output)
	}
}

func replaceEnv(env []string, key, value string) []string {
	prefix := key + "="
	replaced := make([]string, 0, len(env)+1)
	for _, entry := range env {
		if !strings.HasPrefix(entry, prefix) {
			replaced = append(replaced, entry)
		}
	}
	return append(replaced, prefix+value)
}

//go:nocheckptr
func triggerSynchronousSIGSEGV() {
	*(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(nil)) + 0x123)) = 1
}

func TestNotifyConcurrent(t *testing.T) {
	const workers = 8
	channels := make([]chan os.Signal, workers)
	var wg sync.WaitGroup
	for i := range channels {
		channels[i] = make(chan os.Signal, 1)
		wg.Add(1)
		go func(c chan os.Signal) {
			defer wg.Done()
			signal_.Notify(c, syscall.SIGUSR2)
		}(channels[i])
	}
	wg.Wait()
	for _, c := range channels {
		defer signal.Stop(c)
	}
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR2); err != nil {
		t.Fatal(err)
	}
	for i, c := range channels {
		select {
		case <-c:
		case <-time.After(5 * time.Second):
			t.Fatalf("channel %d did not receive SIGUSR2", i)
		}
	}
}

func TestForwardSAInfoSpecialDispositions(t *testing.T) {
	cxxOutput, err := exec.Command("go", "env", "CXX").Output()
	if err != nil {
		t.Fatal(err)
	}
	cxx := strings.TrimSpace(string(cxxOutput))
	goosOutput, err := exec.Command("go", "env", "GOOS").Output()
	if err != nil {
		t.Fatal(err)
	}
	goos := strings.TrimSpace(string(goosOutput))

	testDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cgoDir := filepath.Join(testDir, "cgo")
	dir := t.TempDir()
	source := filepath.Join(dir, "signal_disposition_test.cpp")
	program := `#include <csignal>
#include <cstring>

#include "signal_handler_unix.hpp"

int main(int argc, char **argv) {
  if (argc != 2) {
    return 2;
  }
  struct sigaction action;
  memset(&action, 0, sizeof(action));
  sigemptyset(&action.sa_mask);
  action.sa_flags = SA_SIGINFO;
  if (strcmp(argv[1], "ignore") == 0) {
    action.sa_sigaction =
        reinterpret_cast<decltype(action.sa_sigaction)>(SIG_IGN);
  } else if (strcmp(argv[1], "default") == 0) {
    action.sa_sigaction =
        reinterpret_cast<decltype(action.sa_sigaction)>(SIG_DFL);
  } else {
    return 3;
  }
  if (sigaction(SIGUSR1, &action, nullptr) != 0 ||
      searking::SignalHandler::SetSig(SIGUSR1) != 0) {
    return 4;
  }
  raise(SIGUSR1);
  return strcmp(argv[1], "ignore") == 0 ? 0 : 5;
}
`
	if err := os.WriteFile(source, []byte(program), 0o600); err != nil {
		t.Fatal(err)
	}
	helper := filepath.Join(dir, "signal_disposition_test")
	args := []string{
		"-std=c++11",
		"-D_GNU_SOURCE",
		"-I", cgoDir,
		"-I", filepath.Join(cgoDir, "include"),
		source,
		filepath.Join(cgoDir, "signal_handler_unix.cpp"),
		filepath.Join(cgoDir, "write_int.cpp"),
		"-o", helper,
	}
	if goos == "linux" {
		args = append(args, "-ldl")
	}
	if output, err := exec.Command(cxx, args...).CombinedOutput(); err != nil {
		t.Fatalf("compile C++ disposition helper: %v\n%s", err, output)
	}

	t.Run("SIG_IGN", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if output, err := exec.CommandContext(ctx, helper, "ignore").CombinedOutput(); err != nil {
			t.Fatalf("SIG_IGN helper failed: %v\n%s", err, output)
		}
	})

	t.Run("SIG_DFL", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, helper, "default")
		output, err := cmd.CombinedOutput()
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("SIG_DFL helper error = %v, want signal termination\n%s", err, output)
		}
		status, ok := exitErr.Sys().(syscall.WaitStatus)
		if !ok || !status.Signaled() || status.Signal() != syscall.SIGUSR1 {
			t.Fatalf("SIG_DFL helper status = %v, want SIGUSR1 termination\n%s", status, output)
		}
	})
}

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

func TestDumpStacktracePathConcurrentUpdate(t *testing.T) {
	dir := t.TempDir()
	paths := []string{filepath.Join(dir, "a.dump"), filepath.Join(dir, "b.dump")}
	c := make(chan os.Signal, 32)
	signal_.Notify(c, syscall.SIGUSR2)
	defer signal.Stop(c)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 32; i++ {
			signal_.DumpStacktraceTo(paths[i%len(paths)])
		}
	}()
	for i := 0; i < 32; i++ {
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR2); err != nil {
			t.Fatal(err)
		}
		select {
		case <-c:
		case <-time.After(5 * time.Second):
			t.Fatal("SIGUSR2 not relayed")
		}
	}
	wg.Wait()
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
