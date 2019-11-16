#include "signal_wrap.h"
#include "signal_handler.h"
#include "backtrace.h"
#include <signal.h>

namespace searking {

int setsig(int signum, SIGNAL_SA_ACTION_CALLBACK action,
           SIGNAL_SA_HANDLER_CALLBACK handler);

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
