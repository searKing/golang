// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import "io"

type eofReader struct{}

func (eofReader) Read([]byte) (int, error) {
	return 0, io.EOF
}

// EOFReader returns a Reader that return EOF anytime.
func EOFReader() io.Reader {
	return eofReader{}
}
