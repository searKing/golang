// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

// +build cgo
// +build windows

#include "signal_handler_windows.hpp"

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

// trigger once, and user must register this signum again
void SignalHandler::operator()(int signum) {
  WriteSignalStacktrace(signum);
  void *on_signal_ctx = on_signal_ctx_;
  auto on_signal = on_signal_;

  if (on_signal) {
    on_signal(on_signal_ctx, signal_dump_to_fd_, signum);
  }
  DoSignalChan(signum);

  InvokeGoSignalHandler(signum);
}

void SignalHandler::DoSignalChan(int signum) {
  gotSignals[signum] = true;
  auto it = sig_invoke_signal_chains_.find(signum);
  if (it == sig_invoke_signal_chains_.end()) {
    return;
  }
  auto &sig_chain = it->second;
  int from = std::get<0>(sig_chain);
  // consist validation_
  if (from != signum) {
    return;
  }
  int to = std::get<1>(sig_chain);
  int wait = std::get<2>(sig_chain);
  int sleepInSeconds = std::get<3>(sig_chain);
  if (to >= 0 && to != signum) {
    InvokeGoSignalHandler(to);
  }
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

void SignalHandler::InvokeGoSignalHandler(int signum) {
  auto it = go_registered_handlers_.find(signum);
  if (it != go_registered_handlers_.end()) {
    SignalHandlerSignalHandler signalHandler = it->second;
    ::signal(signum, signalHandler);
    ::raise(signum);
  }
}

void SignalHandler::RegisterOnSignal(
    std::function<void(void *ctx, int fd, int signum)> callback, void *ctx) {
  on_signal_ctx_ = ctx;
  on_signal_ = callback;
}

void SignalHandler::SetGoRegisteredSignalHandlersIfEmpty(
    int signum, SignalHandlerSignalHandler handler) {
  auto it = go_registered_handlers_.find(signum);

  // register once, avoid go's signal actions are lost.
  if (it == go_registered_handlers_.end()) {
    go_registered_handlers_[signum] = handler;
  }
}

int SignalHandler::SetSig(int signum) {
  SignalHandlerSignalHandler handler = [](int signum) {
    GetInstance()(signum);
  };

  return SetSig(signum, handler);
}

int SignalHandler::SetSig(int signum, SignalHandlerSignalHandler handler) {
  SignalHandlerSignalHandler prev_handler = ::signal(signum, SIG_DFL);

  if (SIG_ERR == prev_handler) {
    return -1;
  }
  GetInstance().SetGoRegisteredSignalHandlersIfEmpty(signum, prev_handler);

  ::signal(signum, handler);
  return 0;
}

}  // namespace searking
