// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring

import (
	"fmt"
	"reflect"
)

var _ Formatter[any] = (FormatterFunc[any])(nil)

type FormatterFunc[Node comparable] func(node Node, repetition int) string

func (f FormatterFunc[Node]) FormatNodeKey(node Node, repetition int) string {
	return f(node, repetition)
}

// Formatter is used to format node for assigning nodes around the ring
type Formatter[Node comparable] interface {
	// FormatNodeKey returns a uniquely identifying key, suitable for hashing by the
	// HashRing algorithm.
	FormatNodeKey(node Node, repetition int) string
}

var _ Formatter[any] = (*KetamaNodeKeyFormatter[any])(nil)

// Format describes known key formats used in Ketama for assigning nodes around the ring
type Format int

const (
	// SpyMemcached uses the format traditionally used by spymemcached to map
	// nodes to names. The format is HOSTNAME/IP:PORT-ITERATION
	//
	// This default implementation uses the socket-address of the Node
	// and concatenates it with a hyphen directly against the repetition number
	// for example a key for a particular server's first repetition may look like:
	// "myhost/10.0.2.1-0", for the second repetition: "myhost/10.0.2.1-1"
	//
	// for a server where reverse lookups are failing the returned keys may look
	// like "/10.0.2.1-0" and "/10.0.2.1-1"
	SpyMemcached Format = iota

	// LibMemcached uses the format traditionally used by libmemcached to map
	// nodes to names. The format is HOSTNAME:[PORT]-ITERATION the PORT is not
	// part of the node identifier if it is the default memcached port (11211)
	LibMemcached
)

type KetamaNodeKeyFormatter[Node comparable] struct {
	format Format

	// Carried over from the DefaultKetamaHashRingConfiguration:
	// Internal lookup map to try to carry forward the optimisation that was
	// previously in HashRing
	keyByNode map[Node]string
}

func (f KetamaNodeKeyFormatter[Node]) GetFormat() Format {
	return f.format
}

func NewKetamaNodeKeyFormatter[Node comparable](format Format) *KetamaNodeKeyFormatter[Node] {
	return &KetamaNodeKeyFormatter[Node]{
		format:    format,
		keyByNode: make(map[Node]string),
	}
}

// FormatNodeKey returns a uniquely identifying key, suitable for hashing by the
// HashRing algorithm.
//
// @param node The Node to use to form the unique identifier
// @param repetition The repetition number for the particular node in question
//
//	(0 is the first repetition)
//
// @return The key that represents the specific repetition of the node
func (f KetamaNodeKeyFormatter[Node]) FormatNodeKey(node Node, repetition int) string {
	// Carried over from the DefaultKetamaHashRingConfiguration:
	// Internal Using the internal map retrieve the socket addresses
	// for given nodes.
	// I'm aware that this code is inherently thread-unsafe as
	// I'm using a HashMap implementation of the map, but the worst
	// case ( I believe) is we're slightly in-efficient when
	// a node has never been seen before concurrently on two different
	// threads, so it the socket-address will be requested multiple times!
	// all other cases should be as fast as possible.
	nodeKey, has := f.keyByNode[node]
	if !has {
		if reflect.TypeOf(node).Implements(reflect.TypeOf((*Formatter[Node])(nil)).Elem()) {
			return any(node).(Formatter[Node]).FormatNodeKey(node, repetition)
		}
		switch f.format {
		case LibMemcached:
		case SpyMemcached:
		default:
			panic(fmt.Errorf("unsupport format %d", f.format))
		}
		nodeKey = fmt.Sprintf("%v", node)
		f.keyByNode[node] = nodeKey
	}
	return fmt.Sprintf("%s-%d", nodeKey, repetition)
}
