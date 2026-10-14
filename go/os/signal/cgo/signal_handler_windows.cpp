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

void SignalHandler::HandleSignal(int signum) { GetInstance()(signum); }

void SignalHandler::InvokeGoSignalHandler(int signum) {
  const SignalHandlerSignalHandler *handler =
      go_registered_handlers_.Load(signum);
  ::signal(signum, handler == nullptr ? SIG_DFL : *handler);
  ::raise(signum);
}

int SignalHandler::SetSig(int signum) {
  SignalHandler &instance = GetInstance();
  std::lock_guard<std::mutex> lock(instance.install_mutex_);

  const SignalHandlerSignalHandler *saved =
      instance.go_registered_handlers_.Load(signum);
  instance.go_registered_handlers_.Store(signum, nullptr);

  SignalHandlerSignalHandler previous = ::signal(signum, &HandleSignal);
  if (previous == SIG_ERR) {
    instance.go_registered_handlers_.Store(signum, saved);
    return -1;
  }
  if (previous == &HandleSignal) {
    instance.go_registered_handlers_.Store(signum, saved);
  } else {
    instance.go_registered_handlers_.Store(
        signum, new SignalHandlerSignalHandler(previous));
  }
  return 0;
}

}  // namespace searking
