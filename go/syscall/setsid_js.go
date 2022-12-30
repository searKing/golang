// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import "syscall"

// SysProcAttrSetsid run a program in a new session, is used to detach the process from the parent (normally a shell)
//
// The disowning of a child process is accomplished by executing the system call
// setpgrp() or setsid(), (both of which have the same functionality) as soon as
// the child is forked. These calls create a new process session group, make the
// child process the session leader, and set the process group ID to the process
// ID of the child. https://bsdmag.org/unix-kernel-system-calls/
func SysProcAttrSetsid(attr *syscall.SysProcAttr) {

}
