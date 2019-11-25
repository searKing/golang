/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#ifndef GO_OS_SIGNAL_CGO_SIGNAL_CGO_H_
#define GO_OS_SIGNAL_CGO_SIGNAL_CGO_H_
#include <signal.h>
#include <stdbool.h>
#ifdef __cplusplus
extern "C" {
#endif

// Callbacks Predefinations

typedef void (*CGO_SignalHandlerSigActionHandler)(void *ctx, int fd, int signum,
                                                  siginfo_t *info,
                                                  void *context);
int CGO_SignalHandlerSetSig(int signum);
void CGO_SignalHandlerSetSignalDumpToFd(int fd);
void CGO_SignalHandlerSetStacktraceDumpToFile(char *name);
void CGO_SignalHandlerRegisterOnSignal(
    CGO_SignalHandlerSigActionHandler callback, void *ctx);
void CGO_SignalHandlerDumpPreviousStacktrace();
char *CGO_PreviousStacktrace();
void CGO_SetSigInvokeChain(const int from, const int to, const int wait,
                           const int sleepInSeconds);

#ifdef __cplusplus
}
#endif

#endif  // GO_OS_SIGNAL_CGO_SIGNAL_CGO_H_
