// package stream provides type-safe streams, functional helper tools and processing
// operations
package stream

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

	// stream operations

	// Filter seturns a Stream consisting of the items of this stream that match the given
	// predicate (this is, applying the predicate function over the item returns true)
	Filter(predicate func(T) bool) Stream[T]

	// Limit returns a stream consisting of the elements of this stream, truncated to
	// be no longer than maxSize in length.
	Limit(maxSize int) Stream[T]

	// Map returns a Stream consisting of the results of individually applying
	// the mapper function to each element of this Stream. The argument and return
	// value of the mapper function must belong to the same type. If you need that
	// the input and output Stream contain elements from different types, you need
	// to invoke the standalone function Map[IN,OUT](Stream[IN], func(IN)OUT) Stream[OUT].
	Map(mapper func(T) T) Stream[T]

	// terminal operations

	// Foreach invokes the consumer function for each item of the Stream.
	ForEach(consumer func(T))

	// ToSlice returns a Slice Containing all the elements of this Stream.
	ToSlice() []T
}

// ComparableStream adds functionalities to a Stream that would require comparing the
// items of the stream between them, so the elements type must fulfill the 'comparable'
// constraint (e.g. define the == and != operators)
type ComparableStream[T comparable] interface {
	Stream[T]
	// Distinct returns a stream consisting of the distinct elements (according to equality operator)
	// of the input stream.
	Distinct() ComparableStream[T]
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
	supply iteratorSupplier[T]
}

func (is *iterableStream[T]) iterator() iterator[T] {
	return is.supply()
}

type comparableStream[T comparable] struct {
	iterableStream[T]
}
