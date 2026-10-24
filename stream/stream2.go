package stream

// Stream2 is a sequence of key-value-like elements supporting different processing and aggregation functionalities.
// Stream2 main goal is the stream processing of map-like structures, hereby its first
// type argument must fulfill the "comparable" constraint.
//
// For sequential, unconstrained data structures, see Stream.
type Stream2[K comparable, V any] struct {
	supply iterator2Supplier[K, V]
}

type iterator2[K comparable, V any] func() (K, V, bool)

func finishedIterator2[K comparable, V any]() (K, V, bool) {
	return *new(K), *new(V), false
}

type iterator2Supplier[K comparable, V any] func() iterator2[K, V]

func (is *Stream2[K, V]) iterator() iterator2[K, V] {
	return is.supply()
}
