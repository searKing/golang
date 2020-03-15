// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tag

import "reflect"

func invalidValueTagFunc(_ *tagState, _ reflect.Value, _ tagOpts) (isUserDefined bool) { return false }
