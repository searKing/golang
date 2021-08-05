// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import "context"

// NewErrPicker returns a Picker that always returns err on Pick().
func NewErrPicker(err error) Picker {
	return &errPicker{err: err}
}

type errPicker struct {
	err error // Pick() always returns this err.
}

func (p *errPicker) Pick(ctx context.Context, addrs []Address, opts ...PickOption) (Address, error) {
	return Address{}, p.err
}
