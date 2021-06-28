// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/searKing/golang/tools/go-import/testdata/cgo/include/has_go"

const Name = "string"

func main() {
	has_go.HasGo()
}
