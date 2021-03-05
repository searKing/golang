// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sort

import (
	"sort"
	"time"
)

// Convenience types for common cases

// TimeSlice attaches the methods of Interface to []time.Time, sorting in increasing order.
type TimeSlice []time.Time

func (x TimeSlice) Len() int           { return len(x) }
func (x TimeSlice) Less(i, j int) bool { return x[i].Before(x[j]) }
func (x TimeSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x TimeSlice) Sort() { sort.Sort(x) }

// TimeDurationSlice attaches the methods of Interface to []time.Duration, sorting in increasing order.
type TimeDurationSlice []time.Duration

func (x TimeDurationSlice) Len() int           { return len(x) }
func (x TimeDurationSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x TimeDurationSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x TimeDurationSlice) Sort() { sort.Sort(x) }
