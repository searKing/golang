// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.19

package atomic

import (
	"math"
	"sync/atomic"
	"time"
)

// Int32 is an atomic wrapper around an int32.
// Deprecated: Use atomic.Int32 instead since go1.19.
type Int32 atomic.Int32

// NewInt32 creates an Int32.
func NewInt32(i int32) *Int32 {
	var v Int32
	v.Store(i)
	return &v
}

// Load atomically loads the wrapped value.
func (i *Int32) Load() int32 {
	return (*atomic.Int32)(i).Load()
}

// Add atomically adds to the wrapped int32 and returns the new value.
func (i *Int32) Add(n int32) int32 {
	return (*atomic.Int32)(i).Add(n)
}

// Sub atomically subtracts from the wrapped int32 and returns the new value.
func (i *Int32) Sub(n int32) int32 {
	return (*atomic.Int32)(i).Add(-n)
}

// Inc atomically increments the wrapped int32 and returns the new value.
func (i *Int32) Inc() int32 {
	return (*atomic.Int32)(i).Add(1)
}

// Dec atomically decrements the wrapped int32 and returns the new value.
func (i *Int32) Dec() int32 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Int32) CAS(old, new int32) bool {
	return (*atomic.Int32)(i).CompareAndSwap(old, new)
}

// Store atomically stores the passed value.
func (i *Int32) Store(n int32) {
	(*atomic.Int32)(i).Store(n)
}

// Swap atomically swaps the wrapped int32 and returns the old value.
func (i *Int32) Swap(n int32) int32 {
	return (*atomic.Int32)(i).Swap(n)
}

// Int64 is an atomic wrapper around an int64.
// Deprecated: Use atomic.Int64 instead since go1.19.
type Int64 atomic.Int64

// NewInt64 creates an Int64.
func NewInt64(i int64) *Int64 {
	var v Int64
	v.Store(i)
	return &v
}

// Load atomically loads the wrapped value.
func (i *Int64) Load() int64 {
	return (*atomic.Int64)(i).Load()
}

// Add atomically adds to the wrapped int64 and returns the new value.
func (i *Int64) Add(n int64) int64 {
	return (*atomic.Int64)(i).Add(n)
}

// Sub atomically subtracts from the wrapped int64 and returns the new value.
func (i *Int64) Sub(n int64) int64 {
	return (*atomic.Int64)(i).Add(-n)
}

// Inc atomically increments the wrapped int64 and returns the new value.
func (i *Int64) Inc() int64 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped int64 and returns the new value.
func (i *Int64) Dec() int64 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Int64) CAS(old, new int64) bool {
	return (*atomic.Int64)(i).CompareAndSwap(old, new)
}

// Store atomically stores the passed value.
func (i *Int64) Store(n int64) {
	(*atomic.Int64)(i).Store(n)
}

// Swap atomically swaps the wrapped int64 and returns the old value.
func (i *Int64) Swap(n int64) int64 {
	return (*atomic.Int64)(i).Swap(n)
}

// Uint32 is an atomic wrapper around an uint32.
// Deprecated: Use atomic.Uint32 instead since go1.19.
type Uint32 atomic.Uint32

// NewUint32 creates a Uint32.
func NewUint32(i uint32) *Uint32 {
	var v Uint32
	v.Store(i)
	return &v
}

// Load atomically loads the wrapped value.
func (i *Uint32) Load() uint32 {
	return (*atomic.Uint32)(i).Load()
}

// Add atomically adds to the wrapped uint32 and returns the new value.
func (i *Uint32) Add(n uint32) uint32 {
	return (*atomic.Uint32)(i).Add(n)
}

// Sub atomically subtracts from the wrapped uint32 and returns the new value.
func (i *Uint32) Sub(n uint32) uint32 {
	return (*atomic.Uint32)(i).Add(^(n - 1))
}

// Inc atomically increments the wrapped uint32 and returns the new value.
func (i *Uint32) Inc() uint32 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped int32 and returns the new value.
func (i *Uint32) Dec() uint32 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Uint32) CAS(old, new uint32) bool {
	return (*atomic.Uint32)(i).CompareAndSwap(old, new)
}

// Store atomically stores the passed value.
func (i *Uint32) Store(n uint32) {
	(*atomic.Uint32)(i).Store(n)
}

// Swap atomically swaps the wrapped uint32 and returns the old value.
func (i *Uint32) Swap(n uint32) uint32 {
	return (*atomic.Uint32)(i).Swap(n)
}

// Uint64 is an atomic wrapper around a uint64.
// Deprecated: Use atomic.Uint64 instead since go1.19.
type Uint64 atomic.Uint64

// NewUint64 creates a Uint64.
func NewUint64(i uint64) *Uint64 {
	var v Uint64
	v.Store(i)
	return &v
}

// Load atomically loads the wrapped value.
func (i *Uint64) Load() uint64 {
	return (*atomic.Uint64)(i).Load()
}

// Add atomically adds to the wrapped uint64 and returns the new value.
func (i *Uint64) Add(n uint64) uint64 {
	return (*atomic.Uint64)(i).Add(n)
}

// Sub atomically subtracts from the wrapped uint64 and returns the new value.
func (i *Uint64) Sub(n uint64) uint64 {
	return (*atomic.Uint64)(i).Add(^(n - 1))
}

// Inc atomically increments the wrapped uint64 and returns the new value.
func (i *Uint64) Inc() uint64 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped uint64 and returns the new value.
func (i *Uint64) Dec() uint64 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Uint64) CAS(old, new uint64) bool {
	return (*atomic.Uint64)(i).CompareAndSwap(old, new)
}

// Store atomically stores the passed value.
func (i *Uint64) Store(n uint64) {
	(*atomic.Uint64)(i).Store(n)
}

// Swap atomically swaps the wrapped uint64 and returns the old value.
func (i *Uint64) Swap(n uint64) uint64 {
	return (*atomic.Uint64)(i).Swap(n)
}

// Bool is an atomic Boolean.
// Deprecated: Use atomic.Bool instead since go1.19.
type Bool atomic.Bool

// NewBool creates a Bool.
func NewBool(i bool) *Bool {
	var v Bool
	v.Store(i)
	return &v
}

// Load atomically loads the Boolean.
func (b *Bool) Load() bool {
	return (*atomic.Bool)(b).Load()
}

// CAS is an atomic compare-and-swap.
func (b *Bool) CAS(old, new bool) bool {
	return (*atomic.Bool)(b).CompareAndSwap(old, new)
}

// Store atomically stores the passed value.
func (b *Bool) Store(new bool) {
	(*atomic.Bool)(b).Store(new)
}

// Swap sets the given value and returns the previous value.
func (b *Bool) Swap(new bool) bool {
	return (*atomic.Bool)(b).Swap(new)
}

// Toggle atomically negates the Boolean and returns the previous value.
func (b *Bool) Toggle() bool {
	for {
		old := b.Load()
		if b.CAS(old, !old) {
			return old
		}
	}
}

// Float32 is an atomic wrapper around float32.
type Float32 atomic.Uint32

// NewFloat32 creates a Float32.
func NewFloat32(i float32) *Float32 {
	var v Float32
	v.Store(i)
	return &v
}

// Load atomically loads the wrapped value.
func (f *Float32) Load() float32 {
	return math.Float32frombits((*atomic.Uint32)(f).Load())
}

// Store atomically stores the passed value.
func (f *Float32) Store(s float32) {
	(*atomic.Uint32)(f).Store(math.Float32bits(s))
}

// Add atomically adds to the wrapped float32 and returns the new value.
func (f *Float32) Add(s float32) float32 {
	for {
		old := f.Load()
		new := old + s
		if f.CAS(old, new) {
			return new
		}
	}
}

// Sub atomically subtracts from the wrapped float32 and returns the new value.
func (f *Float32) Sub(s float32) float32 {
	return f.Add(-s)
}

// CAS is an atomic compare-and-swap.
func (f *Float32) CAS(old, new float32) bool {
	return (*atomic.Uint32)(f).CompareAndSwap(math.Float32bits(old), math.Float32bits(new))
}

// Float64 is an atomic wrapper around float64.
type Float64 atomic.Uint64

// NewFloat64 creates a Float64.
func NewFloat64(i float64) *Float64 {
	var v Float64
	v.Store(i)
	return &v
}

// Load atomically loads the wrapped value.
func (f *Float64) Load() float64 {
	return math.Float64frombits((*atomic.Uint64)(f).Load())
}

// Store atomically stores the passed value.
func (f *Float64) Store(s float64) {
	(*atomic.Uint64)(f).Store(math.Float64bits(s))
}

// Add atomically adds to the wrapped float64 and returns the new value.
func (f *Float64) Add(s float64) float64 {
	for {
		old := f.Load()
		new := old + s
		if f.CAS(old, new) {
			return new
		}
	}
}

// Sub atomically subtracts from the wrapped float64 and returns the new value.
func (f *Float64) Sub(s float64) float64 {
	return f.Add(-s)
}

// CAS is an atomic compare-and-swap.
func (f *Float64) CAS(old, new float64) bool {
	return (*atomic.Uint64)(f).CompareAndSwap(math.Float64bits(old), math.Float64bits(new))
}

// Duration is an atomic wrapper around time.Duration
// https://godoc.org/time#Duration
type Duration atomic.Int64

// NewDuration creates a Duration.
func NewDuration(d time.Duration) *Duration {
	var v Duration
	v.Store(d)
	return &v
}

// Load atomically loads the wrapped value.
func (d *Duration) Load() time.Duration {
	return time.Duration((*atomic.Int64)(d).Load())
}

// Store atomically stores the passed value.
func (d *Duration) Store(n time.Duration) {
	(*atomic.Int64)(d).Store(int64(n))
}

// Add atomically adds to the wrapped time.Duration and returns the new value.
func (d *Duration) Add(n time.Duration) time.Duration {
	return time.Duration((*atomic.Int64)(d).Add(int64(n)))
}

// Sub atomically subtracts from the wrapped time.Duration and returns the new value.
func (d *Duration) Sub(n time.Duration) time.Duration {
	return time.Duration((*atomic.Int64)(d).Add(int64(-n)))
}

// Swap atomically swaps the wrapped time.Duration and returns the old value.
func (d *Duration) Swap(n time.Duration) time.Duration {
	return time.Duration((*atomic.Int64)(d).Swap(int64(n)))
}

// CAS is an atomic compare-and-swap.
func (d *Duration) CAS(old, new time.Duration) bool {
	return (*atomic.Int64)(d).CompareAndSwap(int64(old), int64(new))
}

// Value shadows the type of the same name from sync/atomic
// https://godoc.org/sync/atomic#Value
type Value atomic.Value
