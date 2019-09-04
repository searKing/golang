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
