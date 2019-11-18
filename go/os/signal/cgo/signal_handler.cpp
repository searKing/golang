/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#include "signal_handler.hpp"

#include <string.h>

#include <boost/stacktrace.hpp>
#include <fstream>
#include <memory>
#include <sstream>

#include "write_int.hpp"

namespace searking {

// https://github.com/boostorg/stacktrace/blob/5c6740b68067cbd7070d2965bfbce32e81f680c9/example/terminate_handler.cpp
void SignalHandler::operator()(int signum, siginfo_t *info, void *context) {
  if (signal_dump_to_fd_ >= 0) {
    write(signal_dump_to_fd_, "Signal received(", strlen("Signal received("));
    WriteInt(signal_dump_to_fd_, signum);
    write(signal_dump_to_fd_, ").\n", strlen(").\n"));
    // binary format,not human readable.mute this.
    //    write(signal_dump_to_fd_, "stacktrace dumped in binary format:\n",
    //          strlen("stacktrace dumped in binary format:\n"));
    //    boost::stacktrace::safe_dump_to(signal_dump_to_fd_);
    //    write(signal_dump_to_fd_, "\n", strlen("\n"));
  }

  if (!stacktrace_dump_to_file_.empty()) {
    if (signal_dump_to_fd_ >= 0) {
      write(signal_dump_to_fd_,
            "Backtrace dumped to file: ", strlen("Backtrace dumped to file: "));
      write(signal_dump_to_fd_, stacktrace_dump_to_file_.c_str(),
            stacktrace_dump_to_file_.length());
      write(signal_dump_to_fd_, ".\n", strlen(".\n"));
    }
    boost::stacktrace::safe_dump_to(stacktrace_dump_to_file_.c_str());
  }

  auto it = cgo_sigaction_handlers_.find(signum);
  if (it != cgo_sigaction_handlers_.end()) {
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

  void *on_signal_ctx = on_signal_ctx_;
  auto on_signal = on_signal_;

  if (on_signal) {
    on_signal(on_signal_ctx, signal_dump_to_fd_, signum, info, context);
  }
}

void SignalHandler::RegisterOnSignal(
    std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                       void *context)>
        callback,
    void *ctx) {
  std::lock_guard<std::mutex> lock(mutex_);
  on_signal_ctx_ = ctx;
  on_signal_ = callback;
}

void SignalHandler::SetSigactionHandlers(int signum,
                                         SIGNAL_SA_ACTION_CALLBACK action,
                                         SIGNAL_SA_HANDLER_CALLBACK handler) {
  std::lock_guard<std::mutex> lock(mutex_);
  auto it = cgo_sigaction_handlers_.find(signum);

  // register once, avoid go's signal actions are lost.
  if (it == cgo_sigaction_handlers_.end()) {
    cgo_sigaction_handlers_[signum] = std::make_pair(action, handler);
  }
}

void SignalHandler::SetSignalDumpToFd(int fd) {
  std::lock_guard<std::mutex> lock(mutex_);
  signal_dump_to_fd_ = fd;
}

void SignalHandler::SetSignalDumpToFd(FILE *fd) {
  SetSignalDumpToFd(fileno(fd));
}

void SignalHandler::SetStacktraceDumpToFile(const std::string &name) {
  std::lock_guard<std::mutex> lock(mutex_);
  stacktrace_dump_to_file_ = name;
}

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

ssize_t SignalHandler::DumpPreviousHumanReadableStacktrace() {
  std::ostringstream msg;
  msg << "Previous run crashed:" << std::endl;

  msg << SignalHandler::PreviousHumanReadableStacktrace();
  auto fd = GetInstance().signal_dump_to_fd_;
  if (fd < 0) {
    return 0;
  }

  auto m = msg.str();
  return write(fd, m.c_str(), m.length());
}

std::string SignalHandler::PreviousHumanReadableStacktrace() {
  auto name = GetInstance().stacktrace_dump_to_file_;
  if (name.empty()) {
    return "";
  }
  std::ifstream ifs(name);
  if (!ifs.good()) {
    return "";
  }

  std::shared_ptr<int> deferFileCLose(nullptr, [&ifs](int *) {
    // cleaning up
    ifs.close();
  });

  // there is a backtrace
  boost::stacktrace::stacktrace st =
      boost::stacktrace::stacktrace::from_dump(ifs);
  std::ostringstream msg;

  msg << st << std::endl;
  return msg.str();
}
}  // namespace searking
