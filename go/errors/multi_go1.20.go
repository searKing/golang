// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.20

package errors

import "errors"

// Deprecated: Use errors.Join instead since go1.20.
func Multi(errs ...error) error {
	return errors.Join(errs...)
}
