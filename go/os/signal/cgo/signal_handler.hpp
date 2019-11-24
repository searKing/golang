/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#ifndef GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_HPP_
#define GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_HPP_
#include <string>
namespace searking {
class SignalHandler {
 public:
  static int SetSig(int signum);
  static void SetSignalDumpToFd(int fd);
  static void SetStacktraceDumpToFile(char *name);
  static void RegisterOnSignal(
      std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                         void *context)>
          callback,
      void *ctx);

  static void DumpPreviousStacktrace();
  static std::string PreviousStacktrace();
};
}  // namespace searking
#endif  // GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_HPP_
