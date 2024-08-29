// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package burstlimit

// fullChan returns a channel of the given size, filled with empty structs.
// simple, but effective.
//
// Also see: https://github.com/searKing/golang/blob/go/v1.2.118/go/time/rate/rate.go#L87
func fullChan(b int) (limiter chan struct{}) {
	if b > 0 {
		limiter = make(chan struct{}, b)
		for i := 0; i < b; i++ {
			limiter <- struct{}{}
		}
	}
	return limiter
}
