/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#include "signal.cgo.h"

#include <algorithm>

#include "signal_handler.hpp"

int CGO_SignalHandlerSetSig(int signum) {
  return searking::SignalHandler::SetSig(signum);
}

void CGO_SignalHandlerSetSignalDumpToFd(int fd) {
  searking::SignalHandler::SetSignalDumpToFd(fd);
}

void CGO_SignalHandlerSetStacktraceDumpToFile(char* name) {
  searking::SignalHandler::SetStacktraceDumpToFile(name);
}

void CGO_SignalHandlerRegisterOnSignal(
    CGO_SignalHandlerSigActionHandler callback, void* ctx) {
  searking::SignalHandler::RegisterOnSignal(callback, ctx);
}

void CGO_SignalHandlerDumpPreviousStacktrace() {
  searking::SignalHandler::DumpPreviousStacktrace();
}

// don't forget to free the string after finished using it
char* CGO_PreviousStacktrace() {
  auto str = searking::SignalHandler::PreviousStacktrace();

  char* writable = static_cast<char*>(malloc((str.size() + 1) * sizeof(char)));
  std::copy(str.begin(), str.end(), writable);
  writable[str.size()] = '\0';  // don't forget the terminating 0
  return writable;
}

void CGO_SetSigInvokeChain(const int from, const int to, const int wait,
                           const int sleepInSeconds) {
  searking::SignalHandler::SetSigInvokeChain(from, to, wait, sleepInSeconds);
}