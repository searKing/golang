// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo

#ifndef GO_OS_SIGNAL_CGO_SIGNAL_CGO_H_
#define GO_OS_SIGNAL_CGO_SIGNAL_CGO_H_
#include <signal.h>
#include <stdbool.h>
#include <stdint.h>
#ifdef __cplusplus
extern "C" {
#endif

// Callbacks Predefinations
int CGO_SignalHandlerSetSig(int signum);
void CGO_SignalHandlerSetSignalDumpToFd(int fd);
void CGO_SignalHandlerSetStacktraceDumpToFile(char *name);
void CGO_SignalHandlerDumpPreviousStacktrace();
char *CGO_PreviousStacktrace();
void CGO_SetSigInvokeChain(const int from, const int to, const int wait,
                           const int sleepInSeconds);
// CGO_NewSigInvokeChain returns the chain for CGO_CancelSigInvokeChain.
uintptr_t CGO_NewSigInvokeChain(const int from, const int to, const int wait,
                                const int sleepInSeconds);
void CGO_CancelSigInvokeChain(const int from, uintptr_t chain);

#ifdef __cplusplus
}
#endif

#endif  // GO_OS_SIGNAL_CGO_SIGNAL_CGO_H_
