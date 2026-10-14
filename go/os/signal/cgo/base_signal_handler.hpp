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
#include <csignal>
#include <cstdio>
#include <cstring>
#include <fstream>
#include <memory>
#include <sstream>
#include <string>
#include <time.h>

#include "write_int.hpp"

// Signal handlers read the state by atomics, which must not lock.
static_assert(ATOMIC_POINTER_LOCK_FREE == 2, "atomic pointers must be lock free");
static_assert(ATOMIC_INT_LOCK_FREE == 2, "atomic ints must be lock free");
static_assert(ATOMIC_BOOL_LOCK_FREE == 2, "atomic bools must be lock free");

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

  // RemoveIf removes the value of signum only if it is v, rather than one
  // stored after v.
  void RemoveIf(int signum, const V *v) {
    if (signum < 0 || signum >= kMaxSignum) {
      return;
    }
    const V *expected = v;
    entries_[signum].compare_exchange_strong(expected, nullptr,
                                             std::memory_order_acq_rel);
  }

  // StoreIfEmpty stores v only if no value of signum is stored.
  void StoreIfEmpty(int signum, const V &v) {
    if (signum < 0 || signum >= kMaxSignum) {
      return;
    }
    const V *expected = nullptr;
    const V *value = new V(v);
    if (!entries_[signum].compare_exchange_strong(expected, value,
                                                  std::memory_order_acq_rel)) {
      delete value;
    }
  }

 private:
  std::atomic<const V *> entries_[kMaxSignum];
};

// SigChain raises {to} when {from} is received, then waits for {wait} until
// canceled, and then sleeps {sleepInSeconds}s.
struct SigChain {
  SigChain(int to, int wait, int sleepInSeconds)
      : to(to), wait(wait), sleepInSeconds(sleepInSeconds), canceled(false) {}

  const int to;
  const int wait;
  const int sleepInSeconds;
  // canceled is set when the chain is removed, so that its wait returns.
  mutable std::atomic<bool> canceled;
};

template <class T>
class BaseSignalHandler : private boost::noncopyable,
                          public std::enable_shared_from_this<T> {
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
    stacktrace_dump_to_file_.store(new std::string(name),
                                   std::memory_order_release);
  }

  void WriteSignalStacktrace(int signum) {
    int fd = SignalDumpToFd();
    const std::string *file = StacktraceDumpToFile();
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

    if (file != nullptr && !file->empty()) {
      if (fd >= 0) {
        (void)!write(fd, "Stacktrace dumped to file: ",
                     strlen("Stacktrace dumped to file: "));
        (void)!write(fd, file->c_str(), file->length());
        (void)!write(fd, ".\n", strlen(".\n"));
      }
      boost::stacktrace::safe_dump_to(file->c_str());
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
    const std::string *file = StacktraceDumpToFile();
    if (file == nullptr || file->empty()) {
      return "";
    }

    std::ifstream ifs(*file);
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

  void SetSigInvokeChain(const int from, const int to, const int wait,
                         const int sleepInSeconds) {
    NewSigInvokeChain(from, to, wait, sleepInSeconds);
  }

  // NewSigInvokeChain sets a SigChain of {from}, and returns it for
  // CancelSigInvokeChain. It is never freed, as signal handlers may be
  // reading it.
  const SigChain *NewSigInvokeChain(const int from, const int to,
                                    const int wait, const int sleepInSeconds) {
    const SigChain *chain = new SigChain(to, wait, sleepInSeconds);
    sig_invoke_chains_.Store(from, chain);
    return chain;
  }

  // CancelSigInvokeChain removes chain of {from}, if not replaced yet, and
  // returns its wait in progress.
  void CancelSigInvokeChain(const int from, const SigChain *chain) {
    if (chain == nullptr) {
      return;
    }
    chain->canceled.store(true, std::memory_order_release);
    sig_invoke_chains_.RemoveIf(from, chain);
  }

 protected:
  int SignalDumpToFd() const {
    return signal_dump_to_fd_.load(std::memory_order_acquire);
  }

  const std::string *StacktraceDumpToFile() const {
    return stacktrace_dump_to_file_.load(std::memory_order_acquire);
  }

  // RunSigChain runs the SigChain of signum, see SigChain.
  // invoke raises a signal to the Go signal handler.
  template <class Invoke>
  void RunSigChain(int signum, Invoke invoke) {
    if (signum >= 0 && signum < kMaxSignum) {
      GotSignals()[signum] = 1;
    }
    const SigChain *chain = sig_invoke_chains_.Load(signum);
    if (chain == nullptr) {
      return;
    }
    int to = chain->to;
    int wait = chain->wait;
    bool waiting = wait >= 0 && wait < kMaxSignum && wait != signum;

    // Reset before raising {to}, or {wait} received in between is lost.
    if (waiting) {
      GotSignals()[wait] = 0;
    }
    if (to >= 0 && to != signum) {
      invoke(to);
    }
    if (waiting && !WaitSignal(wait, chain)) {
      int fd = SignalDumpToFd();
      if (fd >= 0) {
        (void)!write(fd, "wait Signal(", strlen("wait Signal("));
        (void)!WriteInt(fd, wait);
        (void)!write(fd, ") for Signal(", strlen(") for Signal("));
        (void)!WriteInt(fd, signum);
        (void)!write(fd, ") canceled.\n", strlen(") canceled.\n"));
      }
    }
    if (chain->sleepInSeconds > 0) {
      sleep(chain->sleepInSeconds);
    }
  }

 private:
  // GotSignals marks signals received, set by signal handlers of any thread.
  static volatile sig_atomic_t *GotSignals() {
    static volatile sig_atomic_t got_signals[kMaxSignum];
    return got_signals;
  }

  // WaitSignal waits {wait} until chain is canceled, by polling rather than
  // sigsuspend, which can not be canceled, and blocks forever sometimes, as
  // {wait} may be handled by another thread.
  static bool WaitSignal(int wait, const SigChain *chain) {
    volatile sig_atomic_t *got = GotSignals();
    for (;;) {
      if (got[wait]) {
        got[wait] = 0;
        return true;
      }
      if (chain->canceled.load(std::memory_order_acquire)) {
        return false;
      }
      SleepMillis(10);
    }
  }

  // SleepMillis is async-signal-safe.
  static void SleepMillis(int ms) {
#if defined(_WIN32)
    usleep(ms * 1000);
#else
    struct timespec ts;
    ts.tv_sec = 0;
    ts.tv_nsec = ms * 1000000L;
    nanosleep(&ts, nullptr);
#endif
  }

  std::atomic<int> signal_dump_to_fd_;
  std::atomic<const std::string *> stacktrace_dump_to_file_;
  SignalTable<SigChain> sig_invoke_chains_;
};
}  // namespace searking
#endif  // GO_OS_SIGNAL_CGO_BASE_SIGNAL_HANDLER_HPP_
