// Copyright 2025 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package generic

import (
	"context"
)

//go:generate go-option -type "Number"
type Number[T context.Context] struct {
	GenericType T
}
