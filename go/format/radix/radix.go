// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package radix

//go:generate go-enum -type Radix
type Radix int // 进制

const (
	Binary      Radix = iota // 二进制
	Octonary                 // 八进制
	Decimal                  // 十进制
	Hexadecimal              // 十六进制
)
