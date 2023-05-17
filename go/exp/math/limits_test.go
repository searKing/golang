// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math_test

import (
	"math"
	"testing"

	math_ "github.com/searKing/golang/go/exp/math"
)

func TestMaxInt(t *testing.T) {
	{
		got := math_.MaxInt[int]()
		want := math.MaxInt
		if got != want {
			t.Errorf("math_.MaxInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MaxInt[int8]()
		want := int8(math.MaxInt8)
		if got != want {
			t.Errorf("math_.MaxInt8[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MaxInt[int16]()
		want := int16(math.MaxInt16)
		if got != want {
			t.Errorf("math_.MaxInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MaxInt[int32]()
		want := int32(math.MaxInt32)
		if got != want {
			t.Errorf("math_.MaxInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MaxInt[int64]()
		want := int64(math.MaxInt64)
		if got != want {
			t.Errorf("math_.MaxInt[%T] = %v, want %v", got, got, want)
		}
	}

	{
		got := math_.MaxInt[uint]()
		want := uint(math.MaxUint)
		if got != want {
			t.Errorf("math_.MaxInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MaxInt[uint8]()
		want := uint8(math.MaxUint8)
		if got != want {
			t.Errorf("math_.MaxInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MaxInt[uint16]()
		want := uint16(math.MaxUint16)
		if got != want {
			t.Errorf("math_.MaxInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MaxInt[uint32]()
		want := uint32(math.MaxUint32)
		if got != want {
			t.Errorf("math_.MaxInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MaxInt[uint64]()
		want := uint64(math.MaxUint64)
		if got != want {
			t.Errorf("math_.MaxInt[%T] = %v, want %v", got, got, want)
		}
	}
}

func TestMinInt(t *testing.T) {
	{
		got := math_.MinInt[int]()
		want := math.MinInt
		if got != math.MinInt {
			t.Errorf("math_.MinInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MinInt[int8]()
		want := int8(math.MinInt8)
		if got != want {
			t.Errorf("math_.MinInt8[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MinInt[int16]()
		want := int16(math.MinInt16)
		if got != want {
			t.Errorf("math_.MinInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MinInt[int32]()
		want := int32(math.MinInt32)
		if got != want {
			t.Errorf("math_.MinInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MinInt[int64]()
		want := int64(math.MinInt64)
		if got != want {
			t.Errorf("math_.MinInt[%T] = %v, want %v", got, got, want)
		}
	}

	{
		got := math_.MinInt[uint]()
		want := uint(0)
		if got != want {
			t.Errorf("math_.MinInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MinInt[uint8]()
		want := uint8(0)
		if got != want {
			t.Errorf("math_.MinInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MinInt[uint16]()
		want := uint16(0)
		if got != want {
			t.Errorf("math_.MinInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MinInt[uint32]()
		want := uint32(0)
		if got != want {
			t.Errorf("math_.MinInt[%T] = %v, want %v", got, got, want)
		}
	}
	{
		got := math_.MinInt[uint64]()
		want := uint64(0)
		if got != want {
			t.Errorf("math_.MinInt[%T] = %v, want %v", got, got, want)
		}
	}
}
