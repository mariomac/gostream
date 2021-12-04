// Package str provides type-safe streams, functional helper tools and processing
// operations
package str

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
	attach() Stream[T]

	// stream operations

	// Filter seturns a Stream consisting of the items of this stream that match the given
	// predicate (this is, applying the predicate function over the item returns true)
	Filter(predicate func(T) bool) Stream[T]
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

type iterator[T any] interface {
	// it there are more items to iterate, returns the next item and true.
	// if the iterator has iterated all the stream items, returns the zero value and false.
	next() (T, bool)
}

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
	return &sliceIterator[T]{items: si.items}
}

// iterator for slices
type sliceIterator[T any] struct {
	items []T
}

func (si *sliceIterator[T]) next() (T, bool) {
	if len(si.items) == 0 {
		var zeroValue T
		return zeroValue, false
	}
	n := si.items[0]
	si.items = si.items[1:]
	return n, true
}

// mapperStream takes the items from an input stream and
// transforms them according to a mapper function.
type mapperStream[IN, OUT any] struct {
	abstractStream[OUT]
	mapper func(IN) OUT
	input  Stream[IN]
}

func (si *mapperStream[IN, OUT]) attach() Stream[OUT] {
	si.abstractStream.implementor = si
	return si
}

func (si *mapperStream[IN, OUT]) iterator() iterator[OUT] {
	return &mapperIterator[IN, OUT]{
		input:  si.input.iterator(),
		mapper: si.mapper,
	}
}

// iterator for a mapperStream
type mapperIterator[IN, OUT any] struct {
	mapper func(IN) OUT
	input  iterator[IN]
}

func (c *mapperIterator[IN, OUT]) next() (OUT, bool) {
	n, ok := c.input.next()
	if !ok {
		var zeroVal OUT
		return zeroVal, ok
	}
	return c.mapper(n), true
}

// filterStream takes the items from an input stream and only forwards them to the next
// stage of the pipeline if they fulfill a given predicate.
type filterStream[T any] struct {
	abstractStream[T]
	predicate func(T) bool
	input     Stream[T]
}

func (si *filterStream[T]) attach() Stream[T] {
	si.abstractStream.implementor = si
	return si
}

func (si *filterStream[T]) iterator() iterator[T] {
	return &filterIterator[T]{
		input:     si.input.iterator(),
		predicate: si.predicate,
	}
}

// iterator for a filterStream
type filterIterator[T any] struct {
	predicate func(T) bool
	input     iterator[T]
}

func (c *filterIterator[T]) next() (T, bool) {
	for {
		n, ok := c.input.next()
		if !ok {
			var zeroVal T
			return zeroVal, ok
		}
		if c.predicate(n) {
			return n, true
		}
	}
}
