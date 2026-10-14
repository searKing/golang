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
#include <atomic>
#include <boost/core/noncopyable.hpp>
#include <boost/stacktrace.hpp>
#include <cstdio>
#include <cstring>
#include <fstream>
#include <sstream>
#include <string>

#include "write_int.hpp"

// Signal handlers read the state by atomics, which must not lock.
static_assert(ATOMIC_POINTER_LOCK_FREE == 2,
              "atomic pointers must be lock free");
static_assert(ATOMIC_INT_LOCK_FREE == 2, "atomic ints must be lock free");

namespace searking {

// Signal numbers must be in [0, kMaxSignum).
constexpr int kMaxSignum = 256;

// SignalTable maps signal numbers to immutable values, which are read in
// signal handlers while set concurrently.
// Values replaced are never freed, as signal handlers may be reading them.
template <class V>
class SignalTable : private boost::noncopyable {
 public:
  SignalTable() {
    for (auto &entry : entries_) {
      entry.store(nullptr, std::memory_order_relaxed);
    }
  }

  const V *Load(int signum) const {
    if (signum < 0 || signum >= kMaxSignum) {
      return nullptr;
    }
    return entries_[signum].load(std::memory_order_acquire);
  }

  // Store stores v, which must be allocated by new, and is never freed.
  void Store(int signum, const V *v) {
    if (signum < 0 || signum >= kMaxSignum) {
      return;
    }
    entries_[signum].store(v, std::memory_order_release);
  }

 private:
  std::atomic<const V *> entries_[kMaxSignum];
};

struct DumpPath {
  explicit DumpPath(const std::string &name)
      : length(name.size()), path(new char[length + 1]) {
    memcpy(path, name.data(), length);
    path[length] = '\0';
  }

  const size_t length;
  char *const path;
};

class BaseSignalHandler : private boost::noncopyable {
 protected:
  BaseSignalHandler()
      : signal_dump_to_fd_(-1), stacktrace_dump_to_file_(nullptr) {}

 public:
  void SetSignalDumpToFd(int fd) {
    signal_dump_to_fd_.store(fd, std::memory_order_release);
  }

  void SetSignalDumpToFd(FILE *fd) { SetSignalDumpToFd(fileno(fd)); }

  // The previous name is never freed, as signal handlers may be reading it.
  void SetStacktraceDumpToFile(const std::string &name) {
    stacktrace_dump_to_file_.store(new DumpPath(name),
                                   std::memory_order_release);
  }

  void WriteSignalStacktrace(int signum) {
    int fd = SignalDumpToFd();
    const DumpPath *file = StacktraceDumpToFile();
    if (fd >= 0) {
      (void)!write(fd, "Signal received(", strlen("Signal received("));
      (void)!WriteInt(fd, signum);
      (void)!write(fd, ").\n", strlen(").\n"));
      // binary format,not human readable.mute this.
      //    write(fd, "stacktrace dumped in binary format:\n",
      //          strlen("stacktrace dumped in binary format:\n"));
      //    boost::stacktrace::safe_dump_to(fd);
      //    write(fd, "\n", strlen("\n"));
    }

    if (file != nullptr && file->length != 0) {
      if (fd >= 0) {
        (void)!write(fd, "Stacktrace dumped to file: ",
                     strlen("Stacktrace dumped to file: "));
        (void)!write(fd, file->path, file->length);
        (void)!write(fd, ".\n", strlen(".\n"));
      }
      boost::stacktrace::safe_dump_to(file->path);
    }
  }

  ssize_t DumpPreviousStacktrace() {
    int fd = SignalDumpToFd();
    if (fd < 0) {
      return 0;
    }

    std::ostringstream msg;
    msg << "Previous run crashed:" << std::endl;
    msg << PreviousStacktrace();

    auto m = msg.str();
    return write(fd, m.c_str(), m.length());
  }

  std::string PreviousStacktrace() {
    const DumpPath *file = StacktraceDumpToFile();
    if (file == nullptr || file->length == 0) {
      return "";
    }

    std::ifstream ifs(file->path);
    if (!ifs.good()) {
      return "";
    }

    // there is a backtrace
    boost::stacktrace::stacktrace st =
        boost::stacktrace::stacktrace::from_dump(ifs);
    std::ostringstream msg;

    msg << st << std::endl;
    return msg.str();
  }

 protected:
  int SignalDumpToFd() const {
    return signal_dump_to_fd_.load(std::memory_order_acquire);
  }

  const DumpPath *StacktraceDumpToFile() const {
    return stacktrace_dump_to_file_.load(std::memory_order_acquire);
  }

 private:
  std::atomic<int> signal_dump_to_fd_;
  std::atomic<const DumpPath *> stacktrace_dump_to_file_;
};
}  // namespace searking
#endif  // GO_OS_SIGNAL_CGO_BASE_SIGNAL_HANDLER_HPP_
