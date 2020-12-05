package os

import (
	"os"
	"syscall"
)

var ShutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
