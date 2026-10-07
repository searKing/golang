// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_test

import (
	"testing"

	os_ "github.com/searKing/golang/go/os"
)

func TestDiskUsage(t *testing.T) {
	total, free, avail, inodes, inodesFree, err := os_.DiskUsage(t.TempDir())
	if err != nil {
		t.Fatalf("DiskUsage() error = %v", err)
	}
	if total <= 0 || avail < 0 || avail > total || free < avail {
		t.Errorf("DiskUsage() total = %d, free = %d, avail = %d, want 0 <= avail <= free and avail <= total, total > 0", total, free, avail)
	}
	if inodes < 0 || inodesFree < 0 || inodesFree > inodes {
		t.Errorf("DiskUsage() inodes = %d, inodesFree = %d, want 0 <= inodesFree <= inodes", inodes, inodesFree)
	}

	if _, _, _, _, _, err := os_.DiskUsage("/not/exist/path"); err == nil {
		t.Errorf("DiskUsage() of not exist path error = nil, want an error")
	}
}
