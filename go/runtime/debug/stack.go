package debug

import (
	"runtime"
	"strings"
)

// GoroutineID returns goroutine id of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func GoroutineID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	return strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
}
