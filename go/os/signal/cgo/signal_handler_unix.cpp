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

void SignalHandler::HandleSignal(int signum, siginfo_t *info, void *context) {
  GetInstance()(signum, info, context);
}

void SignalHandler::InvokeGoSignalHandler(int signum, siginfo_t *info,
                                          void *context) {
  const SignalDisposition *disposition =
      go_registered_handlers_.Load(signum);
  if (disposition == nullptr) {
    return;
  }

  const struct sigaction &action = disposition->action;
  // sa_handler and sa_sigaction may share storage. Inspect only the member
  // selected by SA_SIGINFO.
  if ((action.sa_flags & SA_SIGINFO) != 0) {
    using SigAction = void (*)(int, siginfo_t *, void *);
    const SigAction handler = action.sa_sigaction;
    if (handler == reinterpret_cast<SigAction>(SIG_IGN)) {
      return;
    }
    if (handler == reinterpret_cast<SigAction>(SIG_DFL)) {
      if (sigaction(signum, &action, nullptr) == 0) {
        raise(signum);
      }
      return;
    }
    if (handler != &SignalHandler::HandleSignal) {
      handler(signum, info, context);
    }
    return;
  }
  if (action.sa_handler == SIG_IGN) {
    return;
  }
  if (action.sa_handler == SIG_DFL) {
    if (sigaction(signum, &action, nullptr) == 0) {
      raise(signum);
    }
    return;
  }
  if (action.sa_handler != nullptr) {
    action.sa_handler(signum);
  }
}

// for CGO

int SignalHandler::SetSig(int signum) {
  SignalHandler &instance = GetInstance();
  std::lock_guard<std::mutex> lock(instance.install_mutex_);

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

  struct sigaction previous;
  memset(&previous, 0, sizeof(previous));
  if (sigaction(signum, nullptr, &previous) != 0) {
    return -1;
  }
  if ((previous.sa_flags & SA_SIGINFO) != 0 &&
      previous.sa_sigaction == &SignalHandler::HandleSignal) {
    return 0;
  }

  instance.go_registered_handlers_.Store(
      signum, new SignalDisposition(previous));

  struct sigaction wrapper = previous;
  wrapper.sa_flags |= SA_SIGINFO | SA_ONSTACK | SA_RESTART;
  wrapper.sa_sigaction = &SignalHandler::HandleSignal;
  return sigaction(signum, &wrapper, nullptr);
}

}  // namespace searking
