package rate_test

import (
	"context"
	"fmt"
	"time"

	"github.com/searKing/golang/go/time/rate"
)

func ExampleNewBurstLimiter() {
	limiter := rate.NewFullBurstLimiter(3)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// expect droped, as limiter is inited with full tokens(3)
	limiter.PutToken()

	for i := 0; ; i++ {
		//fmt.Printf("%03d %s\n", i, time.Now().Format(time.RFC3339))
		fmt.Printf("%03d\n", i)
		err := limiter.Wait(ctx)
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}

		if i == 0 {
			// refill one token
			limiter.PutToken()
		}
	}
	// Output:
	// 000
	// 001
	// 002
	// 003
	// 004
	// 005
	// err: context deadline exceeded

}
