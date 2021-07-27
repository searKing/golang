// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rand implements math/rand functions in a concurrent-safe way
// with a global random source, independent of math/rand's global source.
package rand

import (
	"math/rand"
	"sync"
	"time"
)

var (
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu         sync.Mutex
)

// Seed implements rand.Seed on the global source.
func Seed(seed int64) {
	mu.Lock()
	defer mu.Unlock()
	rand.Int31()
	globalRand.Seed(seed)
}

// Int63 implements rand.Int63 on the global source.
func Int63() int64 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Int63()
}

// Uint32 implements rand.Uint32 on the global source.
func Uint32() uint32 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Uint32()
}

// Uint64 implements rand.Uint64 on the global source.
func Uint64() uint64 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Uint64()
}

// Int31 implements rand.Int31 on the global source.
func Int31() int32 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Int31()
}

// Int implements rand.Int on the global source.
func Int() int {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Int()
}

// Int63n implements rand.Int63n on the global source.
func Int63n(n int64) int64 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Int63n(n)
}

// Int31n implements rand.Int31n on the global source.
func Int31n(n int32) int32 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Int31n(n)
}

// Intn implements rand.Intn on the global source.
func Intn(n int) int {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Intn(n)
}

// Float64 implements rand.Float64 on the global source.
func Float64() float64 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Float64()
}

// Float32 implements rand.Float32 on the global source.
func Float32() float32 {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Float32()
}

// Perm implements rand.Perm on the global source.
func Perm(n int) []int {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Perm(n)
}

// Shuffle implements rand.Shuffle on the global source.
func Shuffle(n int, swap func(i, j int)) {
	mu.Lock()
	defer mu.Unlock()
	globalRand.Shuffle(n, swap)
}

// Read implements rand.Read on the global source.
func Read(p []byte) (n int, err error) {
	mu.Lock()
	defer mu.Unlock()
	return globalRand.Read(p)
}
