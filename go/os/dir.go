// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"os"
	"sort"
)

// ReadDirN reads the named directory,
// returning a slice of its directory entries sorted by filename.
// If an error occurs reading the directory,
// ReadDirN returns a slice of entries it was able to read before the error,
// along with the error.
//
// If n > 0, ReadDirN returns at most n DirEntry records.
// In this case, if ReadDirN returns an empty slice, it will return an error explaining why.
// At the end of a directory, the error is [io.EOF].
//
// If n <= 0, ReadDirN returns all the DirEntry records remaining in the directory.
// When it succeeds, it returns a nil error (not [io.EOF]).
// To read all dirs, see [os.ReadDir], or set n as -1.
func ReadDirN(name string, n int) ([]os.DirEntry, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dirs, err := f.ReadDir(n)
	if len(dirs) > 1 {
		sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name() < dirs[j].Name() })
	}
	return dirs, err
}
