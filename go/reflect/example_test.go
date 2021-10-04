// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflect

import (
	"fmt"
	"os"
	osexec "os/exec"
	"testing"
)

func TestGetppid(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		fmt.Print(os.Getppid())
		os.Exit(0)
	}

	cmd := osexec.Command(os.Args[0], "-test.run=TestGetppid")
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")

	// verify that Getppid() from the forked process reports our process id
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to spawn child process: %v %q", err, string(output))
	}

	childPpid := string(output)
	ourPid := fmt.Sprintf("%d", os.Getpid())
	if childPpid != ourPid {
		t.Fatalf("Child process reports parent process id '%v', expected '%v'", childPpid, ourPid)
	}
}

func TestTruncated(t *testing.T) {
	var info = struct {
		Name string
		Desc []byte
	}{
		Name: "ALPHA",
		Desc: []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	}
	TruncateBytes(&info, 3)
	TruncateString(&info, 3)
	fmt.Printf("info truncated\n")
	fmt.Printf("info.Name: %s\n", info.Name)
	fmt.Printf("info.Desc: %s\n", info.Desc)
	// Output:
	// info truncated
	// info.Name: size: 5, string: ALP
	// info.Desc: size: 26, bytes: ABC

}
