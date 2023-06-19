// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"errors"
)

// IsAny reports whether any error in err's tree matches any target in targets.
func IsAny(err error, targets ...error) bool {
	if len(targets) == 0 {
		// Avoid scanning all targets.
		return false
	}
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// IsAll reports whether any error in err's tree matches all target in targets.
func IsAll(err error, targets ...error) bool {
	if len(targets) == 0 {
		return false
	}
	for _, target := range targets {
		if !errors.Is(err, target) {
			// Avoid scanning all targets.
			return false
		}
	}
	return true
}
