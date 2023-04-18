// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps

// A NestedMap is a map which, when it cannot find an element in itself, defers to another map.
// Modifications, however, are passed on to the supermap.
// It is used to implement nested namespaces, such as those which store local-variable bindings.
type NestedMap[K comparable] map[K]any

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m NestedMap[K]) Load(keys []K) (value any, ok bool) {
	nm := m
	l := len(keys) - 1
	for i, k := range keys {
		m2, ok := nm[k]
		// last key
		if i == l {
			return m2, ok
		}

		if !ok {
			return nil, false
		}
		m3, ok := m2.(NestedMap[K])
		if !ok {
			return nil, false
		}
		// continue search from here
		nm = m3
	}
	return nil, false
}

// Store sets the value for a key.
func (m NestedMap[K]) Store(keys []K, value any) {
	nm := m
	l := len(keys) - 1
	for i, k := range keys {
		// last key
		if i == l {
			nm[k] = value
			return
		}

		m2, ok := nm[k]
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(NestedMap[K])
			nm[k] = m3
			nm = m3
			continue
		}

		m3, ok := m2.(NestedMap[K])
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(NestedMap[K])
			nm[k] = m3
		}
		// continue search from here
		nm = m3
	}
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m NestedMap[K]) LoadOrStore(keys []K, value any) (actual any, loaded bool) {
	nm := m
	l := len(keys) - 1
	for i, k := range keys {
		m2, ok := nm[k]
		// last key
		if i == l {
			if ok {
				return m2, true
			}
			if !ok {
				nm[k] = value
				return value, false
			}
		}
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(NestedMap[K])
			nm[k] = m3
			nm = m3
			continue
		}

		m3, ok := m2.(NestedMap[K])
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(NestedMap[K])
			nm[k] = m3
		}
		// continue search from here
		nm = m3
	}
	return nil, false
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m NestedMap[K]) LoadAndDelete(keys []K) (value any, loaded bool) {
	nm := m
	l := len(keys) - 1
	for i, k := range keys {
		m2, ok := nm[k]
		if !ok {
			return
		}
		// last key
		if i == l {
			delete(nm, k)
			return m2, true
		}
		m3, ok := m2.(NestedMap[K])
		if !ok {
			return
		}
		// continue search from here
		nm = m3
	}
	return nil, false
}

// Delete deletes the value for a key.
func (m NestedMap[K]) Delete(keys []K) {
	nm := m
	l := len(keys) - 1
	for i, k := range keys {
		m2, ok := nm[k]
		if !ok {
			return
		}
		// last key
		if i == l {
			delete(nm, k)
			return
		}
		m3, ok := m2.(NestedMap[K])
		if !ok {
			return
		}
		// continue search from here
		nm = m3
	}
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m NestedMap[K]) Swap(keys []K, value any) (previous any, loaded bool) {
	nm := m
	l := len(keys) - 1
	for i, k := range keys {
		m2, ok := nm[k]
		// last key
		if i == l {
			nm[k] = value
			return m2, ok
		}
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(NestedMap[K])
			nm[k] = m3
			nm = m3
			continue
		}
		m3, ok := m2.(NestedMap[K])
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(NestedMap[K])
			nm[k] = m3
		}
		// continue search from here
		nm = m3
	}
	return nil, false
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// The old value must be of a comparable type.
func (m NestedMap[K]) CompareAndSwap(keys []K, old, new any) bool {
	nm := m
	l := len(keys) - 1
	for i, k := range keys {
		m2, ok := nm[k]
		// last key
		if i == l {
			if ok && m2 == old {
				nm[k] = new
				return true
			}
			return false
		}
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(NestedMap[K])
			nm[k] = m3
			nm = m3
			continue
		}
		m3, ok := m2.(NestedMap[K])
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(NestedMap[K])
			nm[k] = m3
		}
		// continue search from here
		nm = m3
	}
	return false
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete
// returns false (even if the old value is the nil interface value).
func (m NestedMap[K]) CompareAndDelete(keys []K, old any) (deleted bool) {
	nm := m
	l := len(keys) - 1
	for i, k := range keys {
		m2, ok := nm[k]
		if !ok {
			return false
		}
		// last key
		if i == l {
			if m2 == old {
				delete(nm, k)
				return true
			}
			return false
		}
		m3, ok := m2.(NestedMap[K])
		if !ok {
			return false
		}
		// continue search from here
		nm = m3
	}
	return false
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m NestedMap[K]) Range(f func(keys []K, value any) bool) {
	m.ranges(nil, f)
}

func (m NestedMap[K]) ranges(ks []K, f func(keys []K, value any) bool) {
	for k, v := range m {
		ks2 := append(ks, k)
		m2, ok := v.(NestedMap[K])
		if !ok {
			f(ks2, v)
			continue
		}
		m2.ranges(ks2, f)
	}
}
