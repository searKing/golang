package runtime

import (
	"path"
	"runtime"
	"strings"
)

// GetCaller returns the caller of the function that calls it.
// The argument skip is the number of stack frames
// to skip before recording in pc, with 0 identifying the frame for Callers itself and
// 1 identifying the caller of Callers.
func GetCaller(skip int) string {
	var pc [1]uintptr
	runtime.Callers(skip+1, pc[:])
	f := runtime.FuncForPC(pc[0])
	if f == nil {
		return "Unable to find caller"
	}
	return f.Name()
}

// GetShortCaller returns the short form of GetCaller.
// The argument skip is the number of stack frames
// to skip before recording in pc, with 0 identifying the frame for Callers itself and
// 1 identifying the caller of Callers.
func GetShortCaller(skip int) string {
	return strings.TrimPrefix(path.Ext(GetCaller(skip+1)), ".")
}

// GetCallerFuncFileLine returns the __FUNCTION__, __FILE__ and __LINE__ of the function that calls it.
// The argument skip is the number of stack frames
// to skip before recording in pc, with 0 identifying the frame for Callers itself and
// 1 identifying the caller of Callers.
func GetCallerFuncFileLine(skip int) (caller string, file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(skip + 1)
	if !ok {
		file = "???"
		line = 0
	}
	return GetCaller(skip + 1), file, line
}

// GetShortCallerFuncFileLine returns the short form of GetCallerFuncFileLine.
// The argument skip is the number of stack frames
// to skip before recording in pc, with 0 identifying the frame for Callers itself and
// 1 identifying the caller of Callers.
func GetShortCallerFuncFileLine(skip int) (caller string, file string, line int) {
	caller, file, line = GetCallerFuncFileLine(skip + 1)
	return strings.TrimPrefix(path.Ext(caller), "."), path.Base(file), line
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
