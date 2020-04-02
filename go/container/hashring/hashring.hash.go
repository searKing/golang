// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring

// 迭代生成哈希槽候选点
func (c *KetamaNodeLocator) ketamaNodePositionsAtIteration(node Node, iteration int) []uint32 {
	return KetamaHash.Hash(c.getIterateKeyForNode(node, iteration))
}
