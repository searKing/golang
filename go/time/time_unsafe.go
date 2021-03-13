// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time

import (
	_ "unsafe"
)

//go:linkname nextStdChunk time.nextStdChunk
func nextStdChunk(layout string) (prefix string, std int, suffix string)
