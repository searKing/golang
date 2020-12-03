package runtime

import (
	"runtime"
	"strings"
)

// GetCaller returns the caller of the function that calls it.
func GetCaller() string {
	var pc [1]uintptr
	runtime.Callers(2, pc[:])
	f := runtime.FuncForPC(pc[0])
	if f == nil {
		return "Unable to find caller"
	}
	return f.Name()
}

// GetCallStack Same as stdlib http server code. Manually allocate stack trace buffer size
// to prevent excessively large logs
func GetCallStack(size int) string {
	buf := make([]byte, size)
	stk := string(buf[:runtime.Stack(buf[:], false)])
	lines := strings.Split(stk, "\n")
	if len(lines) < 3 {
		return stk
	}

	// trim GetCallStack
	var stackLines []string
	stackLines = append(stackLines, lines[0])
	stackLines = append(stackLines, lines[3:]...)

	return strings.Join(stackLines, "\n")
}
