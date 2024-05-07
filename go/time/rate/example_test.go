package rate_test

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/searKing/golang/go/time/rate"
)

func ExampleNewReorderBuffer() {
	const n = 10
	// See https://web.archive.org/web/20040724215416/http://lgjohn.okstate.edu/6253/lectures/reorder.pdf for more about Reorder buffer.
	limiter := rate.NewReorderBuffer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		i := i

		// Allocate: The dispatch stage reserves space in the reorder buffer for instructions in program order.
		r := limiter.Reserve(ctx)

		wg.Add(1)
		go func() {
			defer wg.Done()
			// Execute out of order
			runtime.Gosched() // Increase probability of a race for out-of-order mock
			//fmt.Printf("%03d Execute out of order\n", i)

			defer r.PutToken()
			//fmt.Printf("%03d Wait until in order\n", i)
			// Wait: The complete stage must wait for instructions to finish execution.
			err := r.Wait(ctx) // Commit in order
			if err != nil {
				fmt.Printf("err: %s\n", err.Error())
				return
			}
			// Complete: Finished instructions are allowed to write results in order into the architected registers.
			fmt.Printf("%03d Complete in order\n", i)
		}()
	}
	wg.Wait()
	// Output:
	// 000 Complete in order
	// 001 Complete in order
	// 002 Complete in order
	// 003 Complete in order
	// 004 Complete in order
	// 005 Complete in order
	// 006 Complete in order
	// 007 Complete in order
	// 008 Complete in order
	// 009 Complete in order

}

func ExampleNewFullBurstLimiter() {
	const (
		burst = 3
	)
	limiter := rate.NewFullBurstLimiter(burst)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// expect dropped, as limiter is initialized with full tokens(3)
	limiter.PutToken()

	for i := 0; ; i++ {
		// fmt.Printf("%03d %s\n", i, time.Now().Format(time.RFC3339))
		fmt.Printf("Wait %03d, tokens left: %d\n", i, limiter.Tokens())
		err := limiter.Wait(ctx)
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			return
		}
		fmt.Printf("Got %03d, tokens left: %d\n", i, limiter.Tokens())

		// actor mocked by gc
		runtime.GC()

		if i == 0 {
			// refill one token
			limiter.PutToken()
		}
	}
	// Output:
	// Wait 000, tokens left: 3
	// Got 000, tokens left: 2
	// Wait 001, tokens left: 3
	// Got 001, tokens left: 2
	// Wait 002, tokens left: 2
	// Got 002, tokens left: 1
	// Wait 003, tokens left: 1
	// Got 003, tokens left: 0
	// Wait 004, tokens left: 0
	// err: context deadline exceeded
}

func ExampleNewEmptyBurstLimiter() {
	const (
		burst       = 3
		concurrency = 2
	)
	limiter := rate.NewEmptyBurstLimiter(burst)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	fmt.Printf("tokens left: %d\n", limiter.Tokens())

	// expect not allowed, as limiter is initialized with empty tokens(0)
	if limiter.Allow() {
		fmt.Printf("allow passed\n")
	} else {
		fmt.Printf("allow refused\n")
	}
	// fill one token
	limiter.PutToken()
	fmt.Printf("tokens left: %d\n", limiter.Tokens())

	// expect allowed, as limiter is filled with one token(1)
	if limiter.Allow() {
		fmt.Printf("allow passed\n")
	} else {
		fmt.Printf("allow refused\n")
	}
	fmt.Printf("tokens left: %d\n", limiter.Tokens())

	var mu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// fmt.Printf("%03d %s\n", i, time.Now().Format(time.RFC3339))
			mu.Lock()
			fmt.Printf("Wait 1 Token, tokens left: %d\n", limiter.Tokens())
			mu.Unlock()
			err := limiter.Wait(ctx)
			if err != nil {
				mu.Lock()
				fmt.Printf("err: %s\n", err.Error())
				mu.Unlock()
				return
			}

			mu.Lock()
			fmt.Printf("Got 1 Token, tokens left: %d\n", limiter.Tokens())
			mu.Unlock()
		}()
	}

	time.Sleep(10 * time.Millisecond)
	for i := 0; i < concurrency; i++ {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		fmt.Printf("PutToken #%d: before tokens left: %d\n", i, limiter.Tokens())
		// fill one token
		limiter.PutToken()
		fmt.Printf("PutToken #%d: after tokens left: %d\n", i, limiter.Tokens())
		mu.Unlock()
	}
	wg.Wait()
	fmt.Printf("tokens left: %d\n", limiter.Tokens())

	// expect allowed, as limiter is filled with one token(1)
	if limiter.Allow() {
		fmt.Printf("allow passed\n")
	} else {
		fmt.Printf("allow refused\n")
	}
	fmt.Printf("tokens left: %d\n", limiter.Tokens())

	// expect not allowed, as limiter is initialized with empty tokens(0)
	if limiter.Allow() {
		fmt.Printf("allow passed\n")
	} else {
		fmt.Printf("allow refused\n")
	}
	// Output:
	// tokens left: 0
	// allow refused
	// tokens left: 1
	// allow passed
	// tokens left: 0
	// Wait 1 Token, tokens left: 0
	// Wait 1 Token, tokens left: 0
	// PutToken #0: before tokens left: 0
	// PutToken #0: after tokens left: 0
	// Got 1 Token, tokens left: 0
	// PutToken #1: before tokens left: 0
	// PutToken #1: after tokens left: 0
	// Got 1 Token, tokens left: 0
	// tokens left: 0
	// allow refused
	// tokens left: 0
	// allow refused
}

func ExampleBurstLimiter_Reserve() {
	const (
		burst = 1
		n     = 10
	)
	limiter := rate.NewFullBurstLimiter(burst)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	// expect dropped, as limiter is initialized with full tokens(1)
	limiter.PutToken()

	type Reservation struct {
		index int
		r     *rate.Reservation
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	var rs []*Reservation

	for i := 0; i < n; i++ {
		// fmt.Printf("%03d %s\n", i, time.Now().Format(time.RFC3339))
		// fmt.Printf("Reserve %03d\n", i)
		r := &Reservation{
			index: i,
			r:     limiter.Reserve(ctx),
		}
		if i%2 == rand.Intn(2)%2 {
			rs = append(rs, r)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			// fmt.Printf("%03d %s\n", r.index, time.Now().Format(time.RFC3339))
			// fmt.Printf("Wait %03d\n", r.index)
			err := r.r.Wait(ctx)
			if err != nil {
				mu.Lock()
				fmt.Printf("err: %s\n", err.Error())
				mu.Unlock()
			}

			mu.Lock()
			fmt.Printf("%03d Got 1 Token, tokens left: %d\n", r.index, limiter.Tokens())
			mu.Unlock()
			r.r.PutToken()
		}()
	}

	for i := 0; i < len(rs); i++ {
		r := rs[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			// fmt.Printf("%03d %s\n", r.index, time.Now().Format(time.RFC3339))
			// fmt.Printf("Wait %03d\n", r.index)
			err := r.r.Wait(ctx)
			if err != nil {
				mu.Lock()
				fmt.Printf("err: %s\n", err.Error())
				mu.Unlock()
			}

			mu.Lock()
			fmt.Printf("%03d Got 1 Token, tokens left: %d\n", r.index, limiter.Tokens())
			mu.Unlock()
			r.r.PutToken()
		}()
	}
	wg.Wait()
	// Output:
	// 000 Got 1 Token, tokens left: 0
	// 001 Got 1 Token, tokens left: 0
	// 002 Got 1 Token, tokens left: 0
	// 003 Got 1 Token, tokens left: 0
	// 004 Got 1 Token, tokens left: 0
	// 005 Got 1 Token, tokens left: 0
	// 006 Got 1 Token, tokens left: 0
	// 007 Got 1 Token, tokens left: 0
	// 008 Got 1 Token, tokens left: 0
	// 009 Got 1 Token, tokens left: 0

}
