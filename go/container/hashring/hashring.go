// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hashring provides a consistent hashing function.
//
// KetamaNodeLocator hashing is often used to distribute requests to a changing set of servers.  For example,
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
//
package hashring

import (
	"crypto/md5"
	"fmt"
	"math"
	"net"
	"sort"
)

// {}	-> 127.0.0.1:11311 -> 127.0.0.1:11311-0 -> 1234
// Node ->       Key       ->     IterateKey    -> HashKey
//			               ->     IterateKey    -> HashKey
//			               ->     IterateKey    -> HashKey
type Node interface {
	// Get the SocketAddress of the server to which this node is connected.
	GetSocketAddress() net.Addr
}

type StringNode string

// name of the network (for example, "tcp", "udp")
func (n StringNode) Network() string {
	return "string"
}

// string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
func (n StringNode) String() string {
	return string(n)
}
func (n StringNode) GetSocketAddress() net.Addr {
	return n
}

// KetamaNodeLocator holds the information about the allNodes of the consistent hash ketamaNodes.
//go:generate go-option -type "KetamaNodeLocator"
type KetamaNodeLocator struct {
	// The List of nodes to use in the Ketama consistent hash continuum
	// 模拟一致性哈希环的环状结构，存放的都是可用节点
	// 一致性Hash环
	sortedKetamaHashRing Uint32Slice       // []HashKey, Index for ketamaNodes binary search
	ketamaNodes          map[uint32]Node   // <HashKey,Node>
	allNodes             map[Node]struct{} // <Node>

	// The hash algorithm to use when choosing a node in the Ketama consistent hash continuum
	hashAlg HashAlgorithm

	// node weights for ketama, a map from InetSocketAddress to weight as Integer
	weights          map[net.Addr]int
	isWeightedKetama bool

	// the number of discrete hashes that should be defined for each node in the continuum.
	numReps int
	// the format used to name the nodes in Ketama, either SpyMemcached or LibMemcached
	ketamaNodeKeyFormatter *KetamaNodeKeyFormatter
}

// New creates a hash ring of n replicas for each entry.
func New(opts ...KetamaNodeLocatorOption) *KetamaNodeLocator {
	r := &KetamaNodeLocator{
		ketamaNodes:            make(map[uint32]Node),
		allNodes:               make(map[Node]struct{}),
		hashAlg:                KetamaHash,
		weights:                make(map[net.Addr]int),
		numReps:                160,
		ketamaNodeKeyFormatter: NewKetamaNodeKeyFormatter(SpyMemcached),
	}
	r.ApplyOptions(opts...)
	return r
}

// GetAllNodes returns all available nodes
func (c *KetamaNodeLocator) GetAllNodes() []Node {
	var nodes []Node
	for node, _ := range c.allNodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetPrimaryNode returns the first available node for a name, such as “127.0.0.1:11311-0” for "Alice"
func (c *KetamaNodeLocator) GetPrimaryNode(name string) (Node, bool) {
	return c.getNodeForHashKey(c.getHashKey(name))
}

// GetMaxHashKey returns the last available node's HashKey
// that is, Maximum HashKey in the Hash Cycle
func (c *KetamaNodeLocator) GetMaxHashKey() (uint32, error) {
	if len(c.sortedKetamaHashRing) == 0 {
		return 0, fmt.Errorf("NoSuchElementException")
	}
	return c.sortedKetamaHashRing[len(c.sortedKetamaHashRing)-1], nil
}

// getNodeForHashKey returns the first available node since iterateHashKey, such as HASH(“127.0.0.1:11311-0”)
func (c *KetamaNodeLocator) getNodeForHashKey(hash uint32) (Node, bool) {
	rv, has := c.getKetamaNodes()[hash]
	if has {
		return rv, true
	}
	firstKey, found := c.tailSearch(hash)
	if !found {
		firstKey = 0
	}

	hash = c.sortedKetamaHashRing[firstKey]
	rv, has = c.getKetamaNodes()[hash]
	return rv, has
}

// 根据输入物理节点列表，重新构造Hash环，即虚拟节点环
// updateLocator reconstructs the hash ring with the input nodes
func (c *KetamaNodeLocator) updateLocator(nodes ...Node) {
	c.allNodes = make(map[Node]struct{})
	for _, node := range nodes {
		c.allNodes[node] = struct{}{}
	}
	c.SetKetamaNodes(nodes...)
}

// GetNodeRepetitions returns the number of discrete hashes that should be defined for each node
// in the continuum.
func (c *KetamaNodeLocator) getNodeRepetitions() int {
	return c.numReps
}

// getKetamaNodes returns the ketamaNodes
func (c *KetamaNodeLocator) getKetamaNodes() map[uint32]Node {
	return c.ketamaNodes
}

// SetKetamaNodes setups the KetamaNodeLocator with the list of nodes it should use.
// If there are existing nodes not present in nodes, they will be removed.
// @param nodes a List of Nodes for this KetamaNodeLocator to use in
// its continuum
func (c *KetamaNodeLocator) SetKetamaNodes(nodes ...Node) {
	if c.isWeightedKetama {
		c.setWeightKetamaNodes(nodes...)
		return
	}
	c.setNoWeightKetamaNodes(nodes...)
}

func (c *KetamaNodeLocator) setNoWeightKetamaNodes(nodes ...Node) {
	// Set sets all the elements in the hash.
	// If there are existing elements not present in nodes, they will be removed.
	var nodesToBeRemoved []Node
	// remove missing Nodes
	for k := range c.allNodes {
		found := false
		for _, v := range nodes {
			if k.GetSocketAddress().String() == v.GetSocketAddress().String() {
				found = true
				break
			}
		}
		if !found {
			nodesToBeRemoved = append(nodesToBeRemoved, k)
		}
	}
	if len(nodesToBeRemoved) == len(nodes) {
		c.RemoveAllKetamaNodes()
	} else {
		c.removeNoWeightKetamaNodes(nodesToBeRemoved...)
	}
	// add all missing elements present in nodes.
	var nodesToBeAdded []Node
	for _, k := range nodes {
		found := false
		for v := range c.allNodes {
			if k.GetSocketAddress().String() == v.GetSocketAddress().String() {
				found = true
				break
			}
		}
		if !found {
			nodesToBeAdded = append(nodesToBeAdded, k)
		}
	}
	c.addNoWeightKetamaNodes(nodesToBeAdded...)
}

func (c *KetamaNodeLocator) setWeightKetamaNodes(nodes ...Node) {
	c.RemoveAllKetamaNodes()
	numReps := c.getNodeRepetitions()
	nodeCount := len(nodes)
	totalWeight := 0

	for _, node := range nodes {
		totalWeight += c.weights[node.GetSocketAddress()]
	}

	// add all elements present in nodes.
	for _, node := range nodes {
		thisWeight := c.weights[node.GetSocketAddress()]
		percent := float64(thisWeight) / float64(totalWeight)
		// floor(percent * numReps / 4 * nodeCount + 1e10) * 4
		pointerPerServer := (int)((math.Floor(percent*(float64(numReps))/4*float64(nodeCount) + 0.0000000001)) * 4)
		for i := 0; i < pointerPerServer/4; i++ {
			// batch size is 4 = SIZE(CRC)/SIZE(uint32)=16B/4B
			for _, position := range c.ketamaNodePositionsAtIteration(node, i) {
				c.getKetamaNodes()[position] = node
				fmt.Printf("adding node %s with weight %d in position %d\n", node, thisWeight, position)
			}
		}
	}

	for _, node := range nodes {
		c.allNodes[node] = struct{}{}
	}
	c.updateSortedKetamaNodes()
}

// RemoveAllKetamaNodes removes all nodes in the continuum....
func (c *KetamaNodeLocator) RemoveAllKetamaNodes() {
	c.sortedKetamaHashRing = nil
	c.ketamaNodes = make(map[uint32]Node)
	c.allNodes = make(map[Node]struct{})
}

// AddKetamaNodes inserts nodes into the consistent hash cycle.
func (c *KetamaNodeLocator) AddKetamaNodes(nodes ...Node) {
	if c.isWeightedKetama {
		c.addWeightKetamaNodes(nodes...)
		return
	}
	c.addNoWeightKetamaNodes(nodes...)
}

func (c *KetamaNodeLocator) addWeightKetamaNodes(nodes ...Node) {
	c.setWeightKetamaNodes(append(c.GetAllNodes(), nodes...)...)
}

func (c *KetamaNodeLocator) addNoWeightKetamaNodes(nodes ...Node) {
	numReps := c.getNodeRepetitions()

	for _, node := range nodes {
		// Ketama does some special work with md5 where it reuses chunks.
		// Check to be backwards compatible, the hash algorithm does not
		// matter for Ketama, just the placement should always be done using
		// MD5

		// KETAMA_HASH, Special Case, batch mode to speedup

		for i := 0; i < numReps; {
			positions := c.getIterateHashKeyForNode(node, i)
			if len(positions) == 0 {
				i++ // ignore no hash node
				break
			}

			kBound := len(positions)
			if kBound+i > numReps { // bound check [0, numReps)
				kBound = numReps - i
			}
			// [0, kBound)
			for k, pos := range positions {
				if k >= kBound {
					break // bound check
				}
				c.getKetamaNodes()[pos] = node
			}
			i += kBound
		}
	}

	for _, node := range nodes {
		c.allNodes[node] = struct{}{}
	}
	c.updateSortedKetamaNodes()
}

// Remove removes nodes from the consistent hash cycle...
func (c *KetamaNodeLocator) RemoveKetamaNodes(nodes ...Node) {
	if c.isWeightedKetama {
		c.removeWeightKetamaNodes(nodes...)
		return
	}
	c.removeNoWeightKetamaNodes(nodes...)
}

func (c *KetamaNodeLocator) removeWeightKetamaNodes(nodes ...Node) {
	c.setWeightKetamaNodes(c.GetAllNodes()...)
}

func (c *KetamaNodeLocator) removeNoWeightKetamaNodes(nodes ...Node) {
	numReps := c.getNodeRepetitions()

	for _, node := range nodes {
		for i := 0; i < numReps; i++ {

			for i := 0; i < numReps; {
				positions := c.getIterateHashKeyForNode(node, i)
				if len(positions) == 0 {
					i++ // ignore no hash node
					continue
				}
				i += len(positions)
				for _, pos := range positions {
					if n, has := c.ketamaNodes[pos]; has && n == node {
						delete(c.ketamaNodes, pos)
					}
				}
			}
		}
		delete(c.allNodes, node)
	}
	c.updateSortedKetamaNodes()
}

// tailSearch returns the first available node since iterateHashKey's Index, such as Index(HASH(“127.0.0.1:11311-0”))
func (c *KetamaNodeLocator) tailSearch(key uint32) (i int, found bool) {
	found = true
	f := func(x int) bool {
		return c.sortedKetamaHashRing[x] >= key
	}
	i = sort.Search(len(c.sortedKetamaHashRing), f)
	if i >= len(c.sortedKetamaHashRing) {
		found = false
	}
	return
}

// Get returns an element close to where name hashes to in the ketamaNodes.
func (c *KetamaNodeLocator) Get(name string) (Node, bool) {
	if len(c.ketamaNodes) == 0 {
		return nil, false
	}
	return c.GetPrimaryNode(name)
}

// GetTwo returns the two closest distinct elements to the name input in the ketamaNodes.
func (c *KetamaNodeLocator) GetTwo(name string) (Node, Node, bool) {
	if len(c.getKetamaNodes()) == 0 {
		return nil, nil, false
	}
	key := c.getHashKey(name)
	firstKey, found := c.tailSearch(key)
	if !found {
		firstKey = 0
	}
	firstNode, has := c.getKetamaNodes()[c.sortedKetamaHashRing[firstKey]]

	if len(c.allNodes) == 1 {
		return firstNode, nil, has
	}

	start := firstKey
	var secondNode Node
	for i := start + 1; i != start; i++ {
		if i >= len(c.sortedKetamaHashRing) {
			i = 0
		}
		secondNode = c.getKetamaNodes()[c.sortedKetamaHashRing[i]]
		if secondNode != firstNode {
			break
		}
	}
	return firstNode, secondNode, true
}

// GetN returns the N closest distinct elements to the name input in the ketamaNodes.
func (c *KetamaNodeLocator) GetN(name string, n int) ([]Node, bool) {
	if len(c.getKetamaNodes()) == 0 {
		return nil, false
	}

	if len(c.getKetamaNodes()) < n {
		n = len(c.getKetamaNodes())
	}

	key := c.getHashKey(name)
	firstKey, found := c.tailSearch(key)
	if !found {
		firstKey = 0
	}
	firstNode, has := c.getKetamaNodes()[c.sortedKetamaHashRing[firstKey]]

	nodes := make([]Node, 0, n)
	nodes = append(nodes, firstNode)

	if len(nodes) == n {
		return nodes, has
	}

	start := firstKey
	var secondNode Node
	for i := start + 1; i != start; i++ {
		if i >= len(c.sortedKetamaHashRing) {
			i = 0
		}
		secondNode = c.getKetamaNodes()[c.sortedKetamaHashRing[i]]
		if !sliceContainsMember(nodes, secondNode) {
			nodes = append(nodes, secondNode)
		}
		if len(nodes) == n {
			break
		}
	}

	return nodes, true
}

// 16B of md5
// 根据key生成16位的MD5摘要， 因此digest数组共16位
func (c *KetamaNodeLocator) computeMd5(key string) []byte {
	h := md5.New()
	h.Write([]byte(key))
	return h.Sum(nil)
}

func (c *KetamaNodeLocator) updateSortedKetamaNodes() {
	hashes := c.sortedKetamaHashRing[:0]
	// reallocate if we're holding on to too much (1/4th)
	// len(ketamaNodes) * replicas < cap / 4
	if cap(c.sortedKetamaHashRing)/(c.numReps*4) > len(c.ketamaNodes) {
		hashes = nil
	}
	for k := range c.ketamaNodes {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	c.sortedKetamaHashRing = hashes
}

func sliceContainsMember(set []Node, member Node) bool {
	for _, m := range set {
		if m == member {
			return true
		}
	}
	return false
}
