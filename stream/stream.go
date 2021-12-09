// Package stream provides type-safe streams, functional helper tools and processing
// operations
package stream

import (
	"fmt"
	"github.com/mariomac/gostream/order"
)

// Stream is a sequence of elements supporting different processing and aggregation functionalities.
// To perform a computation, Stream operations are composed into a stream pipeline. A stream pipeline
// consists of a source (which might be an array, a collection, a generator function, an I/O channel,
// etc), zero or more intermediate operations (which transform a stream into another stream, such as
// filter(predicate)), and a terminal operation (which produces a result or side-effect, such as
// reduce(function) or forEach(consumer)). Streams are lazy; computation on the source data is
// only performed when the terminal operation is initiated, and source elements are consumed only as
// needed.
type Stream[T any] interface {
	// returns a new iterator to the stream
	iterator() iterator[T]
	// returns whether the stream is infinite or not
	isInfinite() bool

	// transformation operations

	// Filter returns a Stream consisting of the items of this stream that match the given
	// predicate (this is, applying the predicate function over the item returns true).
	Filter(predicate func(T) bool) Stream[T]

	// FlatMap returns a stream consisting of the results of replacing each element of this stream
	// with the contents of a mapped stream produced by applying the provided mapping function to
	// each element. Each mapped stream is closed after its contents have been placed into this
	// stream. (If a mapped stream is null an empty stream is used, instead.)
	//
	// Due to the lazy nature of streams, if any of the mapped streams is infinite it will remain
	// unnoticed and some operations (Count, Reduce, Sorted, AllMatch...) will not end.
	//
	// If you need that the input and output Stream contain elements from different types, you need
	// to invoke the standalone function FlatMap[IN,OUT](Stream[IN], func(IN)OUT) Stream[OUT].
	FlatMap(input Stream[T], mapper func(T) Stream[T]) Stream[T]

	// Limit returns a stream consisting of the elements of this stream, truncated to
	// be no longer than maxSize in length.
	Limit(maxSize int) Stream[T]

	// Map returns a Stream consisting of the results of individually applying
	// the mapper function to each element of this Stream. The argument and return
	// value of the mapper function must belong to the same type. If you need that
	// the input and output Stream contain elements from different types, you need
	// to invoke the standalone function Map[IN,OUT](Stream[IN], func(IN)OUT) Stream[OUT].
	Map(mapper func(T) T) Stream[T]

	// Sorted returns a stream consisting of the elements of this stream, sorted according
	// to the provided order.Comparator.
	Sorted(comparator order.Comparator[T]) Stream[T]

	// terminal operations

	// AllMatch returns whether all elements of this stream match the provided predicate.
	// If this operation finds an item where the predicate is false, it stops processing
	// the rest of the stream.
	AllMatch(predicate func(T) bool) bool

	// AnyMatch returns whether any elements of this stream match the provided predicate.
	// If this operation finds an item where the predicate is true, it stops processing
	// the rest of the stream.
	AnyMatch(predicate func(T) bool) bool

	// ForEach invokes the consumer function for each item of the Stream.
	ForEach(consumer func(T))

	// NoneMatch returns whether no elements of this stream match the provided predicate.
	// If this operation finds an item where the predicate is true, it stops processing
	// the rest of the stream.
	NoneMatch(predicate func(T) bool) bool

	// Reduce performs a reduction on the elements of this stream, using an associative
	// accumulation function, and returns an value describing the reduced value, if any.
	// If no reduced value (e.g. because the stream is empty), the second returned value
	// is false.
	Reduce(accumulator func(a, b T) T) (T, bool)

	// ToSlice returns a Slice Containing all the elements of this Stream.
	ToSlice() []T
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

// iterableStream is a generic stream that is iterated by the iterator returned by the
// supplier function
type iterableStream[T any] struct {
	infinite bool
	supply   iteratorSupplier[T]
}

func (is *iterableStream[T]) iterator() iterator[T] {
	return is.supply()
}

func (is *iterableStream[T]) isInfinite() bool {
	return is.infinite
}

func assertFinite[T any](is Stream[T]) {
	if is.isInfinite() {
		var v T
		panic(fmt.Sprintf("operation not allowed in an infinite Stream[%T]", v))
	}
}
