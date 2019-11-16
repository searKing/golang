package cgo

/*
#cgo CXXFLAGS: -std=c++11
#include "signal.cgo.h"
*/
import "C"

//go:generate make rebuild
