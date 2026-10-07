// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux || darwin

package exec

import "os/exec"

// KillProcByName kills processes named pname by killall with SIGKILL, ignoring errors.
// killall on Solaris and AIX kills all processes, so it is available on Linux and macOS only.
func KillProcByName(pname string) {
	params := []string{
		"killall",
		"-9",
		pname,
	}
	exec.Command(params[0], params[1:]...).CombinedOutput()
}
