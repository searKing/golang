// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo && windows

#include <algorithm>

#include "signal.cgo.h"
#include "signal_handler_windows.hpp"

int CGO_SignalHandlerSetSig(int signum) {
  return searking::SignalHandler::SetSig(signum);
}

void CGO_SignalHandlerSetSignalDumpToFd(int fd) {
  searking::SignalHandler::GetInstance().SetSignalDumpToFd(fd);
}

void CGO_SignalHandlerSetStacktraceDumpToFile(char* name) {
  searking::SignalHandler::GetInstance().SetStacktraceDumpToFile(name);
}

void CGO_SignalHandlerDumpPreviousStacktrace() {
  searking::SignalHandler::GetInstance().DumpPreviousStacktrace();
}

// don't forget to free the string after finished using it
char* CGO_PreviousStacktrace() {
  auto str = searking::SignalHandler::GetInstance().PreviousStacktrace();

  char* writable = static_cast<char*>(malloc((str.size() + 1) * sizeof(char)));
  std::copy(str.begin(), str.end(), writable);
  writable[str.size()] = '\0';  // don't forget the terminating 0
  return writable;
}

void CGO_SetSigInvokeChain(const int from, const int to, const int wait,
                           const int sleepInSeconds) {
  searking::SignalHandler::GetInstance().SetSigInvokeChain(from, to, wait,
                                                           sleepInSeconds);
}
