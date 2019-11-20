/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#ifndef GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_HPP_
#define GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_HPP_
#include <atomic>
#include <csignal>
#include <cstdio>
#include <functional>
#include <map>
#include <mutex>
#include <string>
#include <utility>

namespace searking {

// Callbacks Predefinations
typedef void (*SIGNAL_SA_ACTION_CALLBACK)(int signum, siginfo_t *info,
                                          void *context);
typedef void (*SIGNAL_SA_HANDLER_CALLBACK)(int signum);

typedef void (*SIGNAL_ON_SIGNAL_CALLBACK)(void *ctx, int fd, int signum,
                                          siginfo_t *info, void *context);

typedef void (*SIGNAL_ON_BACKTRACE_DUMP_CALLBACK)(int fd);

class SignalHandler {
 protected:
  SignalHandler()
      : on_signal_ctx_(nullptr), on_signal_(nullptr), signal_dump_to_fd_(-1) {}

 public:
  // Thread safe GetInstance.
  static SignalHandler &GetInstance();

  void operator()(int signum, siginfo_t *info, void *context);
  void RegisterOnSignal(std::function<void(void *ctx, int fd, int signum,
                                           siginfo_t *info, void *context)>
                            callback,
                        void *ctx);

  void SetSignalDumpToFd(int fd);

  void SetSignalDumpToFd(FILE *fd);

  void SetStacktraceDumpToFile(const std::string &name);

  void SetSigactionHandlers(int signum, SIGNAL_SA_ACTION_CALLBACK action,
                            SIGNAL_SA_HANDLER_CALLBACK handler);

  static int SignalAction(int signum);
  static int SignalAction(int signum, SIGNAL_SA_ACTION_CALLBACK action,
                          SIGNAL_SA_HANDLER_CALLBACK handler);
  static ssize_t DumpPreviousHumanReadableStacktrace();
  static std::string PreviousHumanReadableStacktrace();

 private:
  std::mutex mutex_;
  int signal_dump_to_fd_;
  std::string stacktrace_dump_to_file_;

  void *on_signal_ctx_;
  std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                     void *context)>
      on_signal_;
  std::map<int,
           std::pair<SIGNAL_SA_ACTION_CALLBACK, SIGNAL_SA_HANDLER_CALLBACK> >
      cgo_sigaction_handlers_;

 private:
  SignalHandler(const SignalHandler &) = delete;
  void operator=(const SignalHandler &) = delete;
};
}  // namespace searking
#endif  // GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_HPP_
