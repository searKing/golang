/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */

#include "base_signal_handler.hpp"

#include <boost/stacktrace.hpp>
#include <cstring>
#include <fstream>
#include <memory>
#include <sstream>

#include "write_int.hpp"

namespace searking {

void BaseSignalHandler::SetSignalDumpToFd(int fd) { signal_dump_to_fd_ = fd; }

void BaseSignalHandler::SetSignalDumpToFd(FILE *fd) {
  SetSignalDumpToFd(fileno(fd));
}

void BaseSignalHandler::SetStacktraceDumpToFile(const std::string &name) {
  stacktrace_dump_to_file_ = name;
}

void BaseSignalHandler::WriteSignalStacktrace(int signum) {
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
      write(signal_dump_to_fd_, "Stacktrace dumped to file: ",
            strlen("Stacktrace dumped to file: "));
      write(signal_dump_to_fd_, stacktrace_dump_to_file_.c_str(),
            stacktrace_dump_to_file_.length());
      write(signal_dump_to_fd_, ".\n", strlen(".\n"));
    }
    boost::stacktrace::safe_dump_to(stacktrace_dump_to_file_.c_str());
  }
}

ssize_t BaseSignalHandler::DumpPreviousStacktrace() {
  if (signal_dump_to_fd_ < 0) {
    return 0;
  }

  std::ostringstream msg;
  msg << "Previous run crashed:" << std::endl;
  msg << PreviousStacktrace();

  auto m = msg.str();
  return write(signal_dump_to_fd_, m.c_str(), m.length());
}

std::string BaseSignalHandler::PreviousStacktrace() {
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

void BaseSignalHandler::SetSigInvokeChain(const int from, const int to,
                                          const int wait,
                                          const int sleepInSeconds) {
  sig_invoke_chains_[from] = {from, to, wait, sleepInSeconds};
}

}  // namespace searking
