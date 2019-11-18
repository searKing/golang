/*
 *  Copyright 2019 The searKing authors. All Rights Reserved.
 *
 *  Use of this source code is governed by a MIT-style license
 *  that can be found in the LICENSE file in the root of the source
 *  tree. An additional intellectual property rights grant can be found
 *  in the file PATENTS.  All contributing project authors may
 *  be found in the AUTHORS file in the root of the source tree.
 */
#ifndef GO_OS_SIGNAL_CGO_STACKTRACE_H_
#define GO_OS_SIGNAL_CGO_STACKTRACE_H_

#include <string>
namespace searking {
namespace stacktrace {

// This function produces a stack backtrace with demangled function & method
// names.
std::string Stacktrace(int skip = 1);

/// SafeDumpToFd is low-level async-signal-safe
/// functions for dumping call stacks.
void SafeDumpToFd(int fd = 1);

}  // namespace stacktrace
}  // namespace searking
#endif  // GO_OS_SIGNAL_CGO_STACKTRACE_H_
