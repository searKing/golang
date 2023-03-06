// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prettyjson

import (
	"math"
	"strconv"
)

func appendFloat(out []byte, n float64, bitSize int) []byte {
	switch {
	case math.IsNaN(n):
		return append(out, "nan"...)
	case math.IsInf(n, +1):
		return append(out, "inf"...)
	case math.IsInf(n, -1):
		return append(out, "-inf"...)
	default:
		return strconv.AppendFloat(out, n, 'g', -1, bitSize)
	}
}
