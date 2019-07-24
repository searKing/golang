package tag

import "reflect"

type tagOpts struct {
	TagHandler func(val reflect.Value, tag reflect.StructTag) error
}

// Convert wrapper of convertState
func Tag(v interface{}, tagHandler func(val reflect.Value, tag reflect.StructTag) error) error {
	e := newTagState()
	err := e.handle(v, tagOpts{tagHandler})
	if err != nil {
		return err
	}

	e.Reset()
	tagStatePool.Put(e)
	return nil
}

// Tagger is the interface implemented by types that
// can marshal themselves into valid JSON.
type Tagger interface {
	TagDefault() error
}
