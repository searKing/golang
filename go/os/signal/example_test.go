// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build unix

package signal_test

import (
	"fmt"
	"os"
	"syscall"

	signal_ "github.com/searKing/golang/go/os/signal"
)

func ExampleNotify() {
	c := make(chan os.Signal, 1)
	signal_.Notify(c, syscall.SIGINT)

	go func() {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	s := <-c
	fmt.Printf("Got signal: %s\n", s)

	// Output:
	// Got signal: interrupt
}
