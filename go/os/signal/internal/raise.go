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

	#include <signal.h>
	#include <stdio.h>
	#include <stdbool.h>
   	#include <stdlib.h>  // Needed for C.free

	int CallSignalRaise(int signum){
		return raise(signum);
	}

	int CallRaise(int signum){
		return CallSignalRaise(signum);
	}
*/
import "C"
import (
	"syscall"

	"github.com/searKing/golang/go/os/signal"
)

func Raise(sig syscall.Signal) int{
	return int(C.CallRaise(C.int(signal.Signum(sig))))
}
