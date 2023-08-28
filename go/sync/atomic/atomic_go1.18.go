// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.19

package atomic

import (
	"math"
	"sync/atomic"
	"time"
)

// Int32 is an atomic wrapper around an int32.
type Int32 int32

// NewInt32 creates an Int32.
func NewInt32(i int32) *Int32 {
	return (*Int32)(&i)
}

// Load atomically loads the wrapped value.
func (i *Int32) Load() int32 {
	return atomic.LoadInt32((*int32)(i))
}

// Add atomically adds to the wrapped int32 and returns the new value.
func (i *Int32) Add(n int32) int32 {
	return atomic.AddInt32((*int32)(i), n)
}

// Sub atomically subtracts from the wrapped int32 and returns the new value.
func (i *Int32) Sub(n int32) int32 {
	return atomic.AddInt32((*int32)(i), -n)
}

// Inc atomically increments the wrapped int32 and returns the new value.
func (i *Int32) Inc() int32 {
	return i.Add(1)
}

// Dec atomically decrements the wrapped int32 and returns the new value.
func (i *Int32) Dec() int32 {
	return i.Sub(1)
}

// CAS is an atomic compare-and-swap.
func (i *Int32) CAS(old, new int32) bool {
	return atomic.CompareAndSwapInt32((*int32)(i), old, new)
}

// Store atomically stores the passed value.
func (i *Int32) Store(n int32) {
	atomic.StoreInt32((*int32)(i), n)
}

// Swap atomically swaps the wrapped int32 and returns the old value.
func (i *Int32) Swap(n int32) int32 {
	return atomic.SwapInt32((*int32)(i), n)
}

// Int64 is an atomic wrapper around an int64.
type Int64 int64

// NewInt64 creates an Int64.
func NewInt64(i int64) *Int64 {
	return (*Int64)(&i)
}

// Load atomically loads the wrapped value.
func (i *Int64) Load() int64 {
	return atomic.LoadInt64((*int64)(i))
}

// Add atomically adds to the wrapped int64 and returns the new value.
func (i *Int64) Add(n int64) int64 {
	return atomic.AddInt64((*int64)(i), n)
}

// Sub atomically subtracts from the wrapped int64 and returns the new value.
func (i *Int64) Sub(n int64) int64 {
	return atomic.AddInt64((*int64)(i), -n)
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
	return atomic.CompareAndSwapInt64((*int64)(i), old, new)
}

// Store atomically stores the passed value.
func (i *Int64) Store(n int64) {
	atomic.StoreInt64((*int64)(i), n)
}

// Swap atomically swaps the wrapped int64 and returns the old value.
func (i *Int64) Swap(n int64) int64 {
	return atomic.SwapInt64((*int64)(i), n)
}

// Uint32 is an atomic wrapper around an uint32.
type Uint32 uint32

// NewUint32 creates a Uint32.
func NewUint32(i uint32) *Uint32 {
	return (*Uint32)(&i)
}

// Load atomically loads the wrapped value.
func (i *Uint32) Load() uint32 {
	return atomic.LoadUint32((*uint32)(i))
}

// Add atomically adds to the wrapped uint32 and returns the new value.
func (i *Uint32) Add(n uint32) uint32 {
	return atomic.AddUint32((*uint32)(i), n)
}

// Sub atomically subtracts from the wrapped uint32 and returns the new value.
func (i *Uint32) Sub(n uint32) uint32 {
	return atomic.AddUint32((*uint32)(i), ^(n - 1))
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
	return atomic.CompareAndSwapUint32((*uint32)(i), old, new)
}

// Store atomically stores the passed value.
func (i *Uint32) Store(n uint32) {
	atomic.StoreUint32((*uint32)(i), n)
}

// Swap atomically swaps the wrapped uint32 and returns the old value.
func (i *Uint32) Swap(n uint32) uint32 {
	return atomic.SwapUint32((*uint32)(i), n)
}

// Uint64 is an atomic wrapper around a uint64.
type Uint64 uint64

// NewUint64 creates a Uint64.
func NewUint64(i uint64) *Uint64 {
	return (*Uint64)(&i)
}

// Load atomically loads the wrapped value.
func (i *Uint64) Load() uint64 {
	return atomic.LoadUint64((*uint64)(i))
}

// Add atomically adds to the wrapped uint64 and returns the new value.
func (i *Uint64) Add(n uint64) uint64 {
	return atomic.AddUint64((*uint64)(i), n)
}

// Sub atomically subtracts from the wrapped uint64 and returns the new value.
func (i *Uint64) Sub(n uint64) uint64 {
	return atomic.AddUint64((*uint64)(i), ^(n - 1))
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
	return atomic.CompareAndSwapUint64((*uint64)(i), old, new)
}

// Store atomically stores the passed value.
func (i *Uint64) Store(n uint64) {
	atomic.StoreUint64((*uint64)(i), n)
}

// Swap atomically swaps the wrapped uint64 and returns the old value.
func (i *Uint64) Swap(n uint64) uint64 {
	return atomic.SwapUint64((*uint64)(i), n)
}

// Bool is an atomic Boolean.
type Bool Uint32

// NewBool creates a Bool.
func NewBool(initial bool) *Bool {
	b := boolToUint32(initial)
	return (*Bool)(&b)
}

// Load atomically loads the Boolean.
func (b *Bool) Load() bool {

	return truthy((*Uint32)(b).Load())
}

// CAS is an atomic compare-and-swap.
func (b *Bool) CAS(old, new bool) bool {
	return (*Uint32)(b).CAS(boolToUint32(old), boolToUint32(new))
}

// Store atomically stores the passed value.
func (b *Bool) Store(new bool) {
	(*Uint32)(b).Store(boolToUint32(new))
}

// Swap sets the given value and returns the previous value.
func (b *Bool) Swap(new bool) bool {
	return truthy((*Uint32)(b).Swap(boolToUint32(new)))
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

func truthy(n uint32) bool {
	return n&1 == 1
}

func boolToUint32(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// Float32 is an atomic wrapper around float32.
type Float32 uint32

// NewFloat32 creates a Float32.
func NewFloat32(f float32) *Float32 {
	u := math.Float32bits(f)
	return (*Float32)(&u)

}

// Load atomically loads the wrapped value.
func (f *Float32) Load() float32 {
	return math.Float32frombits(atomic.LoadUint32((*uint32)(f)))
}

// Store atomically stores the passed value.
func (f *Float32) Store(s float32) {
	atomic.StoreUint32((*uint32)(f), math.Float32bits(s))
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
	return atomic.CompareAndSwapUint32((*uint32)(f), math.Float32bits(old), math.Float32bits(new))
}

// Float64 is an atomic wrapper around float64.
type Float64 uint64

// NewFloat64 creates a Float64.
func NewFloat64(f float64) *Float64 {
	u := math.Float64bits(f)
	return (*Float64)(&u)

}

// Load atomically loads the wrapped value.
func (f *Float64) Load() float64 {
	return math.Float64frombits(atomic.LoadUint64((*uint64)(f)))
}

// Store atomically stores the passed value.
func (f *Float64) Store(s float64) {
	atomic.StoreUint64((*uint64)(f), math.Float64bits(s))
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
	return atomic.CompareAndSwapUint64((*uint64)(f), math.Float64bits(old), math.Float64bits(new))
}

// Duration is an atomic wrapper around time.Duration
// https://godoc.org/time#Duration
type Duration Int64

// NewDuration creates a Duration.
func NewDuration(d time.Duration) *Duration {
	return (*Duration)(NewInt64(int64(d)))
}

// Load atomically loads the wrapped value.
func (d *Duration) Load() time.Duration {
	return time.Duration((*Int64)(d).Load())
}

// Store atomically stores the passed value.
func (d *Duration) Store(n time.Duration) {
	(*Int64)(d).Store(int64(n))
}

// Add atomically adds to the wrapped time.Duration and returns the new value.
func (d *Duration) Add(n time.Duration) time.Duration {
	return time.Duration((*Int64)(d).Add(int64(n)))
}

// Sub atomically subtracts from the wrapped time.Duration and returns the new value.
func (d *Duration) Sub(n time.Duration) time.Duration {
	return time.Duration((*Int64)(d).Sub(int64(n)))
}

// Swap atomically swaps the wrapped time.Duration and returns the old value.
func (d *Duration) Swap(n time.Duration) time.Duration {
	return time.Duration((*Int64)(d).Swap(int64(n)))
}

// CAS is an atomic compare-and-swap.
func (d *Duration) CAS(old, new time.Duration) bool {
	return (*Int64)(d).CAS(int64(old), int64(new))
}

// Value shadows the type of the same name from sync/atomic
// https://godoc.org/sync/atomic#Value
type Value atomic.Value
