// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package atomic

import (
	"fmt"
	"io/ioutil"
	"os"
)

// File is an atomic wrapper around a file.
type File string

func (m *File) TryLock() error {
	if m == nil {
		return fmt.Errorf("nil pointer")
	}
	if *m == "" {
		temp, err := ioutil.TempFile("", "file_lock_")
		if err != nil {
			return err
		}
		*m = File(temp.Name())
		_ = temp.Close()
		return nil
	}

	f, err := os.OpenFile(string(*m), os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		// Can't lock, just return
		return err
	}
	_ = f.Close()
	return nil
}

func (m *File) TryUnlock() error {
	if m == nil || *m == "" {
		return nil
	}
	return os.Remove(string(*m))
}
