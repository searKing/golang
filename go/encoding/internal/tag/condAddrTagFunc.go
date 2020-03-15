// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import "reflect"

// If CanAddr then get addr and handle else handle directly
type condAddrTagFunc struct {
	canAddrTagFunc, elseTagFunc tagFunc
}

func (ce *condAddrTagFunc) handle(e *tagState, v reflect.Value, opts tagOpts) (isUserDefined bool) {
	if v.CanAddr() {
		return ce.canAddrTagFunc(e, v, opts)
	}
	return ce.elseTagFunc(e, v, opts)
}

// newCondAddrConverter returns an encoder that checks whether its structTag
// CanAddr and delegates to canAddrTagFunc if so, else to elseTagFunc.
func newCondAddrTagFunc(canAddrConvert, elseConvert tagFunc) tagFunc {
	tagFn := &condAddrTagFunc{canAddrTagFunc: canAddrConvert, elseTagFunc: elseConvert}
	return tagFn.handle
}
