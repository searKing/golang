/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#ifndef GO_OS_SIGNAL_CGO_BASE_SIGNAL_HANDLER_HPP_
#define GO_OS_SIGNAL_CGO_BASE_SIGNAL_HANDLER_HPP_
#include <cstdio>
#include <string>
namespace searking {
class BaseSignalHandler {
 protected:
  BaseSignalHandler() : signal_dump_to_fd_(-1) {}

 public:
  void SetSignalDumpToFd(int fd);

  void SetSignalDumpToFd(FILE *fd);

  void SetStacktraceDumpToFile(const std::string &name);

  void WriteSignalStacktrace(int signum);

  ssize_t DumpPreviousStacktrace();
  std::string PreviousStacktrace();

 protected:
  int signal_dump_to_fd_;
  std::string stacktrace_dump_to_file_;
};
}  // namespace searking
#endif  // GO_OS_SIGNAL_CGO_BASE_SIGNAL_HANDLER_HPP_
