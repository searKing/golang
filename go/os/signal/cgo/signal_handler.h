#ifndef SEARKING_GOLANG_GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_H_
#define SEARKING_GOLANG_GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_H_
#include <signal.h>
#include <stdio.h>
#include <unistd.h>

#include <functional>
#include <map>
#include <mutex>
#include <utility>

namespace searking {

// Callbacks Predefinations
typedef void (*SIGNAL_SA_ACTION_CALLBACK)(int signum, siginfo_t *info,
                                          void *context);
typedef void (*SIGNAL_SA_HANDLER_CALLBACK)(int signum);

class SignalHandler {
public:
  SignalHandler()
      : onSignalCtx_(nullptr), onSignal_(nullptr), fd_(fileno(stdout)) {}
  void operator()(int signum, siginfo_t *info, void *context);
  void RegisterOnSignal(std::function<void(void *ctx, int fd, int signum,
                                           siginfo_t *info, void *context)>
                            callback,
                        void *ctx);

  void SetFd(int fd);

  void SetFd(FILE *fd);

private:
  std::mutex mutex_;
  int fd_;
  void *onSignalCtx_;
  std::function<void(void *ctx, int fd, int signum, siginfo_t *info,
                     void *context)>
      onSignal_;
  std::map<int,
           std::pair<SIGNAL_SA_ACTION_CALLBACK, SIGNAL_SA_HANDLER_CALLBACK>>
      sigactionHandlers_;
  friend int setsig(int signum, SIGNAL_SA_ACTION_CALLBACK action,
                    SIGNAL_SA_HANDLER_CALLBACK handler);
};
} // namespace searking
#endif // SEARKING_GOLANG_GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_H_
