// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashring

import (
	"crypto/md5"
	"encoding/binary"
	"hash/crc32"
	"hash/fnv"
)

// HashAlgorithm intents to provide hash for locating a server for a key.
type HashAlgorithm interface {
	// Hash computes the hash for the given key.
	// @return a positive integer hash
	Hash(k string) []uint32
}

type HashFunc func(k string) []uint32

func (f HashFunc) Hash(k string) []uint32 {
	return f(k)
}

// Known hashing algorithms for locating a server for a key. Note that all hash
// algorithms return 64-bits of hash, but only the lower 32-bits are
// significant. This allows a positive 32-bit number to be returned for all
// cases.
var (
	// CRCHash hash algorithm by crc32.
	CRCHash = HashFunc(crcHash)
	// CRCPerlHash as used by the perl API. This will be more consistent both
	// across multiple API users as well as java versions, but is mostly likely
	// significantly slower.
	CRCPerlHash = HashFunc(crcPerlHash)

	// FNV hashes are designed to be fast while maintaining a low collision rate.
	// The FNV speed allows one to quickly hash lots of data while maintaining a
	// reasonable collision rate.
	//
	// See http://www.isthe.com/chongo/tech/comp/fnv/
	// See "http://en.wikipedia.org/wiki/Fowler_Noll_Vo_hash"

	// FNV132Hash hash algorithm by 32-bit FNV1.
	FNV132Hash = HashFunc(fnv132Hash)
	// FNV1a32Hash hash algorithm by Variation of FNV.
	// 32-bit FNV1a.
	FNV1a32Hash = HashFunc(fnv1a32Hash)
	// FNV164Hash hash algorithm by 64-bit FNV1.
	FNV164Hash = HashFunc(fnv164Hash)
	// FNV1a64Hash hash algorithm by FNV1a.
	FNV1a64Hash = HashFunc(fnv1a64Hash)
	// FNV1128Hash hash algorithm by 128-bit FNV1.
	FNV1128Hash = HashFunc(fnv1128Hash)
	// FNV1a128Hash hash algorithm by 128-bit FNV1a.
	FNV1a128Hash = HashFunc(fnv1a128Hash)
	// KetamaHash hash algorithm by MD5-based hash algorithm used by ketama.
	KetamaHash = HashFunc(ketamaHash)
)

func crcHash(k string) []uint32 {
	rv := crc32.ChecksumIEEE([]byte(k))
	return []uint32{rv}
}

func crcPerlHash(k string) []uint32 {
	// return (crc32(shift) >> 16) & 0x7fff;

	//crc32 := crc32.NewIEEE()
	//crc32.Write([]byte(k))
	//rv := (crc32.Sum32() >> 16) & 0x7fff
	rv := crc32.ChecksumIEEE([]byte(k))
	return []uint32{(rv >> 16) & 0x7fff} // Truncate to 16-bits
}

// FNV hashes are designed to be fast while maintaining a low collision rate.
// The FNV speed allows one to quickly hash lots of data while maintaining a
// reasonable collision rate.
//
// @see <a href="http://www.isthe.com/chongo/tech/comp/fnv/">fnv
//
//	comparisons</a>
//
// @see <a href="http://en.wikipedia.org/wiki/Fowler_Noll_Vo_hash">fnv at
//
//	wikipedia</a>
func fnv164Hash(k string) []uint32 {
	hash := fnv.New64()
	_, _ = hash.Write([]byte(k))
	rv := hash.Sum64()
	return []uint32{uint32(rv & 0xffffffff)} // Truncate to 32-bits
}

// Variation of FNV.
func fnv1a64Hash(k string) []uint32 {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(k))
	rv := (hash.Sum64() >> 16) & 0x7fff
	return []uint32{uint32(rv & 0xffffffff)} // Truncate to 32-bits
}

// 32-bit FNV1.
func fnv132Hash(k string) []uint32 {
	hash := fnv.New32()
	_, _ = hash.Write([]byte(k))
	rv := (hash.Sum32() >> 16) & 0x7fff
	return []uint32{rv & 0xffffffff} // Truncate to 32-bits
}

// 32-bit FNV1a.
func fnv1a32Hash(k string) []uint32 {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(k))
	rv := (hash.Sum32() >> 16) & 0x7fff
	return []uint32{rv & 0xffffffff} // Truncate to 32-bits
}

// 128-bit FNV1.
func fnv1128Hash(k string) []uint32 {
	hash := fnv.New128()
	_, _ = hash.Write([]byte(k))
	bKey := hash.Sum(nil)
	rv := binary.LittleEndian.Uint32(bKey)
	return []uint32{rv & 0xffffffff} // Truncate to 32-bits
}

// 128-bit FNV1a.
func fnv1a128Hash(k string) []uint32 {
	hash := fnv.New128a()
	_, _ = hash.Write([]byte(k))
	bKey := hash.Sum(nil)
	rv := binary.LittleEndian.Uint32(bKey)
	return []uint32{rv & 0xffffffff} // Truncate to 32-bits
}

// MD5-based hash algorithm used by ketama.
func ketamaHash(k string) []uint32 {
	h := md5.New()
	_, _ = h.Write([]byte(k))
	digest := h.Sum(nil)
	var rvs []uint32
	// 16B -> [4B, 4B, 4B, 4B]
	for h := 0; h < 4; h++ {
		k := binary.LittleEndian.Uint32(digest[h*4:])
		rvs = append(rvs, k)
	}
	return rvs
}
