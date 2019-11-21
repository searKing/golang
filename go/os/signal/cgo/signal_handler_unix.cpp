/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#include "signal_handler_unix.hpp"
#if defined(USE_UNIX_SIGNAL_HANDLER)

#include <string.h>

#include <boost/stacktrace.hpp>
#include <fstream>
#include <memory>
#include <sstream>

#include "base_signal_handler.hpp"

namespace searking {
SignalHandlerUnix &SignalHandlerUnix::GetInstance() {
  static SignalHandlerUnix instance;
  return instance;
}

// https://github.com/boostorg/stacktrace/blob/5c6740b68067cbd7070d2965bfbce32e81f680c9/example/terminate_handler.cpp
void SignalHandlerUnix::operator()(int signum, siginfo_t *info, void *context) {
  WriteSignalStacktrace(signum);

  void *on_signal_ctx = on_signal_ctx_;
  auto on_signal = on_signal_;

  if (on_signal) {
    on_signal(on_signal_ctx, signal_dump_to_fd_, signum, info, context);
  }

  auto it = go_registered_handlers_.find(signum);
  if (it != go_registered_handlers_.end()) {
    auto handlers = it->second;
    SignalHandlerSigActionHandler sigActionHandler = handlers.first;
    SignalHandlerSignalHandler signalHandler = handlers.second;

    // http://man7.org/linux/man-pages/man7/signal.7.html
    if (sigActionHandler) {
      sigActionHandler(signum, info, context);
      return;
    }
    if (signalHandler == SIG_IGN) {
      return;
    }
    if (signalHandler == SIG_DFL) {
      struct sigaction preSa;
      memset(&preSa, 0, sizeof(preSa));
      sigaction(signum, nullptr, &preSa);

      preSa.sa_sigaction = nullptr;
      preSa.sa_handler = SIG_DFL;

      sigaction(signum, &preSa, nullptr);
      raise(signum);
      return;
    }

    signalHandler(signum);
  }
}

void SignalHandlerUnix::RegisterOnSignal(
    std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                       void *context)>
        callback,
    void *ctx) {
  on_signal_ctx_ = ctx;
  on_signal_ = callback;
}

void SignalHandlerUnix::SetGoRegisteredSignalHandlersIfEmpty(
    int signum, SignalHandlerSigActionHandler action,
    SignalHandlerSignalHandler handler) {
  auto it = go_registered_handlers_.find(signum);

  // register once, avoid go's signal actions are lost.
  if (it == go_registered_handlers_.end()) {
    go_registered_handlers_[signum] = std::make_pair(action, handler);
  }
}

int SignalHandlerUnix::SignalAction(int signum) {
  SignalHandlerSigActionHandler sa_sigaction_action = nullptr;
  SignalHandlerSignalHandler sa_sigaction_handler = nullptr;
  sa_sigaction_action = [](int signum, siginfo_t *info, void *context) {
    GetInstance()(signum, info, context);
  };

  return SignalAction(signum, sa_sigaction_action, sa_sigaction_handler);
}

int SignalHandlerUnix::SignalAction(int signum,
                                    SignalHandlerSigActionHandler action,
                                    SignalHandlerSignalHandler handler) {
  struct sigaction sa;
  memset(&sa, 0, sizeof(sa));
  sigaction(signum, nullptr, &sa);
  //  sigemptyset(&sa.sa_mask);
  //  sigfillset(&sa.sa_mask);
  if (sa.sa_flags | SA_SIGINFO) {
    GetInstance().SetGoRegisteredSignalHandlersIfEmpty(signum, sa.sa_sigaction,
                                                       nullptr);
  } else {
    GetInstance().SetGoRegisteredSignalHandlersIfEmpty(signum, nullptr,
                                                       sa.sa_handler);
  }
  sa.sa_flags = sa.sa_flags & (~SA_SIGINFO);
  sa.sa_flags = sa.sa_flags | SA_ONSTACK | SA_RESTART;
  sa.sa_handler = nullptr;
  if (action) {
    // If SA_SIGINFO is specified in sa_flags, then sa_sigaction (instead of
    // sa_handler) specifies the signal-handling function for signum.  This
    // function receives three arguments, as described below.
    sa.sa_flags = sa.sa_flags | SA_SIGINFO;
    sa.sa_sigaction = action;
  } else if (handler) {
    sa.sa_handler = handler;
  }
  return sigaction(signum, &sa, nullptr);
}
}  // namespace searking
#endif
