// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

package signal_test

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"

	signal_ "github.com/searKing/golang/go/os/signal"
)

func ExampleNotify() {
	signal_.DumpSignalTo(syscall.Stderr)
	signal_.SetSigInvokeChain(syscall.SIGUSR1, syscall.SIGUSR2, 0, syscall.SIGINT)

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	//signal.Notify(c, syscall.SIGINT, syscall.SIGSEGV)
	signal_.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT)

	// simulate to send a SIGINT to this example test
	go func() {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer wg.Done()
		for {
			// Block until a signal is received.
			select {
			case s, ok := <-c:
				if !ok {
					return
				}
				//signal.Stop(c)
				//close(c)
				if s == syscall.SIGUSR1 {
					_, _ = fmt.Fprintf(os.Stderr, "Previous run crashed:\n%s\n", signal_.PreviousStacktrace())
					_ = syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
				} else if s != syscall.SIGUSR2 {
					fmt.Printf("Got signal: %s\n", s)

					// just in case of windows os system, which is based on signal() in C language
					// you can comment below out on unix-like os system.
					signal_.Notify(c, s)
				}
				if s == syscall.SIGINT {
					return
				}

			case <-time.After(time.Minute):
				_, _ = fmt.Fprintf(os.Stderr, "time overseed:\n")
				return
			}
			// set os.Stderr to pass test, for the stacktrace is random.
			// stderr prints something like:
			// Signal received(2).
			// Stacktrace dumped to file: stacktrace.dump.
			// Previous run crashed:
			//  0# searking::SignalHandler::operator()(int, __siginfo*, void*) in /private/var/folders/12/870qx8rd0_d96nt6g078wp080000gn/T/___ExampleSignalAction_in_github_com_searKing_golang_go_os_signal
			//  1# _sigtramp in /usr/lib/system/libsystem_platform.dylib

		}
	}()
	wg.Wait()

	// Output:
	// Got signal: interrupt
}
