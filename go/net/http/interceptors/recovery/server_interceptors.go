package recovery

import (
	"io"
	"log"
	"net/http"
)

// ServerInterceptor returns a new server interceptors with recovery from panic.
func ServerInterceptor(next http.Handler, out io.Writer, f func(w http.ResponseWriter, r *http.Request, err interface{})) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var logger *log.Logger
		if out != nil {
			logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
		}
		defer func() {
			handleRecover(logger, func(err interface{}) {
				if f == nil {
					return
				}
				f(w, r, err)
			})
		}()
		next.ServeHTTP(w, r)
	})
}
