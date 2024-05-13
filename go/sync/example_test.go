// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync_test

import (
	"context"
	"fmt"

	"github.com/searKing/golang/go/sync"
)

func ExampleWalk() {

	// chan WalkInfo
	walkChan := make(chan any, 0)

	p := sync.Walk{
		Burst: 1,
	}
	defer p.Wait()

	p.Walk(context.Background(), walkChan, func(name any) error {
		fmt.Printf("%s\n", name)
		return nil
	})

	for i := 0; i < 5; i++ {
		walkChan <- fmt.Sprintf("%d", i)
	}
	close(walkChan)
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}
