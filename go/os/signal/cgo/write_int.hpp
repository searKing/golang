// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo

#ifndef GO_OS_SIGNAL_CGO_WRITE_INT_HPP_
#define GO_OS_SIGNAL_CGO_WRITE_INT_HPP_
#include <unistd.h>
namespace searking {
ssize_t WriteInt(int fd, int n);
}  // namespace searking
#endif  // GO_OS_SIGNAL_CGO_WRITE_INT_HPP_
