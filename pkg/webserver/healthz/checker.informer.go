// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package healthz

import (
	"fmt"
	"net/http"
	"reflect"
)

type cacheSyncWaiter interface {
	WaitForCacheSync(stopCh <-chan struct{}) map[reflect.Type]bool
}

type informerSync struct {
	cacheSyncWaiter cacheSyncWaiter
}

var _ HealthChecker = &informerSync{}

// NewInformerSyncHealthz returns a new HealthChecker that will pass only if all informers in the given cacheSyncWaiter sync.
func NewInformerSyncHealthz(cacheSyncWaiter cacheSyncWaiter) HealthChecker {
	return &informerSync{
		cacheSyncWaiter: cacheSyncWaiter,
	}
}

func (i *informerSync) Name() string {
	return "informer-sync"
}

func (i *informerSync) Check(_ *http.Request) error {
	stopCh := make(chan struct{})
	// Close stopCh to force checking if informers are synced now.
	close(stopCh)

	informersByStarted := make(map[bool][]string)
	for informerType, started := range i.cacheSyncWaiter.WaitForCacheSync(stopCh) {
		informersByStarted[started] = append(informersByStarted[started], informerType.String())
	}

	if notStarted := informersByStarted[false]; len(notStarted) > 0 {
		return fmt.Errorf("%d informers not started yet: %v", len(notStarted), notStarted)
	}
	return nil
}
