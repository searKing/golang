package runtime_test

import (
	"fmt"

	"github.com/searKing/golang/go/runtime"
)

func ExampleGetCaller() {
	caller := runtime.GetCaller()
	fmt.Print(caller)

	// Output:
	// github.com/searKing/golang/go/runtime_test.ExampleGetCaller
}
