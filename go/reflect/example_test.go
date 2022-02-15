// Copyright 2022 The searKing Author. All rights reserved.
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

func TestTruncate(t *testing.T) {
	type Human struct {
		Name       string
		Desc       []byte
		Friends    []Human
		FriendById map[string][]Human
	}

	var info = Human{
		Name:       "ALPHA",
		Desc:       []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		Friends:    []Human{{Name: "BETA", Desc: []byte("abcdefghijklmnopqrstuvwxyz")}},
		FriendById: map[string][]Human{"quick brown fox": {{Name: "GRAMMAR", Desc: []byte("The quick brown fox jumps over the lazy dog")}}},
	}
	TruncateBytes(info, 3)
	TruncateString(info, 3)
	fmt.Printf("info truncated\n")
	fmt.Printf("info.Name: %s\n", info.Name)
	fmt.Printf("info.Desc: %s\n", info.Desc)

	for i, friend := range info.Friends {
		fmt.Printf("info.Friends[%d].Name: %s\n", i, friend.Name)
		fmt.Printf("info.Friends[%d].Desc: %s\n", i, friend.Desc)
	}
	for id, friends := range info.FriendById {
		for i, friend := range friends {
			fmt.Printf("info.FriendById[%s][%d].Name: %s\n", id, i, friend.Name)
			fmt.Printf("info.FriendById[%s][%d].Desc: %s\n", id, i, friend.Desc)
		}
	}
	// Output:
	// info truncated
	// info.Name: size: 5, string: ALP
	// info.Desc: size: 26, bytes: ABC
	// info.Friends[0].Name: size: 4, string: BET
	// info.Friends[0].Desc: size: 26, bytes: abc
	// info.FriendById[quick brown fox][0].Name: size: 7, string: GRA
	// info.FriendById[quick brown fox][0].Desc: size: 43, bytes: The

}
