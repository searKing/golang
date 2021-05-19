package rate_test

import (
	"context"
	"fmt"
	"time"

	"github.com/searKing/golang/go/time/rate"
)

func ExampleNewFullBurstLimiter() {
	const (
		burst = 3
	)
	limiter := rate.NewFullBurstLimiter(burst)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// expect dropped, as limiter is inited with full tokens(3)
	limiter.PutToken()

	for i := 0; ; i++ {
		//fmt.Printf("%03d %s\n", i, time.Now().Format(time.RFC3339))
		fmt.Printf("Wait %03d\n", i)
		err := limiter.Wait(ctx)
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}

		fmt.Printf("Got %03d\n", i)
		if i == 0 {
			// refill one token
			limiter.PutToken()
		}
	}
	// Output:
	// Wait 000
	// Got 000
	// Wait 001
	// Got 001
	// Wait 002
	// Got 002
	// Wait 003
	// Got 003
	// Wait 004
	// err: context deadline exceeded

}
