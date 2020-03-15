// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

// Bytes returns securely generated random bytes.
func Bytes(n int) []byte {
	b, err := BytesCrypto(n)
	if err == nil {
		return b
	}
	return BytesMath(n)
}
