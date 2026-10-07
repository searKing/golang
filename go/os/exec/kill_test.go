// Copyright 2026 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux || darwin

package exec_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	exec_ "github.com/searKing/golang/go/os/exec"
)

func TestKillProcByName(t *testing.T) {
	if _, err := exec.LookPath("killall"); err != nil {
		t.Skip("killall not found")
	}
	sleep, err := exec.LookPath("sleep")
	if err != nil {
		t.Skip("sleep not found")
	}
	data, err := os.ReadFile(sleep)
	if err != nil {
		t.Fatal(err)
	}
	// a unique name, so that no other processes are killed
	name := "kpbn" + strconv.Itoa(os.Getpid())
	bin := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(bin, data, 0755); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin, "30")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	exec_.KillProcByName(name)
	select {
	case err := <-done:
		if err == nil {
			t.Errorf("process exited normally, want killed")
		}
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Kill()
		t.Errorf("process not killed by name %q", name)
	}
}
