// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multi

import (
	errors_ "github.com/searKing/golang/go/errors"
)

// New returns an error with the supplied errors.
// If no error contained, return nil.
// Deprecated: Use errors.Multi instead.
func New(errs ...error) error {
	return errors_.Multi(errs...)
}
