package cgo

/*
#cgo CXXFLAGS: -std=c++11
#include "signal.cgo.h"
*/
import "C"

//go:generate make rebuild
//go:generate bash ../../../../tools/scripts/cgo_include_gen.sh -p "github.com/searKing/golang/go/os/signal/cgo/include" "./include"
