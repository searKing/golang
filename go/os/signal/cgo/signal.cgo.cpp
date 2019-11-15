#include <iostream>

#include "signal.cgo.h"
#include "signal_wrap.h"

void CGOSignalAction(bool enable, int signum) {
  searking::SignalAction(enable, signum);
}