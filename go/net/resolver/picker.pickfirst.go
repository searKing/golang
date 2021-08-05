// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"os"
)

func PickFirst(ctx context.Context, addrs []Address, opts ...PickOption) (Address, error) {
	if len(addrs) == 0 {
		return Address{}, os.ErrNotExist
	}
	return addrs[0], nil
}
