// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring

import "net"

func WithNumberNodeRepetitions(n int) KetamaNodeLocatorOption {
	return KetamaNodeLocatorOptionFunc(func(l *KetamaNodeLocator) {
		l.numReps = n
	})
}

func WithHashAlg(hashAlg HashAlgorithm) KetamaNodeLocatorOption {
	return KetamaNodeLocatorOptionFunc(func(l *KetamaNodeLocator) {
		l.hashAlg = hashAlg
	})
}

func WithFormatter(formatter *KetamaNodeKeyFormatter) KetamaNodeLocatorOption {
	return KetamaNodeLocatorOptionFunc(func(l *KetamaNodeLocator) {
		l.ketamaNodeKeyFormatter = formatter
	})
}

func WithWeights(weights map[net.Addr]int) KetamaNodeLocatorOption {
	return KetamaNodeLocatorOptionFunc(func(l *KetamaNodeLocator) {
		l.weights = weights
		l.isWeightedKetama = len(weights) > 0
	})
}
