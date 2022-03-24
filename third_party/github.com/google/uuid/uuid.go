// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"encoding/binary"

	"github.com/google/uuid"
)

func IDUint8() uint8 {
	var u = uuid.New()
	return uint8(u[0])
}

func IDUint16() uint16 {
	var u = uuid.New()
	return binary.BigEndian.Uint16(u[0:2])
}

func IDUint32() uint32 {
	return uuid.New().ID()
}

func IDUint64() uint64 {
	var u = uuid.New()
	return binary.BigEndian.Uint64(u[0:8])
}
