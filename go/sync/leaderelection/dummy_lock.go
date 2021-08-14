// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leaderelection

import "context"

// NewDummyLock returns a locker which returns follower
func NewDummyLock(identity string) *dummyLock {
	return &dummyLock{identity: identity}
}

type dummyLock struct {
	identity string
}

// Get is a dummy to allow us to have a dummyLock for testing.
func (fl *dummyLock) Get(ctx context.Context) (ler *Record, rawRecord []byte, err error) {
	return nil, nil, nil
}

// Create is a dummy to allow us to have a dummyLock for testing.
func (fl *dummyLock) Create(ctx context.Context, ler Record) error {
	return nil
}

// Update is a dummy to allow us to have a dummyLock for testing.
func (fl *dummyLock) Update(ctx context.Context, ler Record) error {
	return nil
}

// RecordEvent is a dummy to allow us to have a dummyLock for testing.
func (fl *dummyLock) RecordEvent(name, event string) {}

// Identity is a dummy to allow us to have a dummyLock for testing.
func (fl *dummyLock) Identity() string {
	return fl.identity
}

// Describe is a dummy to allow us to have a dummyLock for testing.
func (fl *dummyLock) Describe() string {
	return "Dummy implementation of lock for testing"
}
