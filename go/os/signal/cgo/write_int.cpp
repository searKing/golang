// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

// +build cgo

#include "write_int.hpp"

namespace searking {

// C++14
// constexpr int BitNumInDecimal(int n) {
//  int count = 0;
//  do {
//    n = n / 10;
//    count++;
//  } while (n > 0);
//  return count;
//}

// C++11
constexpr int BitNumInDecimal(int n) {
  return n == 0 ? 0 : 1 + BitNumInDecimal(n / 10);
}

ssize_t WriteInt(int fd, int n) {
  unsigned int unsigned_n = n >= 0 ? n : -n;
  // + 1 for sign
  const auto kBits = BitNumInDecimal(unsigned_n) + 1;
  char nums[kBits];
  for (auto i = 0; i < kBits; i++) {
    nums[i] = 0;
  }

  // push bits by a reverse order
  int idx = 0;
  do {
    nums[idx] = '0' + (unsigned_n % 10);
    idx++;
    unsigned_n /= 10;
  } while (unsigned_n && idx < sizeof(nums) / sizeof(nums[0]));
  if (n < 0 && idx < sizeof(nums) / sizeof(nums[0])) {
    nums[idx] = '-';
    idx++;
  }

  // reverse as a stack
  auto cnt = idx;
  for (auto i = 0; i < cnt / 2; i++) {
    nums[i] = nums[i] ^ nums[cnt - 1 - i];
    nums[cnt - 1 - i] = nums[i] ^ nums[cnt - 1 - i];
    nums[i] = nums[i] ^ nums[cnt - 1 - i];
  }
  return write(fd, nums, cnt);
}
}  // namespace searking
