// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

//go:build cgo

#include "write_int.hpp"

#include <climits>

namespace searking {

// WriteInt is async-signal-safe, as it neither allocates nor calls stdio.
ssize_t WriteInt(int fd, int n) {
  // digits of INT_MIN, and a sign
  char nums[sizeof(int) * CHAR_BIT / 3 + 2];
  // negate in unsigned, as -INT_MIN overflows int.
  unsigned int u = n < 0 ? 0u - static_cast<unsigned int>(n)
                         : static_cast<unsigned int>(n);

  // fill digits backwards from the end
  size_t idx = sizeof(nums);
  do {
    nums[--idx] = static_cast<char>('0' + u % 10);
    u /= 10;
  } while (u != 0);
  if (n < 0) {
    nums[--idx] = '-';
  }
  return write(fd, nums + idx, sizeof(nums) - idx);
}
}  // namespace searking
