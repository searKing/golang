package pool_test

import (
	"context"
	"fmt"

	"github.com/searKing/golang/go/x/pool"
)

func ExampleWalk() {

	// chan WalkInfo
	walkChan := make(chan interface{}, 0)

	p := pool.Walk{}
	defer p.Wait()

	p.Walk(context.Background(), walkChan, func(name interface{}) error {
		fmt.Printf("%s\n", name)
		return nil
	})

	walkChan <- "1"
	walkChan <- "2"
	walkChan <- "3"
	walkChan <- "4"
	walkChan <- "5"
	close(walkChan)
	// Output Like:
	// 1
	// 2
	// 3
	// 4
	// 5
}
