#include <iostream>

#include "signal.cgo.h"
#include "signal_wrap.h"

int CGOSignalAction(bool enable, int signum) {
  return searking::SignalAction(enable, signum);
}