// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import "os/exec"

func KillProcByName(pname string) {
	params := []string{
		"taskkill",
		"/F",
		"/IM",
		pname,
		"/T",
	}
	exec.Command(params[0], params[1:]...).CombinedOutput()
}
