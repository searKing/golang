#include "signal_handler.h"
#include "backtrace.h"

namespace searking {
void SignalHandler::operator()(int signum, siginfo_t *info, void *context) {
  // https://github.com/boostorg/stacktrace/blob/5c6740b68067cbd7070d2965bfbce32e81f680c9/example/terminate_handler.cpp
  ::signal(signum, SIG_DFL);
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
  BacktraceFd(fd_);
  write(fd_, "Backtrace End\n", strlen("Backtrace End\n"));

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

  void *onSignalCtx = nullptr;
  std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                     void *context)>
      onSignal;
  {
    std::lock_guard<std::mutex> lock(mutex_);
    onSignalCtx = onSignalCtx_;
    onSignal = onSignal_;
  }

  if (onSignal) {
    onSignal(onSignalCtx, fd_, signum, info, context);
  }

  ::raise(signum);
  //    // SIGBUS, SIGFPE, SIGILL, or SIGSEGV
  //    if (signum == SIGBUS || signum == SIGFPE || signum == SIGILL ||
  //        signum == SIGSEGV) {
  //      exit(1);
  //    }
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

void SignalHandler::SetFd(int fd) {
  std::lock_guard<std::mutex> lock(mutex_);
  fd_ = fd;
}

void SignalHandler::SetFd(FILE *fd) { SetFd(fileno(fd)); }
} // namespace searking
