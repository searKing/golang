// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps_test

// This file contains reference map implementations for unit-tests.

// mapInterface is the interface Map implements.
type mapInterface[K comparable] interface {
	Load(keys []K) (any, bool)
	Store(keys []K, value any)
	LoadOrStore(keys []K, value any) (actual any, loaded bool)
	LoadAndDelete(keys []K) (value any, loaded bool)
	Delete(keys []K)
	Swap(keys []K, value any) (previous any, loaded bool)
	CompareAndSwap(keys []K, old, new any) (swapped bool)
	CompareAndDelete(keys []K, old any) (deleted bool)
	Range(func(keys []K, value any) (shouldContinue bool))
}
