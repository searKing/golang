#include "signal_handler.h"
#include "string.h"

namespace searking {

// https://github.com/boostorg/stacktrace/blob/5c6740b68067cbd7070d2965bfbce32e81f680c9/example/terminate_handler.cpp
void SignalHandler::operator()(int signum, siginfo_t *info, void *context) {
  if (backtrace_dump_to_) {
    // https://stackoverflow.com/questions/16891019/how-to-avoid-using-printf-in-a-signal-handler
    write(fd_, "Sig(", strlen("Sig("));
    int _signum = signum;

    char nums[10] = {0};
    int idx = 0;
    do {
      switch (_signum % 10) {
      case 0:
        nums[idx] = '0';
        break;
      case 1:
        nums[idx] = '1';
        break;
      case 2:
        nums[idx] = '2';
        break;
      case 3:
        nums[idx] = '3';
        break;
      case 4:
        nums[idx] = '4';
        break;
      case 5:
        nums[idx] = '5';
        break;
      case 6:
        nums[idx] = '6';
        break;
      case 7:
        nums[idx] = '7';
        break;
      case 8:
        nums[idx] = '8';
        break;
      case 9:
        nums[idx] = '9';
        break;
      }
      idx++;
      _signum /= 10;
    } while (_signum && idx < sizeof(nums) / sizeof(nums[0]));
    auto cnt = idx;
    for (auto i = 0; i < cnt / 2; i++) {
      nums[i] = nums[i] ^ nums[cnt - 1 - i];
      nums[cnt - 1 - i] = nums[i] ^ nums[cnt - 1 - i];
      nums[i] = nums[i] ^ nums[cnt - 1 - i];
    }
    write(fd_, nums, cnt);
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

} // namespace searking