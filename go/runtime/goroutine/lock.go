// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Defensive debug-only utility to track that functions run on the
// goroutine that they're supposed to.

package goroutine

import (
	"errors"
	"os"

	"github.com/searKing/golang/go/error/must"
)

var DebugGoroutines = os.Getenv("DEBUG_GOROUTINES") == "1"

// Lock represents a goroutine ID, with goroutine ID checked, that is  whether GoRoutines of lock newer and check caller differ.
// disable when DebugGoroutines equals false
type Lock uint64

// NewLock returns a goroutine Lock, that checks whether goroutine of lock newer and check caller differ.
// Code borrowed from https://github.com/golang/go/blob/master/src/net/http/h2_bundle.go
func NewLock() Lock {
	if !DebugGoroutines {
		return 0
	}
	return Lock(ID())
}

// Check if caller's goroutine is locked
func (g Lock) Check() error {
	if !DebugGoroutines {
		return nil
	}
	if ID() != uint64(g) {
		return errors.New("running on the wrong goroutine")
	}
	return nil
}

func (g Lock) MustCheck() {
	must.Must(g.Check())
}

// Check whether caller's goroutine escape lock
func (g Lock) CheckNotOn() error {
	if !DebugGoroutines {
		return nil
	}
	if ID() == uint64(g) {
		return errors.New("running on the wrong goroutine")
	}
	return nil
}

func (g Lock) MustCheckNotOn() {
	must.Must(g.CheckNotOn())
}
