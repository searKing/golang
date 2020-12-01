// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring

func WithNumberNodeRepetitions(n int) NodeLocatorOption {
	return KetamaNodeLocatorOptionFunc(func(l *NodeLocator) {
		l.numReps = n
	})
}

func WithHashAlg(hashAlg HashAlgorithm) NodeLocatorOption {
	return KetamaNodeLocatorOptionFunc(func(l *NodeLocator) {
		l.hashAlg = hashAlg
	})
}

func WithFormatter(formatter *KetamaNodeKeyFormatter) NodeLocatorOption {
	return KetamaNodeLocatorOptionFunc(func(l *NodeLocator) {
		l.nodeKeyFormatter = formatter
	})
}

func WithWeights(weights map[Node]int) NodeLocatorOption {
	return KetamaNodeLocatorOptionFunc(func(l *NodeLocator) {
		l.weightByNode = weights
		l.isWeighted = len(weights) > 0
	})
}
