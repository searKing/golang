// +build aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris windows

package signal

import (
	"os"
	"syscall"
)

const (
	numSig = 65 // max across all systems
)

func Signum(sig os.Signal) int {
	switch sig := sig.(type) {
	case syscall.Signal:
		i := int(sig)
		if i < 0 || i >= numSig {
			return -1
		}
		return i
	default:
		return -1
	}
}