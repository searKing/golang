// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rand

import (
	"math/rand"
	"time"
)

// https://www.calhoun.io/creating-random-strings-in-go/
// we are able to isolate it so that no other code can affect our seed.
// This is important,
// because another piece of code we import might also seed the math/rand package and
// cause all of our “random” functions to not really be that random.
// For example,
// if we seed with rand.Seed(time.Now().UnixNano())
// and then another initializer calls rand.Seed(1) our seed will get overridden,
// and that definitely isn’t what we want.
// By using a rand.Rand instance
// we are able to prevent this from happening to our random number generator.
var seededRandMath = rand.New(rand.NewSource(time.Now().UnixNano()))
