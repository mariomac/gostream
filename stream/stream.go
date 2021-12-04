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

	// connects an abstract stream with its implementor. Stream implementers should implement this
	// and invoke it after each stream instantiation.
	// This allows to automatically implement all the abstract stream operations
	// through the abstractStream
	// For fluent invocation, it returns the target stream
	attach() Stream[T]

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

// if there are more items to iterate, returns the next item and true.
// if the iterator has iterated all the stream items, returns the zero value and false.
// this function will usually depend on an external status (e.g. can be a struct method or
// a function literal that rely on outer variables)
type iterator[T any] func() (T, bool)

// abstractStream provides generic implementations of the Stream methods so implementors don't
// need to reimplement all of them. However, implementors could override the default methods
// for optimizing them according their own nature.
type abstractStream[T any] struct {
	implementor Stream[T]
}

// sliceStream is a stream whose items are stored in a slice
type sliceStream[T any] struct {
	abstractStream[T]
	items []T
}

// sets itself as an implementor of the stream so the abstractStream will
// invoke the implementor methods instead of its own methods.
func (si *sliceStream[T]) attach() Stream[T] {
	si.abstractStream.implementor = si
	return si
}

func (si *sliceStream[T]) iterator() iterator[T] {
	items := si.items
	return func() (T, bool) {
		if len(items) == 0 {
			var zeroValue T
			return zeroValue, false
		}
		n := items[0]
		items = items[1:]
		return n, true
	}
}

type iteratorSupplier[T any] func() iterator[T]

// iterableStream is a generic stream that is iterated by the iterator returned by the
// supplier function
type iterableStream[T any] struct {
	abstractStream[T]
	supply iteratorSupplier[T]
}

func newIterableStream[T any](supplier iteratorSupplier[T]) *iterableStream[T] {
	is := &iterableStream[T]{supply: supplier}
	is.attach()
	return is
}

func (is *iterableStream[T]) attach() Stream[T] {
	is.abstractStream.implementor = is
	return is
}

func (is *iterableStream[T]) iterator() iterator[T] {
	return is.supply()
}
