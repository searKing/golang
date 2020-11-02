// +build go1.15

package main

// Arguments to format are:
//	[1]: map type name
//	[2]: key type name
//	[3]: value type name
//	[4]: nil value of map type
const stringOneRun = `
// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *%[1]s) Load(key %[2]s) (%[3]s, bool) {
    value, ok := (*sync.Map)(m).Load(key)
    if value == nil {
    	return %[4]s, ok
    }
    return value.(%[3]s), ok
}

// Store sets the value for a key.
func (m *%[1]s) Store(key %[2]s, value %[3]s) {
    (*sync.Map)(m).Store(key, value)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *%[1]s) LoadOrStore(key %[2]s, value %[3]s) (%[3]s, bool) {
    actual, loaded := (*sync.Map)(m).LoadOrStore(key, value)
	if actual == nil {
        return %[4]s, loaded
    }
    return actual.(%[3]s), loaded
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *%[1]s) LoadAndDelete(key %[2]s) (value %[3]s, loaded bool) {
	actual, loaded := (*sync.Map)(m).LoadAndDelete(key)
	if actual == nil {
        return %[4]s, loaded
    }
    return actual.(%[3]s), loaded
}

// Delete deletes the value for a key.
func (m *%[1]s) Delete(key %[2]s) {
    (*sync.Map)(m).Delete(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently, Range may reflect any mapping for that key
// from any point during the Range call.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *%[1]s) Range(f func(key %[2]s, value %[3]s) bool) {
    (*sync.Map)(m).Range(func(key, value interface{}) bool {
        return f(key.(%[2]s), value.(%[3]s))
    })
}
`
