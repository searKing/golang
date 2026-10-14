// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo && unix

#include "signal_handler_unix.hpp"

#include <string.h>

#include <cstdlib>

namespace searking {

SignalHandler &SignalHandler::GetInstance() {
  static SignalHandler instance;
  return instance;
}
// https://github.com/boostorg/stacktrace/blob/5c6740b68067cbd7070d2965bfbce32e81f680c9/example/terminate_handler.cpp
void SignalHandler::operator()(int signum, siginfo_t *info, void *context) {
  WriteSignalStacktrace(signum);

  InvokeGoSignalHandler(signum, info, context);
}

void SignalHandler::InvokeGoSignalHandler(int signum, siginfo_t *info,
                                          void *context) {
  auto handlers = go_registered_handlers_.Load(signum);
  if (handlers != nullptr) {
    SignalHandlerSigActionHandler sigActionHandler = handlers->first;
    SignalHandlerSignalHandler signalHandler = handlers->second;

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

void SignalHandler::SetGoRegisteredSignalHandlersIfEmpty(
    int signum, SignalHandlerSigActionHandler action,
    SignalHandlerSignalHandler handler) {
  // register once, avoid go's signal actions are lost.
  if (go_registered_handlers_.Load(signum) == nullptr) {
    go_registered_handlers_.Store(
        signum, new std::pair<SignalHandlerSigActionHandler,
                              SignalHandlerSignalHandler>(action, handler));
  }
}

// for CGO

int SignalHandler::SetSig(int signum) {
  SignalHandlerSigActionHandler sa_sigaction_action = nullptr;
  SignalHandlerSignalHandler sa_sigaction_handler = nullptr;
  sa_sigaction_action = [](int signum, siginfo_t *info, void *context) {
    GetInstance()(signum, info, context);
  };

  return SetSig(signum, sa_sigaction_action, sa_sigaction_handler);
}

int SignalHandler::SetSig(int signum, SignalHandlerSigActionHandler action,
                          SignalHandlerSignalHandler handler) {
  // The alternate signal stack is per thread, and Go installs one for every
  // thread of its own, which must not be replaced.
  stack_t ss;
  if (sigaltstack(nullptr, &ss) == 0 && (ss.ss_flags & SS_DISABLE)) {
    ss.ss_size = SIGSTKSZ * 100;
    ss.ss_sp = malloc(ss.ss_size);
    ss.ss_flags = 0;
    if (ss.ss_sp == nullptr) {
      return EXIT_FAILURE;
    }
    if (sigaltstack(&ss, nullptr) == -1) {
      free(ss.ss_sp);
      return EXIT_FAILURE;
    }
  }
  struct sigaction sa;
  memset(&sa, 0, sizeof(sa));
  sigaction(signum, nullptr, &sa);
  //  sigemptyset(&sa.sa_mask);
  //  sigfillset(&sa.sa_mask);
  if (sa.sa_flags & SA_SIGINFO) {
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
