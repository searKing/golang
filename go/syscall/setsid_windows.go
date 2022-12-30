// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec

import (
	"syscall"

	"golang.org/x/sys/windows"
)

// SysProcAttrSetsid run a program in a new session, is used to detach the process from the parent (normally a shell)
func SysProcAttrSetsid(attr *syscall.SysProcAttr) {
	// CREATE_NEW_PROCESS_GROUP: The new process is the root process of a new process group. The process group includes all processes that are descendants of this root process.
	// 	The process identifier of the new process group is the same as the process identifier,
	// 	which is returned in the lpProcessInformation parameter.
	//  If this flag is specified, CTRL+C signals will be disabled for all processes within the new process group.
	// DETACHED_PROCESS: For console processes, the new process does not inherit its parent's console (the default).
	// https://learn.microsoft.com/en-us/windows/win32/procthread/process-creation-flags
	attr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP | windows.DETACHED_PROCESS
	attr.HideWindow = true
}
