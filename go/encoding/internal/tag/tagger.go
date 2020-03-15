// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import "reflect"

func valueTaggeFunc(v reflect.Value) tagFunc {
	if !v.IsValid() {
		return invalidValueTagFunc
	}
	return typeTagFunc(v.Type())
}
