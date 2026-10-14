// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo && windows

#include "signal_handler_windows.hpp"

#include <string.h>

#include <boost/stacktrace.hpp>
#include <fstream>
#include <memory>
#include <sstream>

namespace searking {

SignalHandler &SignalHandler::GetInstance() {
  static SignalHandler instance;
  return instance;
}

// trigger once, and user must register this signum again
void SignalHandler::operator()(int signum) {
  WriteSignalStacktrace(signum);

  const OnSignal *on_signal = on_signal_.load(std::memory_order_acquire);
  if (on_signal != nullptr && on_signal->callback) {
    on_signal->callback(on_signal->ctx, SignalDumpToFd(), signum);
  }

  RunSigChain(signum, [this](int to) { InvokeGoSignalHandler(to); });

  InvokeGoSignalHandler(signum);
}

void SignalHandler::InvokeGoSignalHandler(int signum) {
  auto handler = go_registered_handlers_.Load(signum);
  if (handler != nullptr) {
    ::signal(signum, *handler);
    ::raise(signum);
  }
}

void SignalHandler::RegisterOnSignal(
    std::function<void(void *ctx, int fd, int signum)> callback, void *ctx) {
  on_signal_.store(new OnSignal{callback, ctx}, std::memory_order_release);
}

void SignalHandler::SetGoRegisteredSignalHandlersIfEmpty(
    int signum, SignalHandlerSignalHandler handler) {
  // register once, avoid go's signal actions are lost.
  go_registered_handlers_.StoreIfEmpty(signum, handler);
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
