// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hashring provides a consistent hashing function.
//
// NodeLocator hashing is often used to distribute requests to a changing set of servers.  For example,
// say you have some cache servers cacheA, cacheB, and cacheC.  You want to decide which cache server
// to use to look up information on a user.
//
// You could use a typical hash table and hash the user id
// to one of cacheA, cacheB, or cacheC.  But with a typical hash table, if you add or remove a server,
// almost all keys will get remapped to different results, which basically could bring your service
// to a grinding halt while the caches get rebuilt.
//
// With a consistent hash, adding or removing a server drastically reduces the number of keys that
// get remapped.
//
// Read more about consistent hashing on wikipedia:  http://en.wikipedia.org/wiki/Consistent_hashing
package hashring

import (
	"fmt"
	"math"
	"sort"
)

// {}	-> 127.0.0.1:11311 -> 127.0.0.1:11311-0 -> 1234
// Node ->       Key       ->     IterateKey    -> HashKey
//
//	->     IterateKey    -> HashKey
//	->     IterateKey    -> HashKey
type Node interface {
	// Get the SocketAddress of the server to which this node is connected.
	fmt.Stringer
}

const defaultNumReps = 160

type StringNode string

func (n StringNode) String() string {
	return string(n)
}

// NodeLocator holds the information about the allNodes of the consistent hash nodes.
//
//go:generate go-option -type "NodeLocator"
type NodeLocator struct {
	// The List of nodes to use in the Ketama consistent hash continuum
	// 模拟一致性哈希环的环状结构，存放的都是可用节点
	// 一致性Hash环
	sortedKeys uint32Slice       // []HashKey, Index for nodes binary search
	nodeByKey  map[uint32]Node   // <HashKey,Node>
	allNodes   map[Node]struct{} // <Node>

	// The hash algorithm to use when choosing a node in the Ketama consistent hash continuum
	hashAlg HashAlgorithm

	// node weights for ketama, a map from InetSocketAddress to weight as Integer
	weightByNode map[Node]int
	isWeighted   bool

	// the number of discrete hashes that should be defined for each node in the continuum.
	numReps int
	// the format used to name the nodes in Ketama, either SpyMemcached or LibMemcached
	nodeKeyFormatter *KetamaNodeKeyFormatter
}

// New creates a hash ring of n replicas for each entry.
func New(opts ...NodeLocatorOption) *NodeLocator {
	r := &NodeLocator{
		nodeByKey:        make(map[uint32]Node),
		allNodes:         make(map[Node]struct{}),
		hashAlg:          KetamaHash,
		weightByNode:     make(map[Node]int),
		numReps:          defaultNumReps,
		nodeKeyFormatter: NewKetamaNodeKeyFormatter(SpyMemcached),
	}
	r.ApplyOptions(opts...)
	return r
}

// GetAllNodes returns all available nodes
func (c *NodeLocator) GetAllNodes() []Node {
	var nodes []Node
	for node := range c.allNodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetPrimaryNode returns the first available node for a name, such as “127.0.0.1:11311-0” for "Alice"
func (c *NodeLocator) GetPrimaryNode(name string) (Node, bool) {
	return c.getNodeForHashKey(c.getHashKey(name))
}

// GetMaxHashKey returns the last available node's HashKey
// that is, Maximum HashKey in the Hash Cycle
func (c *NodeLocator) GetMaxHashKey() (uint32, error) {
	if len(c.sortedKeys) == 0 {
		return 0, fmt.Errorf("NoSuchElementException")
	}
	return c.sortedKeys[len(c.sortedKeys)-1], nil
}

// getNodeForHashKey returns the first available node since iterateHashKey, such as HASH(“127.0.0.1:11311-0”)
func (c *NodeLocator) getNodeForHashKey(hash uint32) (Node, bool) {
	if len(c.sortedKeys) == 0 {
		return nil, false
	}

	rv, has := c.getNodeByKey()[hash]
	if has {
		return rv, true
	}
	firstKey, found := c.tailSearch(hash)
	if !found {
		firstKey = 0
	}

	hash = c.sortedKeys[firstKey]
	rv, has = c.getNodeByKey()[hash]
	return rv, has
}

// 根据输入物理节点列表，重新构造Hash环，即虚拟节点环
// updateLocator reconstructs the hash ring with the input nodes
func (c *NodeLocator) updateLocator(nodes ...Node) {
	c.SetNodes(nodes...)
}

// GetNodeRepetitions returns the number of discrete hashes that should be defined for each node
// in the continuum.
func (c *NodeLocator) getNodeRepetitions() int {
	return c.numReps
}

// getNodeByKey returns the nodes
func (c *NodeLocator) getNodeByKey() map[uint32]Node {
	return c.nodeByKey
}

// SetNodes setups the NodeLocator with the list of nodes it should use.
// If there are existing nodes not present in nodes, they will be removed.
// @param nodes a List of Nodes for this NodeLocator to use in
// its continuum
func (c *NodeLocator) SetNodes(nodes ...Node) {
	if c.isWeighted {
		c.setWeightNodes(nodes...)
		return
	}
	c.setNoWeightNodes(nodes...)
}

func (c *NodeLocator) setNoWeightNodes(nodes ...Node) {
	// Set sets all the elements in the hash.
	// If there are existing elements not present in nodes, they will be removed.
	var nodesToBeRemoved []Node
	// remove missing Nodes
	for k := range c.allNodes {
		var found bool
		for _, v := range nodes {
			if k.String() == v.String() {
				// found
				found = true
				break
			}
		}
		if !found {
			nodesToBeRemoved = append(nodesToBeRemoved, k)
		}
	}
	if len(nodesToBeRemoved) == len(nodes) {
		c.RemoveAllNodes()
	} else {
		c.removeNoWeightNodes(nodesToBeRemoved...)
	}
	// add all missing elements present in nodes.
	var nodesToBeAdded []Node
	for _, k := range nodes {
		var found bool
		for v := range c.allNodes {
			if k.String() == v.String() {
				found = true
				break
			}
		}
		if !found {
			nodesToBeAdded = append(nodesToBeAdded, k)
		}
	}
	c.addNoWeightNodes(nodesToBeAdded...)
}

func (c *NodeLocator) setWeightNodes(nodes ...Node) {
	c.RemoveAllNodes()
	numReps := c.getNodeRepetitions()
	nodeCount := len(nodes)
	totalWeight := 0

	for _, node := range nodes {
		totalWeight += c.weightByNode[node]
	}

	// add all elements present in nodes.
	for _, node := range nodes {
		thisWeight := c.weightByNode[node]
		percent := float64(thisWeight) / float64(totalWeight)
		// floor(percent * numReps * nodeCount + 1e10)
		pointerPerServer := (int)(math.Floor(percent*(float64(numReps))*float64(nodeCount) + 0.0000000001))
		c.addNodeWithoutSort(node, pointerPerServer)
	}

	// sort keys
	c.updateSortedNodes()
}

// RemoveAllNodes removes all nodes in the continuum....
func (c *NodeLocator) RemoveAllNodes() {
	c.sortedKeys = nil
	c.nodeByKey = make(map[uint32]Node)
	c.allNodes = make(map[Node]struct{})
}

// AddNodes inserts nodes into the consistent hash cycle.
func (c *NodeLocator) AddNodes(nodes ...Node) {
	if c.isWeighted {
		c.addWeightNodes(nodes...)
		return
	}
	c.addNoWeightNodes(nodes...)
}

func (c *NodeLocator) addWeightNodes(nodes ...Node) {
	c.setWeightNodes(append(c.GetAllNodes(), nodes...)...)
}

func (c *NodeLocator) addNoWeightNodes(nodes ...Node) {
	numReps := c.getNodeRepetitions()

	for _, node := range nodes {
		c.addNodeWithoutSort(node, numReps)
	}

	c.updateSortedNodes()
}

func (c *NodeLocator) addNodeWithoutSort(node Node, numReps int) {
	// Ketama does some special work with md5 where it reuses chunks.
	// Check to be backwards compatible, the hash algorithm does not
	// matter for Ketama, just the placement should always be done using
	// MD5

	// KETAMA_HASH, Special Case, batch mode to speedup

	for i := 0; i < numReps; {
		positions := c.getIterateHashKeyForNode(node, i)
		if len(positions) == 0 {
			numReps++
			i++ // ignore no hash node
			break
		}

		for j, pos := range positions {
			if i+j > numReps { // out of bound
				break
			}
			if _, has := c.getNodeByKey()[pos]; has {
				// skip this node, duplicated
				numReps++
				continue
			}
			c.getNodeByKey()[pos] = node
		}
		i += len(positions)
	}

	c.allNodes[node] = struct{}{}
}

// Remove removes nodes from the consistent hash cycle...
func (c *NodeLocator) RemoveNodes(nodes ...Node) {
	if c.isWeighted {
		c.removeWeightNodes(nodes...)
		return
	}
	c.removeNoWeightNodes(nodes...)
}

func (c *NodeLocator) removeWeightNodes(nodes ...Node) {
	for _, node := range nodes {
		delete(c.allNodes, node)
	}
	c.setWeightNodes(c.GetAllNodes()...)
}

func (c *NodeLocator) removeNoWeightNodes(nodes ...Node) {
	numReps := c.getNodeRepetitions()

	for _, node := range nodes {
		for i := 0; i < numReps; {
			positions := c.getIterateHashKeyForNode(node, i)
			if len(positions) == 0 {
				// ignore no hash node
				numReps++
				i++
				continue
			}

			for j, pos := range positions {
				if i+j > numReps { // out of bound
					break
				}
				if n, has := c.nodeByKey[pos]; has {
					if n.String() != node.String() {
						numReps++ // ignore no hash node
						continue
					}
					delete(c.nodeByKey, pos)
				}
			}
			i += len(positions)
		}
		delete(c.allNodes, node)
	}
	c.updateSortedNodes()
}

// tailSearch returns the first available node since iterateHashKey's Index, such as Index(HASH(“127.0.0.1:11311-0”))
func (c *NodeLocator) tailSearch(key uint32) (i int, found bool) {
	found = true
	f := func(x int) bool {
		return c.sortedKeys[x] >= key
	}
	// Search uses binary search to find and return the smallest index since iterateHashKey's Index
	i = sort.Search(len(c.sortedKeys), f)
	if i >= len(c.sortedKeys) {
		found = false
	}
	return
}

// Get returns an element close to where name hashes to in the nodes.
func (c *NodeLocator) Get(name string) (Node, bool) {
	if len(c.nodeByKey) == 0 {
		return nil, false
	}
	return c.GetPrimaryNode(name)
}

// GetTwo returns the two closest distinct elements to the name input in the nodes.
func (c *NodeLocator) GetTwo(name string) (Node, Node, bool) {
	if len(c.getNodeByKey()) == 0 {
		return nil, nil, false
	}
	key := c.getHashKey(name)
	firstKey, found := c.tailSearch(key)
	if !found {
		firstKey = 0
	}
	firstNode, has := c.getNodeByKey()[c.sortedKeys[firstKey]]

	if len(c.allNodes) == 1 {
		return firstNode, nil, has
	}

	start := firstKey
	var secondNode Node
	for i := start + 1; i != start; i++ {
		if i >= len(c.sortedKeys) {
			i = 0
		}
		secondNode = c.getNodeByKey()[c.sortedKeys[i]]
		if secondNode.String() != firstNode.String() {
			break
		}
	}
	return firstNode, secondNode, true
}

// GetN returns the N closest distinct elements to the name input in the nodes.
func (c *NodeLocator) GetN(name string, n int) ([]Node, bool) {
	if len(c.getNodeByKey()) == 0 {
		return nil, false
	}

	if len(c.getNodeByKey()) < n {
		n = len(c.getNodeByKey())
	}

	key := c.getHashKey(name)
	firstKey, found := c.tailSearch(key)
	if !found {
		firstKey = 0
	}
	firstNode, has := c.getNodeByKey()[c.sortedKeys[firstKey]]

	nodes := make([]Node, 0, n)
	nodes = append(nodes, firstNode)

	if len(nodes) == n {
		return nodes, has
	}

	start := firstKey
	var secondNode Node
	for i := start + 1; i != start; i++ {
		if i >= len(c.sortedKeys) {
			i = 0
			// take care of i++ after this loop of for
			i--
			continue
		}
		secondNode = c.getNodeByKey()[c.sortedKeys[i]]
		if !sliceContainsMember(nodes, secondNode) {
			nodes = append(nodes, secondNode)
		}
		if len(nodes) == n {
			break
		}
	}

	return nodes, true
}

func (c *NodeLocator) updateSortedNodes() {
	hashes := c.sortedKeys[:0]
	// reallocate if we're holding on to too much (1/4th)
	// len(nodes) * replicas < cap / 4
	// len(c.nodeByKey) ≈ len(c.allNodes)*c.numReps
	if cap(c.sortedKeys)/4 > len(c.nodeByKey) {
		hashes = nil
	}
	for k := range c.nodeByKey {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	c.sortedKeys = hashes
}

func sliceContainsMember(set []Node, member Node) bool {
	for _, m := range set {
		if m.String() == member.String() {
			return true
		}
	}
	return false
}
