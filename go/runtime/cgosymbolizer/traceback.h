// Copyright (c) 2019 The searKing authors. All Rights Reserved.
//
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file in the root of the source
// tree. An additional intellectual property rights grant can be found
// in the file PATENTS.  All contributing project authors may
// be found in the AUTHORS file in the root of the source tree.

// +build cgo

#ifndef GO_RUNTIME_CGOSYMBOLIZER_TRACEBACK_H_
#define GO_RUNTIME_CGOSYMBOLIZER_TRACEBACK_H_
// We want to get a definition for uintptr_t
#include <cstdint>


struct cgoTracebackArg {
	uintptr_t  context;
	uintptr_t  sigContext;
	uintptr_t* buf;
	uintptr_t  max;
};


#ifdef __cplusplus
extern "C" {
#endif

void cgoTraceback(cgoTracebackArg* parg);

#ifdef __cplusplus
}
#endif

#endif  // GO_RUNTIME_CGOSYMBOLIZER_TRACEBACK_H_
