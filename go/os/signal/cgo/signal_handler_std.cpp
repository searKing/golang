// +build non-go
/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */

#include "signal_handler_std.hpp"

#include <string.h>

#include <boost/stacktrace.hpp>
#include <fstream>
#include <memory>
#include <sstream>

#include "base_signal_handler.hpp"

namespace searking {

SignalHandlerStd &SignalHandlerStd::GetInstance() {
  static SignalHandlerStd instance;
  return instance;
}

// https://github.com/boostorg/stacktrace/blob/5c6740b68067cbd7070d2965bfbce32e81f680c9/example/terminate_handler.cpp
void SignalHandlerStd::operator()(int signum) {
  WriteSignalStacktrace(signum);

  void *on_signal_ctx = on_signal_ctx_;
  auto on_signal = on_signal_;

  if (on_signal) {
    on_signal(on_signal_ctx, signal_dump_to_fd_, signum);
  }

  auto it = go_registered_handlers_.find(signum);
  if (it != go_registered_handlers_.end()) {
    SignalHandlerStdSignalHandler handler = it->second;

    if (handler == SIG_IGN) {
      return;
    }
    if (handler == SIG_DFL) {
      ::signal(signum, SIG_DFL);
      ::raise(signum);
      return;
    }
    handler(signum);
  }
}

void SignalHandlerStd::RegisterOnSignal(
    std::function<void(void *ctx, int fd, int signum)> callback, void *ctx) {
  on_signal_ctx_ = ctx;
  on_signal_ = callback;
}

void SignalHandlerStd::SetGoRegisteredSignalHandlersIfEmpty(
    int signum, SignalHandlerStdSignalHandler handler) {
  auto it = go_registered_handlers_.find(signum);

  // register once, avoid go's signal actions are lost.
  if (it == go_registered_handlers_.end()) {
    go_registered_handlers_[signum] = handler;
  }
}

SignalHandlerStdSignalHandler SignalHandlerStd::Signal(int signum) {
  SignalHandlerStdSignalHandler handler = [](int signum) {
    SignalHandlerStdSignalHandler prev_handler = ::signal(signum, SIG_DFL);
    GetInstance()(signum);
    ::signal(signum, prev_handler);
  };

  return Signal(signum, handler);
}

SignalHandlerStdSignalHandler SignalHandlerStd::Signal(
    int signum, SignalHandlerStdSignalHandler handler) {
  SignalHandlerStdSignalHandler prev_handler = ::signal(signum, SIG_DFL);
  SetGoRegisteredSignalHandlersIfEmpty(signum, prev_handler);

  return ::signal(signum, handler);
}

}  // namespace searking
