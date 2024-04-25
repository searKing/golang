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
#include <functional>
#include <map>
#include <utility>

#include "base_signal_handler.hpp"

namespace searking {

// Callbacks Predefinations

typedef void (*SignalHandlerSignalHandler)(int signum);
typedef void (*SignalHandlerOnSignal)(void *ctx, int fd, int signum);
// Never used this Handler in cgo, for go needs
class SignalHandler : public BaseSignalHandler<SignalHandler> {
 protected:
  SignalHandler() : on_signal_ctx_(nullptr), on_signal_(nullptr) {}

  void SetGoRegisteredSignalHandlersIfEmpty(int signum,
                                            SignalHandlerSignalHandler handler);
  void DoSignalChan(int signum);
  void InvokeGoSignalHandler(int signum);

 public:
  // Thread safe GetInstance.
  static SignalHandler &GetInstance();
  void operator()(int signum);
  // never invoke a go function, see
  // https://github.com/golang/go/issues/35814
  void RegisterOnSignal(
      std::function<void(void *ctx, int fd, int signum)> callback, void *ctx);

  static int SetSig(int signum);
  static int SetSig(int signum, SignalHandlerSignalHandler handler);

 private:
  void *on_signal_ctx_;
  std::function<void(void *ctx, int fd, int signum)> on_signal_;
  std::map<int, SignalHandlerSignalHandler> go_registered_handlers_;

 private:
  SignalHandler(const SignalHandler &) = delete;
  void operator=(const SignalHandler &) = delete;
};
}  // namespace searking
#endif  // GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_WINDOWS_HPP__
