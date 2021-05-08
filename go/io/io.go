// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"os"
)

// Stater is the interface that wraps the basic Stat method.
// Stat returns the FileInfo structure describing file.
type Stater interface {
	Stat() (os.FileInfo, error)
}
