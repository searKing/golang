// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring

// Returns a uniquely identifying key, suitable for hashing by the
// NodeLocator algorithm.
// @param node The Node to use to form the unique identifier
// @param repetition The repetition number for the particular node in question
//          (0 is the first repetition)
// @return The key that represents the specific repetition of the node, such as “127.0.0.1:11311-0”
func (c *NodeLocator) getIterateKeyForNode(node Node, repetition int) string {
	return c.nodeKeyFormatter.getKeyForNode(node, repetition)
}

func (c *NodeLocator) getIterateHashKeyForNode(node Node, repetition int) []uint32 {
	return c.hashAlg.Hash(c.getIterateKeyForNode(node, repetition))
}

// 127.0.0.1:11311-0 -> 1122334455
// IterateKey -> IterateHashKey
func (c *NodeLocator) getHashKey(iterateKey string) uint32 {
	return c.hashAlg.Hash(iterateKey)[0]
}
