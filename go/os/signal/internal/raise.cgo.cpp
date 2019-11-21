/*
 * Copyright (c) 2019 The searKing authors. All Rights Reserved.
 *
 * Use of this source code is governed by a MIT-style license
 * that can be found in the LICENSE file in the root of the source
 * tree. An additional intellectual property rights grant can be found
 * in the file PATENTS.  All contributing project authors may
 * be found in the AUTHORS file in the root of the source tree.
 */
#include "raise.cgo.h"

void MustSegmentFault() {
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wc++11-compat-deprecated-writable-strings"
  char* c = "hello world";
  c[1] = 'H';
#pragma GCC diagnostic pop
}
