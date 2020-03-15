// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package match

/**
 * Enum describing quantified match options -- all match, any match, none
 * match.
 */
type Kind struct {
	StopOnPredicateMatches bool
	ShortCircuitResult     bool
}

var (
	/** Do all elements match the predicate? */
	KindAny = Kind{true, true}
	/** Do any elements match the predicate? */
	KindAll = Kind{false, false}
	/** Do no elements match the predicate? */
	KindNone = Kind{true, false}
)
