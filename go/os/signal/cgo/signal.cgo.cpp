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

#include <iostream>

#include "signal_handler.h"
#include "stacktrace.h"
int CGOSignalHandlerSignalAction(int signum) {
  return searking::SignalHandler::SignalAction(signum);
}
void CGOSignalHandlerSetFd(int fd) {
  searking::SignalHandler::GetInstance().SetFd(fd);
}

void CGOSignalHandlerSetBacktraceDump(bool enable) {
  if (enable) {
    searking::SignalHandler::GetInstance().SetBacktraceDumpTo(
        searking::stacktrace::SafeDumpToFd);
    return;
  }
  searking::SignalHandler::GetInstance().SetBacktraceDumpTo(nullptr);
}
