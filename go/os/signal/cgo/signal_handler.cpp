/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#include "signal_handler.hpp"

#if defined(USE_UNIX_SIGNAL_HANDLER)
#include "signal_handler_unix.hpp"
#elif defined(USE_WINDOWS_SIGNAL_HANDLER)
#include "signal_handler_windows.hpp"
#else
#include "signal_handler_std.hpp"
#endif
namespace searking {

int SignalHandler::SetSig(int signum) {
#if defined(USE_UNIX_SIGNAL_HANDLER)
  // Yes it is a UNIX because __unix__ is defined.
  return SignalHandlerUnix::GetInstance().SetSig(signum);
#elif defined(USE_WINDOWS_SIGNAL_HANDLER)
#else
  SignalHandlerStd::GetInstance().Signal(signum);
  return 0;
#endif
}

void SignalHandler::SetSignalDumpToFd(int fd) {
#if defined(USE_UNIX_SIGNAL_HANDLER)
  // Yes it is a UNIX because __unix__ is defined.
  SignalHandlerUnix::GetInstance().SetSignalDumpToFd(fd);
  return;
#elif defined(USE_WINDOWS_SIGNAL_HANDLER)
#else
  SignalHandlerStd::GetInstance().SetSignalDumpToFd(fd);
  return;
#endif
}

void SignalHandler::SetStacktraceDumpToFile(char *name) {
#if defined(USE_UNIX_SIGNAL_HANDLER)
  // Yes it is a UNIX because __unix__ is defined.
  SignalHandlerUnix::GetInstance().SetStacktraceDumpToFile(name);
#elif defined(USE_WINDOWS_SIGNAL_HANDLER)
#else
  SignalHandlerStd::GetInstance().SetStacktraceDumpToFile(name);
  return;
#endif
  return;
}

void SignalHandler::RegisterOnSignal(
    std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                       void *context)>
        callback,
    void *ctx) {
#if defined(USE_UNIX_SIGNAL_HANDLER)
  // Yes it is a UNIX because __unix__ is defined.
  SignalHandlerUnix::GetInstance().RegisterOnSignal(callback, ctx);
#elif defined(USE_WINDOWS_SIGNAL_HANDLER)
#else
  SignalHandlerStd::GetInstance().RegisterOnSignal(callback, ctx);
  return;
#endif
  return;
}

void SignalHandler::DumpPreviousStacktrace() {
#if defined(USE_UNIX_SIGNAL_HANDLER)
  // Yes it is a UNIX because __unix__ is defined.
  SignalHandlerUnix::GetInstance().DumpPreviousStacktrace();
#elif defined(USE_WINDOWS_SIGNAL_HANDLER)
#else
  SignalHandlerStd::GetInstance().DumpPreviousStacktrace();
#endif
  return;
}
std::string SignalHandler::PreviousStacktrace() {
#if defined(USE_UNIX_SIGNAL_HANDLER)
  // Yes it is a UNIX because __unix__ is defined.
  return SignalHandlerUnix::GetInstance().PreviousStacktrace();
#elif defined(USE_WINDOWS_SIGNAL_HANDLER)
#else
  return SignalHandlerStd::GetInstance().PreviousStacktrace();
#endif
}
}  // namespace searking
