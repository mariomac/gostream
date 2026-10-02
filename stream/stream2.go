// Package stream provides type-safe streams, functional helper tools and processing
// operations
package stream

import (
	"fmt"
)

// Stream is a sequence of elements supporting different processing and aggregation functionalities.
// To perform a computation, Stream operations are composed into a stream pipeline. A stream pipeline
// consists of a source (which might be an array, a collection, a generator function, an I/O channel,
// etc), zero or more intermediate operations (which transform a stream into another stream, such as
// filter(predicate)), and a terminal operation (which produces a result or side-effect, such as
// reduce(function) or forEach(consumer)). Streams are lazy; computation on the source data is
// only performed when the terminal operation is initiated, and source elements are consumed only as
// needed.
type MapStream[K comparable, V any] struct {
	infinite bool
	supply   iterator2Supplier[K, V]
}

// if there are more items to iterate, returns the next item and true.
// if the iterator has iterated all the stream items, returns the zero value and false.
// this function will usually depend on an external status (e.g. can be a struct method or
// a function literal that rely on outer variables)
type iterator2[K comparable, V any] func() (K, V, bool)

func finishedIterator2[K comparable, V any]() (K, V, bool) {
	var zeroKey K
	var zeroVal V
	return zeroKey, zeroVal, false
}

type iterator2Supplier[K comparable, V any] func() iterator2[K, V]

func (is *MapStream[K, V]) iterator() iterator2[K, V] {
	return is.supply()
}

func (is *MapStream[K, V]) isInfinite() bool {
	return is.infinite
}

func assertFinite2[K comparable, V any](is *MapStream[K, V]) {
	if is.isInfinite() {
		var k K
		var v V
		panic(fmt.Sprintf("operation not allowed in an infinite MapStream[%T, %T]", k, v))
	}
}
