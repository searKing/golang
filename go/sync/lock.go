// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import "sync"

// WithLock runs while holding lk.
func WithLock(lk sync.Locker, fn func()) {
	lk.Lock()
	defer lk.Unlock() // in case fn panics
	fn()
}
