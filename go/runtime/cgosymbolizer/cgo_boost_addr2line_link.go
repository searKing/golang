// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build BOOST_STACKTRACE_USE_ADDR2LINE_LINK
// +build BOOST_STACKTRACE_USE_ADDR2LINE_LINK

package cgosymbolizer

/*
#cgo windows CXXFLAGS:
#cgo !windows CXXFLAGS: -DBOOST_STACKTRACE_USE_ADDR2LINE
#cgo !windows LDFLAGS: -lboost_stacktrace_addr2line -ldl -rdynamic
*/
import "C"

// https://www.boost.org/doc/libs/develop/doc/html/stacktrace/configuration_and_build.html
