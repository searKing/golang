// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http2

import (
	"bytes"
	"io"

	"golang.org/x/net/http2"
)

var (
	clientPreface = []byte(http2.ClientPreface)
)

func HasClientPreface(r io.Reader) bool {
	// Check the validity of client preface.
	preface := make([]byte, len(clientPreface))
	if _, err := io.ReadFull(r, preface); err != nil {
		return false
	}
	return bytes.Equal(preface, clientPreface)
}
