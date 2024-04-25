// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo

#ifndef GO_OS_SIGNAL_CGO_BASE_SIGNAL_HANDLER_HPP_
#define GO_OS_SIGNAL_CGO_BASE_SIGNAL_HANDLER_HPP_
#include <boost/core/noncopyable.hpp>
#include <boost/stacktrace.hpp>
#include <cstdio>
#include <cstring>
#include <fstream>
#include <map>
#include <memory>
#include <sstream>
#include <string>
#include <tuple>

#include "write_int.hpp"

namespace searking {
template <class T>
class BaseSignalHandler : private boost::noncopyable,
                          public std::enable_shared_from_this<T> {
 protected:
  BaseSignalHandler() : signal_dump_to_fd_(-1) {}

 public:
  void SetSignalDumpToFd(int fd) { signal_dump_to_fd_ = fd; }

  void SetSignalDumpToFd(FILE *fd) { SetSignalDumpToFd(fileno(fd)); }

  void SetStacktraceDumpToFile(const std::string &name) {
    stacktrace_dump_to_file_ = name;
  }

  void WriteSignalStacktrace(int signum) {
    if (signal_dump_to_fd_ >= 0) {
      (void)!write(signal_dump_to_fd_, "Signal received(",
                   strlen("Signal received("));
      (void)!WriteInt(signal_dump_to_fd_, signum);
      (void)!write(signal_dump_to_fd_, ").\n", strlen(").\n"));
      // binary format,not human readable.mute this.
      //    write(signal_dump_to_fd_, "stacktrace dumped in binary format:\n",
      //          strlen("stacktrace dumped in binary format:\n"));
      //    boost::stacktrace::safe_dump_to(signal_dump_to_fd_);
      //    write(signal_dump_to_fd_, "\n", strlen("\n"));
    }

    if (!stacktrace_dump_to_file_.empty()) {
      if (signal_dump_to_fd_ >= 0) {
        (void)!write(signal_dump_to_fd_, "Stacktrace dumped to file: ",
                     strlen("Stacktrace dumped to file: "));
        (void)!write(signal_dump_to_fd_, stacktrace_dump_to_file_.c_str(),
                     stacktrace_dump_to_file_.length());
        (void)!write(signal_dump_to_fd_, ".\n", strlen(".\n"));
      }
      boost::stacktrace::safe_dump_to(stacktrace_dump_to_file_.c_str());
    }
  }

  ssize_t DumpPreviousStacktrace() {
    if (signal_dump_to_fd_ < 0) {
      return 0;
    }

    std::ostringstream msg;
    msg << "Previous run crashed:" << std::endl;
    msg << PreviousStacktrace();

    auto m = msg.str();
    return write(signal_dump_to_fd_, m.c_str(), m.length());
  }
  std::string PreviousStacktrace() {
    if (stacktrace_dump_to_file_.empty()) {
      return "";
    }

    std::ifstream ifs(stacktrace_dump_to_file_);
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

  void SetSigInvokeChain(const int from, const int to, const int wait,
                         const int sleepInSeconds) {
    sig_invoke_signal_chains_[from] =
        std::make_tuple(from, to, wait, sleepInSeconds);
  }

  void SetSigInvokeChain(const int from, const int pipeWriter,
                         const int pipeReader) {
    sig_invoke_pipe_chains_[from] =
        std::make_tuple(from, pipeWriter, pipeReader);
  }

 protected:
  int signal_dump_to_fd_;
  std::string stacktrace_dump_to_file_;
  // <from, <from, to, wait, sleepInSeconds>>
  std::map<int, std::tuple<int, int, int, int>> sig_invoke_signal_chains_;

  // <from, <from, pipeWriter, pipeReader>>
  std::map<int, std::tuple<int, int, int>> sig_invoke_pipe_chains_;
};
}  // namespace searking
#endif  // GO_OS_SIGNAL_CGO_BASE_SIGNAL_HANDLER_HPP_
