// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// js and wasip1 do not support inter-process file locking.
//
//go:build !js && !wasip1

package filelock_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/searKing/golang/go/sync/filelock"
)

func mustTempDir(t *testing.T) (dir string, remove func()) {
	t.Helper()

	dir, err := os.MkdirTemp("", filepath.Base(t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	return dir, func() { os.RemoveAll(dir) }
}

func mustBlockFunc(t *testing.T, desc string, f func()) (wait func(*testing.T)) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()

	timer := time.NewTimer(quiescent)
	defer timer.Stop()
	select {
	case <-done:
		t.Fatalf("%s unexpectedly did not block", desc)
	case <-timer.C:
	}

	return func(t *testing.T) {
		logTimer := time.NewTimer(quiescent)
		defer logTimer.Stop()

		select {
		case <-logTimer.C:
			// We expect the operation to have unblocked by now,
			// but maybe it's just slow. Write to the test log
			// in case the test times out, but don't fail it.
			t.Helper()
			t.Logf("%s is unexpectedly still blocked after %v", desc, quiescent)

			// Wait for the operation to actually complete, no matter how long it
			// takes. If the test has deadlocked, this will cause the test to time out
			// and dump goroutines.
			<-done

		case <-done:
		}
	}
}

func TestMutexExcludes(t *testing.T) {
	t.Parallel()

	dir, remove := mustTempDir(t)
	defer remove()

	path := filepath.Join(dir, "lock")

	mu := filelock.MutexAt(path)
	t.Logf("mu := MutexAt(_)")

	unlock, err := mu.Lock()
	if err != nil {
		t.Fatalf("mu.Lock: %v", err)
	}
	t.Logf("unlock, _  := mu.Lock()")

	mu2 := filelock.MutexAt(mu.Path)
	t.Logf("mu2 := MutexAt(mu.Path)")

	wait := mustBlockFunc(t, "mu2.Lock()", func() {
		unlock2, err := mu2.Lock()
		if err != nil {
			t.Errorf("mu2.Lock: %v", err)
			return
		}
		t.Logf("unlock2, _ := mu2.Lock()")
		t.Logf("unlock2()")
		unlock2()
	})

	t.Logf("unlock()")
	unlock()
	wait(t)
}

func TestReadWaitsForLock(t *testing.T) {
	t.Parallel()

	dir, remove := mustTempDir(t)
	defer remove()

	path := filepath.Join(dir, "timestamp.txt")

	f, err := filelock.Create(path)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	defer f.Close()

	const (
		part1 = "part 1\n"
		part2 = "part 2\n"
	)
	_, err = f.File.WriteString(part1)
	if err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	t.Logf("WriteString(%q) = <nil>", part1)

	wait := mustBlockFunc(t, "Read", func() {
		b, err := filelock.Read(path)
		if err != nil {
			t.Errorf("Read: %v", err)
			return
		}

		const want = part1 + part2
		got := string(b)
		if got == want {
			t.Logf("Read(_) = %q", got)
		} else {
			t.Errorf("Read(_) = %q, _; want %q", got, want)
		}
	})

	_, err = f.File.WriteString(part2)
	if err != nil {
		t.Errorf("WriteString: %v", err)
	} else {
		t.Logf("WriteString(%q) = <nil>", part2)
	}
	f.Close()

	wait(t)
}

func TestCanLockExistingFile(t *testing.T) {
	t.Parallel()

	dir, remove := mustTempDir(t)
	defer remove()
	path := filepath.Join(dir, "existing.txt")

	if err := os.WriteFile(path, []byte("ok"), 0777); err != nil {
		t.Fatalf("os.WriteFile: %v", err)
	}

	f, err := filelock.Edit(path)
	if err != nil {
		t.Fatalf("first Edit: %v", err)
	}

	wait := mustBlockFunc(t, "Edit", func() {
		other, err := filelock.Edit(path)
		if err != nil {
			t.Errorf("second Edit: %v", err)
		}
		other.Close()
	})

	f.Close()
	wait(t)
}
