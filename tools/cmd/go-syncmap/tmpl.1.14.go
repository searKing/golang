// +build !go1.15

package main

// Arguments to format are:
//	[1]: map type name
//	[2]: key type name
//	[3]: value type name
//	[4]: nil value of map type
const stringOneRun = `func (m *%[1]s) Store(key %[2]s, value %[3]s) {
    (*sync.Map)(m).Store(key, value)
}

func (m *%[1]s) LoadOrStore(key %[2]s, value %[3]s) (%[3]s, bool) {
    actual, loaded := (*sync.Map)(m).LoadOrStore(key, value)
	if actual == nil {
        return %[4]s, loaded
    }
    return actual.(%[3]s), loaded
}

func (m *%[1]s) Load(key %[2]s) (%[3]s, bool) {
    value, ok := (*sync.Map)(m).Load(key)
    if value == nil {
    	return %[4]s, ok
    }
    return value.(%[3]s), ok
}

func (m *%[1]s) Delete(key %[2]s) {
    (*sync.Map)(m).Delete(key)
}

func (m *%[1]s) Range(f func(key %[2]s, value %[3]s) bool) {
    (*sync.Map)(m).Range(func(key, value interface{}) bool {
        return f(key.(%[2]s), value.(%[3]s))
    })
}
`
