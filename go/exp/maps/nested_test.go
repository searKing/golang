// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"

	maps_ "github.com/searKing/golang/go/exp/maps"
)

type mapOp string

const (
	opLoad             = mapOp("Load")
	opStore            = mapOp("Store")
	opLoadOrStore      = mapOp("LoadOrStore")
	opLoadAndDelete    = mapOp("LoadAndDelete")
	opDelete           = mapOp("Delete")
	opSwap             = mapOp("Swap")
	opCompareAndSwap   = mapOp("CompareAndSwap")
	opCompareAndDelete = mapOp("CompareAndDelete")
)

var mapOps = [...]mapOp{
	opLoad,
	opStore,
	opLoadOrStore,
	opLoadAndDelete,
	opDelete,
	opSwap,
	opCompareAndSwap,
	opCompareAndDelete,
}

// mapCall is a quick.Generator for calls on mapInterface.
type mapCall[K ~string] struct {
	op mapOp
	k  []K
	v  any
}

func (c mapCall[K]) apply(m mapInterface[K]) (any, bool) {
	switch c.op {
	case opLoad:
		return m.Load(c.k)
	case opStore:
		m.Store(c.k, c.v)
		return nil, false
	case opLoadOrStore:
		return m.LoadOrStore(c.k, c.v)
	case opLoadAndDelete:
		return m.LoadAndDelete(c.k)
	case opDelete:
		m.Delete(c.k)
		return nil, false
	case opSwap:
		return m.Swap(c.k, c.v)
	case opCompareAndSwap:
		if m.CompareAndSwap(c.k, c.v, rand.Int()) {
			m.Delete(c.k)
			return c.v, true
		}
		return nil, false
	case opCompareAndDelete:
		if m.CompareAndDelete(c.k, c.v) {
			if _, ok := m.Load(c.k); !ok {
				return nil, true
			}
		}
		return nil, false
	default:
		panic("invalid mapOp")
	}
}

type mapResult struct {
	value any
	ok    bool
}

func (mr *mapResult) Reset() {
	var z mapResult
	*mr = z
}

func randKey[K ~string](r *rand.Rand) []K {
	b := make([]K, r.Intn(4))
	for i := range b {
		b[i] = K('a' + byte(rand.Intn(26)))
	}
	return b
}

func randValue(r *rand.Rand) any {
	b := make([]byte, r.Intn(4))
	for i := range b {
		b[i] = 'a' + byte(rand.Intn(26))
	}
	return string(b)
}

func (mapCall[K]) Generate(r *rand.Rand, size int) reflect.Value {
	c := mapCall[K]{op: mapOps[rand.Intn(len(mapOps))], k: randKey[K](r)}
	switch c.op {
	case opStore, opLoadOrStore:
		c.v = randValue(r)
	}
	return reflect.ValueOf(c)
}

func applyCalls[K ~string](m mapInterface[K], calls [2]mapCall[K]) (results []mapResult, final map[any]any) {
	for _, c := range calls {
		v, ok := c.apply(m)
		results = append(results, mapResult{v, ok})
	}

	final = make(map[any]any)
	m.Range(func(k []K, v any) bool {
		final[fmt.Sprintf("%v", k)] = v
		return true
	})

	return results, final
}

func applyMap[K ~string](calls [2]mapCall[K]) ([]mapResult, map[any]any) {
	return applyCalls[K](make(maps_.NestedMap[K]), calls)
}

func TestIssue40999(t *testing.T) {
	var m = maps_.NestedMap[*int]{}

	// Since the miss-counting in missLocked (via Delete)
	// compares the miss count with len(m.dirty),
	// add an initial entry to bias len(m.dirty) above the miss count.
	m.Store(nil, struct{}{})

	var finalized uint32

	// Set finalizers that count for collected keys. A non-zero count
	// indicates that keys have not been leaked.
	for atomic.LoadUint32(&finalized) == 0 {
		p := new(int)
		runtime.SetFinalizer(p, func(*int) {
			atomic.AddUint32(&finalized, 1)
		})
		m.Store([]*int{p}, struct{}{})
		m.Delete([]*int{p})
		runtime.GC()
	}
}

func TestNestedMapRangeCall(t *testing.T) { // Issue 46399
	var m = maps_.NestedMap[int]{}
	for i, v := range [3]string{"hello", "world", "Go"} {
		m.Store([]int{i, i, i}, v)
	}

	var dummyKey = []int{42, 42, 42, 42}
	m.Range(func(keys []int, value any) bool {
		m.Range(func(keys []int, value any) bool {
			// We should be able to load the key offered in the Range callback,
			// because there are no concurrent Delete involved in this tested map.
			if v, ok := m.Load(keys); !ok || !reflect.DeepEqual(v, value) {
				t.Fatalf("Nested Range loads unexpected value, got %+v want %+v", v, value)
			}

			// We didn't keep dummyKey and a value into the map before, if somehow we loaded
			// a value from such a key, meaning there must be an internal bug regarding
			// nested range in the Map.
			if _, loaded := m.LoadOrStore(dummyKey, "dummy"); loaded {
				t.Fatalf("Nested Range loads unexpected value, want store a new value")
			}

			// Try to Store then LoadAndDelete the corresponding value with the key
			// 42 to the Map. In this case, the key 42 and associated value should be
			// removed from the Map. Therefore any future range won't observe key 42
			// as we checked in above.
			val := "maps_.NestedMap[int]"
			m.Store(dummyKey, val)
			if v, loaded := m.LoadAndDelete(dummyKey); !loaded || !reflect.DeepEqual(v, val) {
				t.Fatalf("Nested Range loads unexpected value, got %v, want %v", v, val)
			}
			return true
		})

		// Remove key from Map on-the-fly.
		m.Delete(keys)
		return true
	})

	// After a Range of Delete, all keys should be removed and any
	// further Range won't invoke the callback. Hence length remains 0.
	length := 0
	m.Range(func(keys []int, value any) bool {
		length++
		return true
	})

	if length != 0 {
		t.Fatalf("Unexpected maps_.NestedMap[int] size, got %v want %v", length, 0)
	}
}

func TestCompareAndSwap_NonExistingKey(t *testing.T) {
	m := &maps_.NestedMap[int]{}
	if m.CompareAndSwap([]int{404}, nil, 42) {
		t.Fatalf("CompareAndSwap on an non-existing key succeeded")
	}
}

var nestedMapTests = []struct {
	calls []mapCall[string]
	want  []mapResult
}{
	{
		calls: []mapCall[string]{
			{
				op: opStore,
				k:  nil,
				v:  nil,
			},
			{
				op: opLoad,
				k:  nil,
				v:  nil,
			},
			{
				op: opStore,
				k:  []string{"name"},
				v:  "Alice",
			},
			{
				op: opLoad,
				k:  []string{"name"},
			},
			{
				op: opLoad,
				k:  []string{"sex"},
			},
			{
				op: opLoadOrStore,
				k:  []string{"sex"},
				v:  "Male",
			},
			{
				op: opLoadAndDelete,
				k:  []string{"sex_to_delete"},
			},
			{
				op: opStore,
				k:  []string{"sex_to_delete"},
				v:  "Female",
			},
			{
				op: opLoadAndDelete,
				k:  []string{"sex_to_delete"},
			},
			{
				op: opLoadOrStore,
				k:  []string{"sex_to_delete"},
				v:  "Middle",
			},
			{
				op: opCompareAndSwap,
				k:  []string{"country", "province"},
				v:  "Nanjing",
			},
			{
				op: opLoadOrStore,
				k:  []string{"country", "province"},
				v:  "Shanghai",
			},
			{
				op: opLoad,
				k:  []string{"country", "province"},
			},
			{
				op: opCompareAndDelete,
				k:  []string{"country", "province"},
				v:  "Nanjing",
			},
			{
				op: opLoad,
				k:  []string{"country", "province"},
			},
			{
				op: opCompareAndDelete,
				k:  []string{"country", "province"},
				v:  "Shanghai",
			},
			{
				op: opLoad,
				k:  []string{"country", "province"},
			},
			{
				op: opDelete,
				k:  []string{"country", "province"},
			},
			{
				op: opLoad,
				k:  []string{"country", "province"},
			},
		},
		want: []mapResult{
			{
				value: nil,
				ok:    false,
			},
			{
				value: nil,
				ok:    false,
			},
			{
				value: nil,
				ok:    false,
			},
			{
				value: "Alice",
				ok:    true,
			},
			{
				value: nil,
				ok:    false,
			},
			{
				value: "Male",
				ok:    false,
			},
			{
				value: nil,
				ok:    false,
			},
			{
				value: nil,
				ok:    false,
			},
			{
				value: "Female",
				ok:    true,
			},
			{
				value: "Middle",
				ok:    false,
			},
			{
				value: nil,
				ok:    false,
			},
			{
				value: "Shanghai",
				ok:    false,
			},
			{
				value: "Shanghai",
				ok:    true,
			},
			{
				ok: false,
			},
			{
				value: "Shanghai",
				ok:    true,
			},
			{
				ok: true,
			},
			{
				value: nil,
				ok:    false,
			},
			{
				value: nil,
				ok:    false,
			},
			{
				value: nil,
				ok:    false,
			},
		},
	},
	{
		calls: []mapCall[string]{
			{
				op: opStore,
				k:  []string{"name"},
				v:  "Alice",
			},
			{
				op: opLoad,
				k:  []string{"name"},
			},
			{
				op: opStore,
				k:  []string{"name", "sex"},
				v:  "Female",
			},
			{
				op: opLoad,
				k:  []string{"name"},
			},
			{
				op: opLoad,
				k:  []string{"name", "sex"},
			},
			{
				op: opDelete,
				k:  []string{"name"},
			},
			{
				op: opLoad,
				k:  []string{"name"},
			},
			{
				op: opLoad,
				k:  []string{"name", "sex"},
			},
		},
		want: []mapResult{
			{},
			{
				value: "Alice",
				ok:    true,
			},
			{
				value: nil,
				ok:    false,
			},
			{
				value: maps_.NestedMap[string]{"sex": "Female"},
				ok:    true,
			},
			{
				value: "Female",
				ok:    true,
			},
			{},
			{
				value: nil,
				ok:    false,
			},
			{
				value: nil,
				ok:    false,
			},
		},
	},
}

func TestNestedMap(t *testing.T) {
	for i, test := range nestedMapTests {
		m := &maps_.NestedMap[string]{}
		var mr mapResult
		for j, c := range test.calls {
			mr.Reset()
			mr.value, mr.ok = c.apply(m)
			if !reflect.DeepEqual(mr, test.want[j]) {
				t.Errorf("#%d-%d: %s(%v, %v) = %v, want %v", i, j, c.op, c.k, c.v, mr, test.want[j])
				break
			}
		}
	}
}
