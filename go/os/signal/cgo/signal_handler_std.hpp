/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#ifndef GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_STD_HPP__
#define GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_STD_HPP__

#if defined(USE_STD_SIGNAL_HANDLER)

#include <csignal>
#include <functional>
#include <map>
#include <utility>

#include "base_signal_handler.hpp"

namespace searking {

// Callbacks Predefinations

typedef void (*SignalHandlerStdSignalHandler)(int signum);

class SignalHandlerStd : public BaseSignalHandler {
 protected:
  SignalHandlerStd() : on_signal_ctx_(nullptr), on_signal_(nullptr) {}

  void SetGoRegisteredSignalHandlersIfEmpty(
      int signum, SignalHandlerStdSignalHandler handler);

 public:
  // Thread safe GetInstance.
  static SignalHandlerStd &GetInstance();
  void operator()(int signum);
  void RegisterOnSignal(
      std::function<void(void *ctx, int fd, int signum)> callback, void *ctx);

  SignalHandlerStdSignalHandler Signal(int signum);
  SignalHandlerStdSignalHandler Signal(int signum,
                                       SignalHandlerStdSignalHandler handler);

 private:
  void *on_signal_ctx_;
  std::function<void(void *ctx, int fd, int signum)> on_signal_;
  std::map<int, SignalHandlerStdSignalHandler> go_registered_handlers_;

 private:
  SignalHandlerStd(const SignalHandlerStd &) = delete;
  void operator=(const SignalHandlerStd &) = delete;
};
}  // namespace searking
#endif
#endif  // GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_STD_HPP__
