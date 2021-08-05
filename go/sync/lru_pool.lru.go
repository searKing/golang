// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sync

import (
	"container/list"
)

type resourceLRU struct {
	ll *list.List // list.Element.Value type is of *PersistResource
	m  map[*PersistResource]*list.Element
}

// add adds pc to the head of the linked list.
func (cl *resourceLRU) add(pc *PersistResource) {
	if cl.ll == nil {
		cl.ll = list.New()
		cl.m = make(map[*PersistResource]*list.Element)
	}
	ele := cl.ll.PushFront(pc)
	if _, ok := cl.m[pc]; ok {
		panic("PersistResource was already in LRU")
	}
	cl.m[pc] = ele
}

func (cl *resourceLRU) removeOldest() *PersistResource {
	ele := cl.ll.Back()
	pc := ele.Value.(*PersistResource)
	cl.ll.Remove(ele)
	delete(cl.m, pc)
	return pc
}

// remove remove pc from cl.
func (cl *resourceLRU) remove(pc *PersistResource) {
	if ele, ok := cl.m[pc]; ok {
		cl.ll.Remove(ele)
		delete(cl.m, pc)
	}
}

// len returns the number of items in the cache.
func (cl *resourceLRU) len() int {
	return len(cl.m)
}
