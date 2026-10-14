// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo && windows

#ifndef GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_WINDOWS_HPP__
#define GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_WINDOWS_HPP__

#include <csignal>
#include <mutex>

#include "base_signal_handler.hpp"

namespace searking {

typedef void (*SignalHandlerSignalHandler)(int signum);
// Never used this Handler in cgo, for go needs
class SignalHandler : public BaseSignalHandler {
 protected:
  SignalHandler() = default;

 public:
  // Thread safe GetInstance.
  static SignalHandler &GetInstance();
  void operator()(int signum);

  static int SetSig(int signum);

 private:
  static void HandleSignal(int signum);
  void InvokeGoSignalHandler(int signum);

  std::mutex install_mutex_;
  // SIG_DFL is nullptr, so handlers are stored by pointers to tell from unset.
  SignalTable<SignalHandlerSignalHandler> go_registered_handlers_;
};
}  // namespace searking
#endif  // GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_WINDOWS_HPP__
