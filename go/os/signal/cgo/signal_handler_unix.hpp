// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo && unix

#ifndef GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_UNIX_HPP_
#define GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_UNIX_HPP_

#include <unistd.h>
// You can find out the version with _POSIX_VERSION.
// POSIX compliant

#include <csignal>
#include <mutex>

#include "base_signal_handler.hpp"

namespace searking {

struct SignalDisposition {
  explicit SignalDisposition(const struct sigaction &action) : action(action) {}
  const struct sigaction action;
};

class SignalHandler : public BaseSignalHandler {
 protected:
  SignalHandler() = default;

  void InvokeGoSignalHandler(int signum, siginfo_t *info, void *context);

 public:
  // Thread safe GetInstance.
  static SignalHandler &GetInstance();

  void operator()(int signum, siginfo_t *info, void *context);

  static int SetSig(int signum);

 private:
  static void HandleSignal(int signum, siginfo_t *info, void *context);

  std::mutex install_mutex_;
  SignalTable<SignalDisposition> go_registered_handlers_;
};
}  // namespace searking
#endif
