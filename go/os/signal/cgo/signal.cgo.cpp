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

int CGOSignalHandlerSignalAction(int signum) {
  return searking::SignalHandler::SignalAction(signum);
}

void CGOSignalHandlerSetSignalDumpToFd(int fd) {
  searking::SignalHandler::SetSignalDumpToFd(fd);
}

void CGOSignalHandlerSetStacktraceDumpToFile(char* name) {
  searking::SignalHandler::SetStacktraceDumpToFile(name);
}

void CGOSignalHandlerDumpPreviousStacktrace() {
  searking::SignalHandler::DumpPreviousStacktrace();
}

// don't forget to free the string after finished using it
char* CGOPreviousStacktrace() {
  auto str = searking::SignalHandler::PreviousStacktrace();

  char* writable = static_cast<char*>(malloc((str.size() + 1) * sizeof(char)));
  std::copy(str.begin(), str.end(), writable);
  writable[str.size()] = '\0';  // don't forget the terminating 0
  return writable;
}
