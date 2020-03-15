// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import "reflect"

func IsNilObject(obtained interface{}) (result bool) {
	if obtained == nil {
		result = true
	} else {
		return IsNilValue(reflect.ValueOf(obtained))
	}
	return
}
