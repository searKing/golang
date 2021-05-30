// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

// +build cgo

#ifndef GO_RUNTIME_CGOSYMBOLIZER_CGO_H_
#define GO_RUNTIME_CGOSYMBOLIZER_CGO_H_
// We want to get a definition for uintptr_t
#include <cstdint>

struct cgoSymbolizerMore {
  struct cgoSymbolizerMore* more;

  const char* file;
  uintptr_t lineno;
  const char* func;
};

// runtime/traceback.go
struct cgoSymbolizerArg {
  uintptr_t pc;      // program counter to fetch information for
  const char* file;  // file name (NUL terminated)
  uintptr_t lineno;  // line number
  const char* func;  // function name (NUL terminated)
  uintptr_t entry;   // function entry point
  uintptr_t more;    // set non-zero if more info for this PC
  //	uintptr_t   data;// unused by runtime, available for function
  cgoSymbolizerMore* data;  // unused by runtime, available for function
};

#ifdef __cplusplus
extern "C" {
#endif

void cgoSymbolizer(cgoSymbolizerArg* arg);

#ifdef __cplusplus
}
#endif

#endif  // GO_RUNTIME_CGOSYMBOLIZER_CGO_H_
