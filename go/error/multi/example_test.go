package multi_test

import (
	"fmt"
	"github.com/searKing/golang/go/error/multi"
)

func ExampleNew() {
	err := multi.New(fmt.Errorf("whoops"), fmt.Errorf("foo"))
	fmt.Println(err)

	// Output: whoops|foo
}

func ExampleFormat() {
	err := multi.New(fmt.Errorf("whoops"), fmt.Errorf("foo"))
	fmt.Printf("%+v", err)

	// Output:
	// Multiple errors occurred:
	//	whoops|foo
}
