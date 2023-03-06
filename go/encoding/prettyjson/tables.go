// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson

import "unicode/utf8"
import _ "unsafe"

//go:linkname htmlSafeSet encoding/json.htmlSafeSet
var htmlSafeSet [utf8.RuneSelf]bool

//go:linkname safeSet encoding/json.safeSet
var safeSet [utf8.RuneSelf]bool
