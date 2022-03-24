package logrus

import (
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// ShortCallerPrettyfier modify the content of the function and
// file keys in the data when ReportCaller is activated.
// INFO[0000] main.go:23 main() hello world
var ShortCallerPrettyfier = func(f *runtime.Frame) (function string, file string) {
	funcname := path.Base(f.Function)
	filename := path.Base(f.File)
	return fmt.Sprintf("%s()", funcname), fmt.Sprintf("%s:%d", filename, f.Line)
}

// WithReporterCaller enhances logrus log to log caller info
// callerPrettyfier is defined in logrus.TextFormatter and logrus.JSONFormatter.
// callerPrettyfier affects only if formatter is logrus.TextFormatter and logrus.JSONFormatter.
// if callerPrettyfier not set, ShortCallerPrettyfier is set
func WithReporterCaller(log *logrus.Logger, callerPrettyfier func(*runtime.Frame) (function string, file string)) {
	if log == nil {
		return
	}

	log.SetReportCaller(true)
	if callerPrettyfier == nil {
		callerPrettyfier = ShortCallerPrettyfier
	}

	switch f := log.Formatter.(type) {
	case *logrus.TextFormatter:
		if f.CallerPrettyfier == nil {
			f.CallerPrettyfier = callerPrettyfier
		}
	case *logrus.JSONFormatter:
		if f.CallerPrettyfier == nil {
			f.CallerPrettyfier = callerPrettyfier
		}
	}
	return
}
