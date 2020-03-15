// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

import "math/rand"

// BytesCrypto returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func BytesMath(n int) []byte {
	b := make([]byte, n)
	rand.Int()
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		for i := range b {
			b[i] = byte(seededRandMath.Int() & 0xFF)
		}
	}

	return b
}
