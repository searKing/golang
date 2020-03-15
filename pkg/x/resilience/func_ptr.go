// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resilience

type funcPtr struct {
	x     interface{}
	ready func(x interface{}) error
	close func(x interface{})
}

func (r *funcPtr) Value() interface{} {
	if r == nil {
		return nil
	}
	return r.x
}

func (r *funcPtr) Ready() error {
	if r == nil {
		return nil
	}
	if r.ready == nil {
		return nil
	}
	return r.ready(r.x)
}

func (r *funcPtr) Close() {
	if r == nil {
		return
	}
	if r.close == nil {
		return
	}
	r.close(r.x)
}

func WithFunc(x interface{}, ready func(x interface{}) error,
	close func(x interface{})) (Ptr, error) {
	return &funcPtr{
		x:     x,
		ready: ready,
		close: close,
	}, nil
}

func WithFuncNewer(new func() (interface{}, error),
	ready func(x interface{}) error,
	close func(x interface{})) func() (Ptr, error) {
	return func() (Ptr, error) {
		if new == nil {
			return nil, ErrEmptyValue
		}
		x, err := new()
		if err != nil {
			return nil, err
		}
		return &funcPtr{
			x:     x,
			ready: ready,
			close: close,
		}, nil
	}
}
