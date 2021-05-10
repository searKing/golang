package runtime_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/searKing/golang/go/runtime"
)

func TestGetCaller(t *testing.T) {
	// Example:
	// github.com/searKing/golang/go/runtime_test.TestGetCaller
	caller := runtime.GetCaller(1)
	if match, _ := regexp.MatchString(`TestGetCaller(.*)`, caller); !match {
		t.Errorf("mismatch symbolized function name: %s", caller)
	}
}

func TestGetShortCaller(t *testing.T) {
	// Example:
	// TestGetCaller
	caller := runtime.GetShortCaller(1)
	if match, _ := regexp.MatchString(`TestGetShortCaller`, caller); !match {
		t.Errorf("mismatch symbolized function name: %s", caller)
	}
}

func caller() string {
	function, file, line := runtime.GetCallerFuncFileLine(1)
	return fmt.Sprintf("%s() %s:%d", function, file, line)
}

func TestGetCallerFunctionLine(t *testing.T) {
	// Example:
	// github.com/searKing/golang/go/runtime_test.caller() /Users/.../workspace/src/github.com/searKing/golang/go/runtime/extern_test.go:29
	cfl := caller()
	if match, _ := regexp.MatchString(
		`github\.com/searKing/golang/go/runtime_test\.caller\(\) .*github.com/searKing/golang/go/runtime/extern_test.go:([0-9]+)`, cfl); !match {
		t.Errorf("mismatch caller's caller file line: %s", cfl)
	}
}

func shortCaller() string {
	function, file, line := runtime.GetShortCallerFuncFileLine(1)
	return fmt.Sprintf("%s() %s:%d", function, file, line)
}

func TestGetShortCallerFuncFileLine(t *testing.T) {
	// Example:
	// shortCaller() extern_test.go:44
	cfl := shortCaller()
	if match, _ := regexp.MatchString(
		`shortCaller\(\) extern_test.go:([0-9]+)`, cfl); !match {
		t.Errorf("mismatch caller's caller file line: %s", cfl)
	}
}

func TestGetCallStack(t *testing.T) {
	stk := runtime.GetCallStack(2 << 20)

	// Example log:
	//
	// goroutine 19 [running]:
	// github.com/searKing/golang/go/runtime_test.TestGetCallStack(0xc000082900)
	//	 .../src/github.com/searKing/golang/go/runtime/extern_test.go:21 +0x3f
	// testing.tRunner(0xc000082900, 0x12ffb18)
	//	 /usr/local/go/src/testing/testing.go:1123 +0x1a3
	// created by testing.(*T).Run
	//	 /usr/local/go/src/testing/testing.go:1168 +0x648
	lines := strings.Split(stk, "\n")
	if len(lines) < 4 {
		t.Fatalf("panic log should have 1 line of message, 1 line per goroutine and 2 lines per function call")
	}

	// The following regexp's verify that Kubernetes panic log matches Golang stdlib
	// stacktrace pattern. We need to update these regexp's if stdlib changes its pattern.
	if match, _ := regexp.MatchString(`goroutine [0-9]+ \[.+\]:`, lines[0]); !match {
		t.Errorf("mismatch goroutine: %s", lines[1])
	}
	if match, _ := regexp.MatchString(`TestGetCallStack(.*)`, lines[1]); !match {
		t.Errorf("mismatch symbolized function name: %s", lines[1])
	}
	if match, _ := regexp.MatchString(`extern_test\.go:[0-9]+ \+0x`, lines[2]); !match {
		t.Errorf("mismatch file/line/offset information: %s", lines[2])
	}
	if match, _ := regexp.MatchString(`TestGetCallStack(.*)`, stk); !match {
		t.Errorf("mismatch symbolized function name: %s", stk)
	}
}
