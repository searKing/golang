/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#ifndef GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_UNIX_HPP_
#define GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_UNIX_HPP_

#if defined(USE_UNIX_SIGNAL_HANDLER)

#include <unistd.h>
// You can find out the version with _POSIX_VERSION.
// POSIX compliant

#include <csignal>
#include <functional>
#include <map>
#include <utility>

#include "base_signal_handler.hpp"

namespace searking {

// Callbacks Predefinations

typedef void (*SignalHandlerSigActionHandler)(int signum, siginfo_t *info,
                                              void *context);
typedef void (*SignalHandlerSignalHandler)(int signum);

class SignalHandlerUnix : public BaseSignalHandler {
 protected:
  SignalHandlerUnix() : on_signal_ctx_(nullptr), on_signal_(nullptr) {}

  void SetGoRegisteredSignalHandlersIfEmpty(
      int signum, SignalHandlerSigActionHandler action,
      SignalHandlerSignalHandler handler);

 public:
  // Thread safe GetInstance.
  static SignalHandlerUnix &GetInstance();

  void operator()(int signum, siginfo_t *info, void *context);
  void RegisterOnSignal(std::function<void(void *ctx, int fd, int signum,
                                           siginfo_t *info, void *context)>
                            callback,
                        void *ctx);

  static int SignalAction(int signum);
  static int SignalAction(int signum, SignalHandlerSigActionHandler action,
                          SignalHandlerSignalHandler handler);

 private:
  void *on_signal_ctx_;
  std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                     void *context)>
      on_signal_;
  std::map<int, std::pair<SignalHandlerSigActionHandler,
                          SignalHandlerSignalHandler> >
      go_registered_handlers_;

 private:
  SignalHandlerUnix(const SignalHandlerUnix &) = delete;
  void operator=(const SignalHandlerUnix &) = delete;
};
}  // namespace searking
#endif
#endif  // GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_UNIX_HPP_
