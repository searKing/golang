// +build cgo

/*
 * Copyright (c) 2019 The searKing authors. All Rights Reserved.
 *
 * Use of this source code is governed by a MIT-style license
 * that can be found in the LICENSE file in the root of the source
 * tree. An additional intellectual property rights grant can be found
 * in the file PATENTS.  All contributing project authors may
 * be found in the AUTHORS file in the root of the source tree.
 */

package internal
/*
	#cgo CXXFLAGS: -I${SRCDIR}/include/
	#cgo windows CXXFLAGS: -g -DUSE_WINDOWS_SIGNAL_HANDLER
	#cgo darwin CXXFLAGS: -g -D_GNU_SOURCE -DUSE_UNIX_SIGNAL_HANDLER
	#cgo !windows,!darwin CXXFLAGS: -g -DUSE_UNIX_SIGNAL_HANDLER
	#cgo linux LDFLAGS: -ldl

	#include "raise.cgo.h"
	#include <stdio.h>
	#include <stdbool.h>
   	#include <stdlib.h>  // Needed for C.free
*/
import "C"

// MustSegmentFault must send a SIGSEGV from cgo
func MustSegmentFault(){
	C.MustSegmentFault()
}
