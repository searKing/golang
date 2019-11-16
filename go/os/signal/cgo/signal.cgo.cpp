#include <iostream>

#include "backtrace.h"
#include "signal.cgo.h"
#include "signal_handler.h"
int CGOSignalHandlerSignalAction(int signum) {
  return searking::SignalHandler::SignalAction(signum);
}
void CGOSignalHandlerSetFd(int fd) {
  searking::SignalHandler::GetInstance().SetFd(fd);
}

void CGOSignalHandlerSetBacktraceDump(bool enable) {
  if (enable) {
    searking::SignalHandler::GetInstance().SetBacktraceDumpTo(
        searking::BacktraceFd);
    return;
  }
  searking::SignalHandler::GetInstance().SetBacktraceDumpTo(nullptr);
}