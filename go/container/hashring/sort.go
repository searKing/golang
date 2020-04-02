// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring

import "sort"

// Convenience types for common cases

// IntSlice attaches the methods of Interface to []int, sorting in increasing order.
type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p Uint32Slice) Sort() { sort.Sort(p) }
