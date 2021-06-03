// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

// +build cgo
// +build linux

#include "signal_handler_linux.hpp"

#include <string.h>

#include <boost/stacktrace.hpp>
#include <fstream>
#include <memory>
#include <sstream>

namespace searking {
// sig nums must be in [0,255)
static volatile sig_atomic_t gotSignals[256];

SignalHandler &SignalHandler::GetInstance() {
  static SignalHandler instance;
  return instance;
}
// https://github.com/boostorg/stacktrace/blob/5c6740b68067cbd7070d2965bfbce32e81f680c9/example/terminate_handler.cpp
void SignalHandler::operator()(int signum, siginfo_t *info, void *context) {
  WriteSignalStacktrace(signum);

  void *on_signal_ctx = on_signal_ctx_;
  auto on_signal = on_signal_;
  if (on_signal) {
    on_signal(on_signal_ctx, signal_dump_to_fd_, signum, info, context);
  }

  DoSignalChan(signum, info, context);

  InvokeGoSignalHandler(signum, info, context);
}

void SignalHandler::DoSignalChan(int signum, siginfo_t *info, void *context) {
  gotSignals[signum] = true;
  auto it = sig_invoke_signal_chains_.find(signum);
  if (it == sig_invoke_signal_chains_.end()) {
    return;
  }
  auto &sig_chain = it->second;
  int from = std::get<0>(sig_chain);
  // consist validation
  if (from != signum) {
    return;
  }
  int to = std::get<1>(sig_chain);
  int wait = std::get<2>(sig_chain);
  int sleepInSeconds = std::get<3>(sig_chain);

  if (to >= 0 && to != signum) {
    InvokeGoSignalHandler(to, info, context);
  }

  // I don't know why suspend will block forever sometimes, so use an ugly
  // implement as a loop, the same as in windows
  if (wait >= 0 && wait != signum) {
    gotSignals[wait] = false;
    for (;;) {
      bool got = gotSignals[wait];
      if (got) {
        gotSignals[wait] = false;
        break;
      }
      // sleep 1s at most, will awake when an unmasked signal is received
      sleep(1);
    }
  }

  if (sleepInSeconds > 0) {
    sleep(sleepInSeconds);
  }
}
void SignalHandler::InvokeGoSignalHandler(int signum, siginfo_t *info,
                                          void *context) {
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

void SignalHandler::RegisterOnSignal(
    std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                       void *context)>
        callback,
    void *ctx) {
  on_signal_ctx_ = ctx;
  on_signal_ = callback;
}

void SignalHandler::SetGoRegisteredSignalHandlersIfEmpty(
    int signum, SignalHandlerSigActionHandler action,
    SignalHandlerSignalHandler handler) {
  auto it = go_registered_handlers_.find(signum);

  // register once, avoid go's signal actions are lost.
  if (it == go_registered_handlers_.end()) {
    go_registered_handlers_[signum] = std::make_pair(action, handler);
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
  stack_t ss;
  sigaltstack(NULL, &ss);
  ss.ss_sp = malloc(SIGSTKSZ * 100);
  ss.ss_size = SIGSTKSZ * 100;
  ss.ss_flags = 0;
  if (sigaltstack(&ss, NULL) == -1) {
    return EXIT_FAILURE;
  }
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
