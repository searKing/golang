#include "signal_wrap.h"
#include "backtrace.h"
#include <unistd.h>

#include <signal.h>
#include <string.h>

#include <map>
#include <mutex>
#include <utility>

namespace searking {

// Callbacks Predefinations
typedef void (*SIGNAL_SA_ACTION_CALLBACK)(int signum, siginfo_t *info,
                                          void *context);
typedef void (*SIGNAL_SA_HANDLER_CALLBACK)(int signum);
int setsig(int signum, SIGNAL_SA_ACTION_CALLBACK action,
           SIGNAL_SA_HANDLER_CALLBACK handler);

class SignalHandler {
public:
  SignalHandler() : onSignalCtx_(nullptr), onSignal_(nullptr) {}
  void operator()(int signum, siginfo_t *info, void *context) {
    int fd = 1;
    // https://stackoverflow.com/questions/16891019/how-to-avoid-using-printf-in-a-signal-handler
    write(fd, "Sig(", strlen("Sig("));
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
    write(fd, nums, cnt);
    write(fd, ") Backtrace:\n", strlen(") Backtrace:\n"));
    BacktraceFd(fd);
    write(fd, "Backtrace End\n", strlen("Backtrace End\n"));

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
    std::function<void(void *ctx, int signum, siginfo_t *info, void *context)>
        onSignal;
    {
      std::lock_guard<std::mutex> lock(mutex_);
      onSignalCtx = onSignalCtx_;
      onSignal = onSignal_;
    }

    if (onSignal) {
      onSignal(onSignalCtx, signum, info, context);
    }

    // SIGBUS, SIGFPE, SIGILL, or SIGSEGV
    if (signum == SIGBUS || signum == SIGFPE || signum == SIGILL ||
        signum == SIGSEGV) {
      exit(1);
    }
  }
  void RegisterOnSignal(
      std::function<void(void *ctx, int signum, siginfo_t *info, void *context)>
          callback,
      void *ctx) {
    std::lock_guard<std::mutex> lock(mutex_);
    onSignalCtx_ = ctx;
    onSignal_ = callback;
  }

private:
  std::mutex mutex_;
  void *onSignalCtx_;
  std::function<void(void *ctx, int signum, siginfo_t *info, void *context)>
      onSignal_;
  std::map<int,
           std::pair<SIGNAL_SA_ACTION_CALLBACK, SIGNAL_SA_HANDLER_CALLBACK>>
      sigactionHandlers_;
  friend int setsig(int signum, SIGNAL_SA_ACTION_CALLBACK action,
                    SIGNAL_SA_HANDLER_CALLBACK handler);
};
SignalHandler gSignalHandler;

int SignalAction(bool enable, int signum) {
  SIGNAL_SA_ACTION_CALLBACK sa_sigaction_action = nullptr;
  SIGNAL_SA_HANDLER_CALLBACK sa_sigaction_handler = nullptr;
  if (enable) {
    sa_sigaction_action = [](int signum, siginfo_t *info, void *context) {
      gSignalHandler(signum, info, context);
    };
  } else {
    sa_sigaction_handler = SIG_DFL;
  }

  return setsig(signum, sa_sigaction_action, sa_sigaction_handler);
}

int setsig(int signum, SIGNAL_SA_ACTION_CALLBACK action,
           SIGNAL_SA_HANDLER_CALLBACK handler) {
  struct sigaction sa;
  memset(&sa, 0, sizeof(sa));
  sigaction(signum, nullptr, &sa);
  //  sigemptyset(&sa.sa_mask);
  //  sigfillset(&sa.sa_mask);
  if (sa.sa_flags | SA_SIGINFO) {
    gSignalHandler.sigactionHandlers_[signum] =
        std::make_pair(sa.sa_sigaction, nullptr);
  } else {
    gSignalHandler.sigactionHandlers_[signum] =
        std::make_pair(nullptr, sa.sa_handler);
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
