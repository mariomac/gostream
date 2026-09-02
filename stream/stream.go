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
type Stream[T any] struct {
	infinite bool
	supply   iteratorSupplier[T]
}

// if there are more items to iterate, returns the next item and true.
// if the iterator has iterated all the stream items, returns the zero value and false.
// this function will usually depend on an external status (e.g. can be a struct method or
// a function literal that rely on outer variables)
type iterator[T any] func() (T, bool)

func finishedIterator[T any]() (T, bool) {
	var zeroVal T
	return zeroVal, false
}

type iteratorSupplier[T any] func() iterator[T]

func (is *Stream[T]) iterator() iterator[T] {
	return is.supply()
}

func (is *Stream[T]) isInfinite() bool {
	return is.infinite
}

func assertFinite[T any](is *Stream[T]) {
	if is.isInfinite() {
		var v T
		panic(fmt.Sprintf("operation not allowed in an infinite Stream[%T]", v))
	}
}
