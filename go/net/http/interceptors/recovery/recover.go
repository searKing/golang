package recovery

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

func handleRecover(logger *log.Logger, f func(err interface{})) {
	if err := recover(); err != nil {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe = ErrorIsBrokenPipe(err)
		if logger != nil {
			reset := string([]byte{27, 91, 48, 109})

			goErr := errors.Wrapf(errors.New("panic"), "%v", err)

			httpRequest, _ := httputil.DumpRequest(r, false)
			headers := strings.Split(string(httpRequest), "\r\n")
			for idx, header := range headers {
				current := strings.Split(header, ":")
				if current[0] == "Authorization" {
					headers[idx] = current[0] + ": *"
				}
			}
			if brokenPipe {
				logger.Printf("[Recovery] brokenPipe %+v\n%s%s", goErr, string(httpRequest), reset)
			} else if gin.IsDebugging() {
				logger.Printf("[Recovery] %s panic recovered:\n%s\n%+v%s",
					timeFormat(time.Now()), strings.Join(headers, "\r\n"), goErr, reset)
			} else {
				logger.Printf("[Recovery] %s panic recovered:\n%+v%s",
					timeFormat(time.Now()), goErr, reset)
			}
		}
		if f != nil {
			f(err)
		}
	}
}
func timeFormat(t time.Time) string {
	var timeString = t.Format("2006/01/02 - 15:04:05")
	return timeString
}

func ErrorIsBrokenPipe(err interface{}) bool {
	var brokenPipe bool
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				brokenPipe = true
			}
		}
	}
	return brokenPipe
}
