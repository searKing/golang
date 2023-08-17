// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// Any returns its argument v or nil if and only if v is a typed nil or an untyped nil.
// For example, if v was created by set with `var p *int` or calling [Any]((*int)(nil)),
// [Any] returns nil.
func Any[T any](v *T) any {
	if v == nil {
		return nil
	}
	return v
}
