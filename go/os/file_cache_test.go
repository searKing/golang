// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_test

import (
	"os"
	"strings"
	"testing"
	"time"

	os_ "github.com/searKing/golang/go/os"
)

func TestCacheFile_NeverExpireByDefault(t *testing.T) {
	f := os_.NewCacheFile(os_.WithCacheFileBucketRootDir(t.TempDir()))
	const key = "http://foo.com/kitty.jpg"

	path, refreshed, err := f.Put(key, strings.NewReader("data"))
	if err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if !refreshed {
		t.Errorf("Put() refreshed = false on first put, want true")
	}

	cachePath, _, hit, err := f.Get(key)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !hit || cachePath != path {
		t.Errorf("Get() = %q, hit %v, want %q, hit true", cachePath, hit, path)
	}
	if data, err := os.ReadFile(path); err != nil || string(data) != "data" {
		t.Errorf("cache file = %q, %v, want %q", data, err, "data")
	}

	if _, refreshed, err = f.Put(key, strings.NewReader("data")); err != nil || refreshed {
		t.Errorf("Put() again refreshed = %v, %v, want false, nil", refreshed, err)
	}
}

func TestCacheFile_Expire(t *testing.T) {
	f := os_.NewCacheFile(os_.WithCacheFileBucketRootDir(t.TempDir()),
		os_.WithCacheFileCacheExpiredAfter(10*time.Millisecond))
	const key = "http://foo.com/kitty.jpg"

	if _, _, err := f.Put(key, strings.NewReader("data")); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	_, _, hit, err := f.Get(key)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if hit {
		t.Errorf("Get() hit = true after expired, want false")
	}
}
