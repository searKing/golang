// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo && windows

#include "signal_handler_windows.hpp"

namespace searking {

SignalHandler &SignalHandler::GetInstance() {
  static SignalHandler instance;
  return instance;
}

// trigger once, and user must register this signum again
void SignalHandler::operator()(int signum) {
  WriteSignalStacktrace(signum);

  InvokeGoSignalHandler(signum);
}

void SignalHandler::InvokeGoSignalHandler(int signum) {
  auto handler = go_registered_handlers_.Load(signum);
  if (handler != nullptr) {
    ::signal(signum, *handler);
    ::raise(signum);
  }
}

void SignalHandler::SetGoRegisteredSignalHandlersIfEmpty(
    int signum, SignalHandlerSignalHandler handler) {
  // register once, avoid go's signal actions are lost.
  if (go_registered_handlers_.Load(signum) == nullptr) {
    go_registered_handlers_.Store(signum,
                                  new SignalHandlerSignalHandler(handler));
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
