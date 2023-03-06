// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson

import _ "unsafe"

//go:linkname foldFunc encoding/json.foldFunc
func foldFunc(s []byte) func(s, t []byte) bool
