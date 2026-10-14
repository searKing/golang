// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"
)

// Locked contains a lock, which must not be copied by WithLocked(v Locked), skipped by -skip-self.
//
//go:generate go-option -type "Locked" -skip-self
type Locked struct {
	Name string
	mu   sync.Mutex `option:"-"`
}

func main() {
	l := &Locked{}
	l.ApplyOptions(WithLockedName("locked"))
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.Name != "locked" {
		panic(fmt.Sprintf("Name = %q, want %q", l.Name, "locked"))
	}
}
