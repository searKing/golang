package runtime_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/searKing/golang/go/runtime"
)

func TestPanic_Recover(t *testing.T) {
	defer func() {
		if x := recover(); x == nil {
			t.Errorf("Expected a panic to recover from")
		}
	}()
	defer runtime.DefaultPanic.Recover()
	panic("Test Panic")
}

func TestPanicWith(t *testing.T) {
	var result interface{}
	func() {
		defer func() {
			if x := recover(); x == nil {
				t.Errorf("Expected a panic to recover from")
			}
		}()
		defer runtime.HandlePanicWith(func(r interface{}) {
			result = r
		}).Recover()
		panic("test")
	}()
	if result != "test" {
		t.Errorf("did not receive custom handler")
	}
}

func TestPanic_Recover_LogPanic(t *testing.T) {
	log, err := captureStderr(func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected a panic to recover from")
			}
		}()
		defer runtime.LogPanic.Recover()
		panic("test panic")
	})
	if err != nil {
		t.Fatalf("%v", err)
	}
	// Example log:
	//
	// ...] Observed a panic: test panic
	// goroutine 6 [running]:
	// github.com/searKing/golang/go/runtime.logPanic(0x12a8b80, 0x130e590)
	//	.../src/github.com/searKing/golang/go/runtime/panic.go:86 +0xda
	lines := strings.Split(log, "\n")
	if len(lines) < 4 {
		t.Fatalf("panic log should have 1 line of message, 1 line per goroutine and 2 lines per function call")
	}
	if match, _ := regexp.MatchString("Observed a panic: test panic", lines[0]); !match {
		t.Errorf("mismatch panic message: %s", lines[0])
	}
	// The following regexp's verify that Kubernetes panic log matches Golang stdlib
	// stacktrace pattern. We need to update these regexp's if stdlib changes its pattern.
	if match, _ := regexp.MatchString(`goroutine [0-9]+ \[.+\]:`, lines[1]); !match {
		t.Errorf("mismatch goroutine: %s", lines[1])
	}
	if match, _ := regexp.MatchString(`logPanic(.*)`, lines[2]); !match {
		t.Errorf("mismatch symbolized function name: %s", lines[2])
	}
	if match, _ := regexp.MatchString(`panic\.go:[0-9]+ \+0x`, lines[3]); !match {
		t.Errorf("mismatch file/line/offset information: %s", lines[3])
	}
}

func TestPanic_Recover_LogPanicSilenceHTTPErrAbortHandler(t *testing.T) {
	log, err := captureStderr(func() {
		defer func() {
			if r := recover(); r != http.ErrAbortHandler {
				t.Fatalf("expected to recover from http.ErrAbortHandler")
			}
		}()
		defer runtime.LogPanic.Recover()
		panic(http.ErrAbortHandler)
	})
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(log) > 0 {
		t.Fatalf("expected no stderr log, got: %s", log)
	}
}

// captureStderr redirects stderr to result string, and then restore stderr from backup
func captureStderr(f func()) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	bak := os.Stderr
	os.Stderr = w
	log.SetOutput(os.Stderr)
	defer func() {
		os.Stderr = bak
		log.SetOutput(os.Stderr)
	}()

	resultCh := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		resultCh <- buf.String()
	}()

	f()
	_ = w.Close()

	return <-resultCh, nil
}
