/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#include "signal_handler.h"

#include <string.h>

#include "write_int.h"

namespace searking {

// https://github.com/boostorg/stacktrace/blob/5c6740b68067cbd7070d2965bfbce32e81f680c9/example/terminate_handler.cpp
void SignalHandler::operator()(int signum, siginfo_t *info, void *context) {
  if (backtrace_dump_to_) {
    // https://stackoverflow.com/questions/16891019/how-to-avoid-using-printf-in-a-signal-handler
    write(fd_, "Sig(", strlen("Sig("));
    WriteInt(fd_, signum);
    write(fd_, ") Backtrace:\n", strlen(") Backtrace:\n"));
    backtrace_dump_to_(fd_);
    write(fd_, "Backtrace End\n", strlen("Backtrace End\n"));
  }

  auto it = sigactionHandlers_.find(signum);
  if (it != sigactionHandlers_.end()) {
    auto handlers = it->second;
    SIGNAL_SA_ACTION_CALLBACK sa_sigaction_action = handlers.first;
    SIGNAL_SA_HANDLER_CALLBACK sa_sigaction_handler = handlers.second;
    if (sa_sigaction_action) {
      sa_sigaction_action(signum, info, context);
    }
    if (sa_sigaction_handler) {
      sa_sigaction_handler(signum);
    }
  }

  void *onSignalCtx = onSignalCtx_;
  auto onSignal = onSignal_;

  if (onSignal) {
    onSignal(onSignalCtx, fd_, signum, info, context);
  }
}

void SignalHandler::RegisterOnSignal(
    std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                       void *context)>
        callback,
    void *ctx) {
  std::lock_guard<std::mutex> lock(mutex_);
  onSignalCtx_ = ctx;
  onSignal_ = callback;
}

void SignalHandler::SetSigactionHandlers(int signum,
                                         SIGNAL_SA_ACTION_CALLBACK action,
                                         SIGNAL_SA_HANDLER_CALLBACK handler) {
  std::lock_guard<std::mutex> lock(mutex_);
  sigactionHandlers_[signum] = std::make_pair(action, handler);
}

void SignalHandler::SetFd(int fd) {
  std::lock_guard<std::mutex> lock(mutex_);
  fd_ = fd;
}

void SignalHandler::SetBacktraceDumpTo(
    std::function<void(int fd)> safe_dump_to) {
  std::lock_guard<std::mutex> lock(mutex_);
  backtrace_dump_to_ = safe_dump_to;
}

void SignalHandler::SetFd(FILE *fd) { SetFd(fileno(fd)); }

SignalHandler &SignalHandler::GetInstance() {
  static SignalHandler instance;
  return instance;
}

int SignalHandler::SignalAction(int signum) {
  SIGNAL_SA_ACTION_CALLBACK sa_sigaction_action = nullptr;
  SIGNAL_SA_HANDLER_CALLBACK sa_sigaction_handler = nullptr;
  sa_sigaction_action = [](int signum, siginfo_t *info, void *context) {
    GetInstance()(signum, info, context);
  };

  return SignalAction(signum, sa_sigaction_action, sa_sigaction_handler);
}

int SignalHandler::SignalAction(int signum, SIGNAL_SA_ACTION_CALLBACK action,
                                SIGNAL_SA_HANDLER_CALLBACK handler) {
  struct sigaction sa;
  memset(&sa, 0, sizeof(sa));
  sigaction(signum, nullptr, &sa);
  //  sigemptyset(&sa.sa_mask);
  //  sigfillset(&sa.sa_mask);
  if (sa.sa_flags | SA_SIGINFO) {
    GetInstance().SetSigactionHandlers(signum, sa.sa_sigaction, nullptr);
  } else {
    GetInstance().SetSigactionHandlers(signum, nullptr, sa.sa_handler);
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
