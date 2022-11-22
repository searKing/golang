// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expvar

import (
	"expvar"
	"fmt"
)

var _ expvar.Var = (*Leak)(nil)

// Leak is a pair of 64-bit integer variables that satisfies the Var interface.
type Leak struct {
	Leak expvar.Int // Leak = New - Delete
	New  expvar.Int // New
}

// Value returns the Leak counter.
func (v *Leak) Value() int64 {
	return v.Leak.Value()
}

func (v *Leak) String() string {
	return fmt.Sprint([]string{v.Leak.String(), v.New.String()})
}

// Add adds delta, which may be negative, to the Leak counter.
func (v *Leak) Add(delta int64) {
	if delta >= 0 {
		v.New.Add(delta)
	}
	v.Leak.Add(delta)
}

// Done decrements the Leak counter by one.
func (v *Leak) Done() {
	v.Add(-1)
}

// Convenience functions for creating new exported variables.

func NewLeak(name string) *Leak {
	v := new(Leak)
	expvar.Publish(name, v)
	return v
}
