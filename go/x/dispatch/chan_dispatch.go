// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dispatch

import "context"

type ChanDispatch struct {
	readerChan chan<- interface{}
	*Dispatch
}

func NewChanDispatch(handler Handler, concurrentReadMax int) *ChanDispatch {
	return NewChanDispatch3(handler, concurrentReadMax, -1)
}
func NewChanDispatch3(handler Handler, concurrentReadMax int, concurrentHandleMax int) *ChanDispatch {
	if concurrentReadMax < 0 {
		return nil
	}
	readerChan := make(chan interface{}, concurrentReadMax)
	return &ChanDispatch{
		readerChan: readerChan,
		Dispatch: NewDispatch3(ReaderFunc(func(ctx context.Context) (interface{}, error) {
			select {
			case msg := <-readerChan:
				return msg, nil
			case <-ctx.Done():
				return nil, nil
			}
		}), handler, concurrentHandleMax),
	}
}
func (thiz *ChanDispatch) SendMessage(message interface{}) bool {
	select {
	case thiz.readerChan <- message:
		return true
	default:
		return false
	}
}
