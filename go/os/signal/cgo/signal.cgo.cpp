#include <iostream>

#include "signal.cgo.h"
#include "signal_handler.h"

int CGOSignalAction(int signum) {
  return searking::SignalHandler::SignalAction(signum);
}