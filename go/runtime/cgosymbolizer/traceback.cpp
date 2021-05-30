#include "traceback.h"

#include <stdio.h>
#include <string.h>

#include <boost/stacktrace.hpp>
#include <string>
// Gather addresses from the call stack.
void cgoTraceback(cgoTracebackArg* arg) {
  try {
    // We can only unwind the current stack.
    if (arg->context != 0) {
      arg->buf[0] = 0;
      return;
    }

    std::size_t skip = 3;
    std::size_t max_depth = arg->max;

    boost::stacktrace::stacktrace stacktrace(skip, max_depth);
//    std::cout << boost::stacktrace::stacktrace();

    std::size_t i = 0;
    for (auto it = stacktrace.cbegin(); it != stacktrace.cend(); it++) {
      arg->buf[i++] = (uintptr_t)(it->address());
    }
    auto frames_count = stacktrace.size();
    // The list of addresses terminates at a 0, so make sure there is one.
    if (frames_count < 0) {
      arg->buf[0] = 0;
    } else if (frames_count < arg->max) {
      arg->buf[frames_count] = 0;
    }
  } catch (...) {
    // ignore exception
  }

}