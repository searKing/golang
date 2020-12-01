// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring

import (
	"fmt"
	"strings"
)

// Known key formats used in Ketama for assigning nodes around the ring
type Format int

const (
	// SpyMemcached uses the format traditionally used by spymemcached to map
	// nodes to names. The format is HOSTNAME/IP:PORT-ITERATION
	//
	// <p>
	// This default implementation uses the socket-address of the Node
	// and concatenates it with a hyphen directly against the repetition number
	// for example a key for a particular server's first repetition may look like:
	// <p>
	//
	// <p>
	// <code>myhost/10.0.2.1-0</code>
	// </p>
	//
	// <p>
	// for the second repetition
	// </p>
	//
	// <p>
	// <code>myhost/10.0.2.1-1</code>
	// </p>
	//
	// <p>
	// for a server where reverse lookups are failing the returned keys may look
	// like
	// </p>
	//
	// <p>
	// <code>/10.0.2.1-0</code> and <code>/10.0.2.1-1</code>
	// </p>
	SpyMemcached Format = iota

	// LibMemcached uses the format traditionally used by libmemcached to map
	// nodes to names. The format is HOSTNAME:[PORT]-ITERATION the PORT is not
	// part of the node identifier if it is the default memcached port (11211)
	LibMemcached
)

type KetamaNodeKeyFormatter struct {
	format Format

	// Carried over from the DefaultKetamaNodeLocatorConfiguration:
	// Internal lookup map to try to carry forward the optimisation that was
	// previously in NodeLocator
	keyByNode map[Node]string
}

func (f KetamaNodeKeyFormatter) GetFormat() Format {
	return f.format
}

func NewKetamaNodeKeyFormatter(format Format) *KetamaNodeKeyFormatter {
	return &KetamaNodeKeyFormatter{
		format:    format,
		keyByNode: make(map[Node]string),
	}
}

// Returns a uniquely identifying key, suitable for hashing by the
// NodeLocator algorithm.
//
// @param node The Node to use to form the unique identifier
// @param repetition The repetition number for the particular node in question
//          (0 is the first repetition)
// @return The key that represents the specific repetition of the node
func (f KetamaNodeKeyFormatter) getKeyForNode(node Node, repetition int) string {
	// Carrried over from the DefaultKetamaNodeLocatorConfiguration:
	// Internal Using the internal map retrieve the socket addresses
	// for given nodes.
	// I'm aware that this code is inherently thread-unsafe as
	// I'm using a HashMap implementation of the map, but the worst
	// case ( I believe) is we're slightly in-efficient when
	// a node has never been seen before concurrently on two different
	// threads, so it the socketaddress will be requested multiple times!
	// all other cases should be as fast as possible.
	nodeKey, has := f.keyByNode[node]
	if !has {
		switch f.format {
		case LibMemcached:
			nodeKey = node.String()
			break
		case SpyMemcached:
			nodeKey = node.String()
			if strings.Index(nodeKey, "/") == 0 {
				nodeKey = strings.TrimLeft(nodeKey, "/")
			}
			break
		default:
			panic(fmt.Errorf("unsupport format %d", f.format))
		}
		f.keyByNode[node] = nodeKey
	}
	return fmt.Sprintf("%s-%d", nodeKey, repetition)
}
