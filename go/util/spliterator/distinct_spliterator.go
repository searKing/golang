package spliterator

import "sync"

/**
 * A wrapping spliterator that only reports distinct elements of the
 * underlying spliterator. Does not preserve size and encounter order.
 */
type DistinctSpliterator struct {
	// The underlying spliterator
	s Spliterator
	// ConcurrentHashMap holding distinct elements as keys
	sync.Map
}
