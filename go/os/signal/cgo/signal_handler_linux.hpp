// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

// +build cgo
// +build linux

#ifndef GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_LINUX_HPP_
#define GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_LINUX_HPP_

#include <unistd.h>
// You can find out the version with _POSIX_VERSION.
// POSIX compliant

#include <csignal>
#include <functional>
#include <map>
#include <memory>
#include <utility>

#include "base_signal_handler.hpp"

namespace searking {

// Callbacks Predefinations

typedef void (*SignalHandlerSigActionHandler)(int signum, siginfo_t *info,
                                              void *context);
typedef void (*SignalHandlerSignalHandler)(int signum);
typedef void (*SignalHandlerOnSignal)(void *ctx, int fd, int signum,
                                      siginfo_t *info, void *context);

class SignalHandler : public BaseSignalHandler<SignalHandler> {
 protected:
  SignalHandler() : on_signal_ctx_(nullptr), on_signal_(nullptr) {}

  void SetGoRegisteredSignalHandlersIfEmpty(
      int signum, SignalHandlerSigActionHandler action,
      SignalHandlerSignalHandler handler);
  void DoSignalChan(int signum, siginfo_t *info, void *context);
  void InvokeGoSignalHandler(int signum, siginfo_t *info, void *context);

 public:
  // Thread safe GetInstance.
  static SignalHandler &GetInstance();

  void operator()(int signum, siginfo_t *info, void *context);
  // never invoke a go function, see
  // https://github.com/golang/go/issues/35814
  void RegisterOnSignal(std::function<void(void *ctx, int fd, int signum,
                                           siginfo_t *info, void *context)>
                            callback,
                        void *ctx);

  static int SetSig(int signum);
  static int SetSig(int signum, SignalHandlerSigActionHandler action,
                    SignalHandlerSignalHandler handler);

 private:
  void *on_signal_ctx_;
  std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                     void *context)>
      on_signal_;
  std::map<int, std::pair<SignalHandlerSigActionHandler,
                          SignalHandlerSignalHandler> >
      go_registered_handlers_;
};
}  // namespace searking
#endif
