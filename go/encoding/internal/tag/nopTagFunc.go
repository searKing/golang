// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import "reflect"

type nopTagFunc struct {
}

func (_ *nopTagFunc) handle(e *tagState, v reflect.Value, opts tagOpts) (isUserDefined bool) {
	// nop
	return false
}

func newNopConverter(t reflect.Type) tagFunc {
	tagFn := &nopTagFunc{}
	return tagFn.handle
}
