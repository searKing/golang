// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"errors"
)

var (
	// ErrNoAddrAvailable indicates no Addr is available for pick().
	ErrNoAddrAvailable = errors.New("no Addr is available")
)

// Pick includes additional information for Pick.
//go:generate go-option -type "Pick"
type Pick struct{}

// Picker is used to pick an Address.
type Picker interface {
	Pick(ctx context.Context, addrs []Address, opts ...PickOption) (Address, error)
}

// The PickerFunc type is an adapter to allow the use of
// ordinary functions as Picker handlers. If f is a function
// with the appropriate signature, PickerFunc(f) is a
// Handler that calls f.
type PickerFunc func(ctx context.Context, addrs []Address, opts ...PickOption) (Address, error)

// Pick calls f(ctx, addrs, opts...).
func (f PickerFunc) Pick(ctx context.Context, addrs []Address, opts ...PickOption) (Address, error) {
	return f(ctx, addrs, opts...)
}
