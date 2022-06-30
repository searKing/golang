// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build BOOST_STACKTRACE_USE_WINDBG_CACHED
// +build BOOST_STACKTRACE_USE_WINDBG_CACHED

package cgosymbolizer

/*
#cgo windows CXXFLAGS: -DBOOST_STACKTRACE_USE_WINDBG_CACHED
#cgo windows LDFLAGS: -lole32 -ldbgeng
#cgo !windows CXXFLAGS:
*/
import "C"

// https://www.boost.org/doc/libs/develop/doc/html/stacktrace/configuration_and_build.html
