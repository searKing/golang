// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import "os"

// FileModeTimeSlice sorts filenames by mode time in increase order
type FileModeTimeSlice []string

func (s FileModeTimeSlice) Len() int {
	return len(s)
}
func (s FileModeTimeSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FileModeTimeSlice) Less(i, j int) bool {
	fi, err := os.Stat(s[i])
	if err != nil {
		return false
	}
	fj, err := os.Stat(s[j])
	if err != nil {
		return false
	}
	return fi.ModTime().Before(fj.ModTime())
}

// FileModeTimeDescSlice sorts filenames by mode time in decrease order
type FileModeTimeDescSlice []string

func (s FileModeTimeDescSlice) Len() int {
	return len(s)
}
func (s FileModeTimeDescSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FileModeTimeDescSlice) Less(i, j int) bool {
	fi, err := os.Stat(s[i])
	if err != nil {
		return false
	}
	fj, err := os.Stat(s[j])
	if err != nil {
		return false
	}
	return fi.ModTime().After(fj.ModTime())
}
