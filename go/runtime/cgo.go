// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

// BytesPointer returns a pointer which can be used by CGO
func BytesPointer(data []byte) unsafe.Pointer {
	if len(data) == 0 {
		return nil
	}
	return unsafe.Pointer(&data[0])
}
