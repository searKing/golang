// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gin

import "github.com/gin-gonic/gin"

// WithRecoveryHandler set a recover handler if panic
func WithRecoveryHandler(f func(c *gin.Context, err interface{}) error) RecoveryOption {
	return RecoveryOptionFunc(func(r *recovery) {
		r.recoveryHandler = f
	})
}
