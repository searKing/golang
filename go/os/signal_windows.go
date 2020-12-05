package os

import (
	"os"
)

var ShutdownSignals = []os.Signal{os.Interrupt}
