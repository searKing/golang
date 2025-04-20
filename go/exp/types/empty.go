// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "reflect"

// IsEmptyStruct returns true if the type of v is struct{}.
func IsEmptyStruct[T any]() bool {
	return reflect.TypeFor[T]() == reflect.TypeOf(struct{}{})
}
