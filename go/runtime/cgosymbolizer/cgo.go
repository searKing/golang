// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cgosymbolizer contains runtime support for code generated
// by the cgo tool.  See the documentation for the cgo command
// for details on using cgo.
package cgosymbolizer

/*

#cgo darwin,!arm,!arm64 LDFLAGS: -lpthread
#cgo darwin,arm LDFLAGS: -framework CoreFoundation
#cgo darwin,arm64 LDFLAGS: -framework CoreFoundation
#cgo dragonfly LDFLAGS: -lpthread
#cgo freebsd LDFLAGS: -lpthread
#cgo android LDFLAGS: -llog
#cgo !android,linux LDFLAGS: -lpthread
#cgo netbsd LDFLAGS: -lpthread
#cgo openbsd LDFLAGS: -lpthread
#cgo aix LDFLAGS: -Wl,-berok
#cgo solaris LDFLAGS: -lxnet

//#cgo CFLAGS: -Wall -Werror

#cgo solaris CPPFLAGS: -D_POSIX_PTHREAD_SEMANTICS
#cgo CXXFLAGS: -I${SRCDIR}/include/
#cgo windows CXXFLAGS: -g
#cgo !windows CXXFLAGS: -g -D_GNU_SOURCE
#cgo linux LDFLAGS: -ldl

#include <stdio.h>
#include <stdlib.h>  // Needed for C.free
*/
import "C"
