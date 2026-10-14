// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo

#ifndef GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_HPP_
#define GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_HPP_

// SignalHandler of the platform, by sigaction on unix, or by signal on windows.
#if defined(_WIN32)
#include "signal_handler_windows.hpp"
#else
#include "signal_handler_unix.hpp"
#endif

#endif  // GO_OS_SIGNAL_CGO_SIGNAL_HANDLER_HPP_
