package stream

import (
	"github.com/mariomac/gostream/order"
)

// Map returns a Stream consisting of the results of individually applying
// the mapper function to each elements of the input Stream.
// When both the input and output type are the same, the operation can be
// invoked as the method input.Map(mapper).
func Map[IT, OT any](input Stream[IT], mapper func(IT) OT) Stream[OT] {
	return &iterableStream[OT]{
		infinite: input.isInfinite(),
		supply: func() iterator[OT] {
			next := input.iterator()
			return func() (OT, bool) {
				n, ok := next()
				if !ok {
					var zeroVal OT
					return zeroVal, false
				}
				return mapper(n), true
			}
		},
	}
}

func (is *iterableStream[T]) Map(mapper func(T) T) Stream[T] {
	return Map[T, T](is, mapper)
}

// Filter returns a Stream consisting of the items of this stream that match the given
// predicate (this is, applying the predicate function over the item returns true).
// This function is equivalent to invoking input.Filter(predicate) as method.
func Filter[T any](input Stream[T], predicate func(T) bool) Stream[T] {
	return input.Filter(predicate)
}

func (is *iterableStream[T]) Filter(predicate func(T) bool) Stream[T] {
	return &iterableStream[T]{
		infinite: is.infinite,
		supply: func() iterator[T] {
			next := is.iterator()
			return func() (T, bool) {
				for {
					n, ok := next()
					if !ok {
						var zeroVal T
						return zeroVal, false
					}
					if predicate(n) {
						return n, true
					}
				}
			}
		}}
}

// Limit returns a stream consisting of the elements of this stream, truncated to
// be no longer than maxSize in length.
// This function is equivalent to invoking input.Limit(maxSize) as method.
func Limit[T any](input Stream[T], maxSize int) Stream[T] {
	return input.Limit(maxSize)
}

func (is *iterableStream[T]) Limit(maxSize int) Stream[T] {
	return &iterableStream[T]{
		infinite: false,
		supply: func() iterator[T] {
			next := is.iterator()
			count := 0
			return func() (T, bool) {
				if count == maxSize {
					var zeroVal T
					return zeroVal, false
				}
				n, ok := next()
				if !ok {
					count = maxSize
					var zeroVal T
					return zeroVal, false
				}
				count++
				return n, true
			}
		}}
}

// Distinct returns a stream consisting of the distinct elements (according to equality operator)
// of the input stream.
func Distinct[T comparable](input Stream[T]) Stream[T] {
	return &iterableStream[T]{supply: func() iterator[T] {
		next := input.iterator()
		elems := map[T]struct{}{}
		return func() (T, bool) {
			for {
				n, ok := next()
				if !ok {
					var zeroVal T
					return zeroVal, false
				}
				if _, ok := elems[n]; !ok {
					elems[n] = struct{}{}
					return n, true
				}
			}
		}
	}}
}

// Sorted returns a stream consisting of the elements of this stream, sorted according
// to the provided order.Comparator.
// This function is equivalent to invoking input.Sorted(comparator) as method.
func Sorted[T any](input Stream[T], comparator order.Comparator[T]) Stream[T] {
	return input.Sorted(comparator)
}
func (is *iterableStream[T]) Sorted(comparator order.Comparator[T]) Stream[T] {
	assertFinite[T](is)
	return &iterableStream[T]{
		infinite: false,
		supply: func() iterator[T] {
			items := is.ToSlice()
			order.SortSlice(items, comparator)
			return func() (T, bool) {
				if len(items) == 0 {
					return finishedIterator[T]()
				}
				n := items[0]
				items = items[1:]
				return n, true
			}
		},
	}
}

// FlatMap returns a stream consisting of the results of replacing each element of this stream
// with the contents of a mapped stream produced by applying the provided mapping function to
// each element. Each mapped stream is closed after its contents have been placed into this
// stream. (If a mapped stream is null an empty stream is used, instead.)
//
// Due to the lazy nature of streams, if any of the mapped streams is infinite it will remain
// unnoticed and some operations (Count, Reduce, Sorted, AllMatch...) will not end.
//
// When both the input and output type are the same, the operation can be
// invoked as the method input.FlatMap(mapper).
func FlatMap[IN, OUT any](input Stream[IN], mapper func(IN) Stream[OUT]) Stream[OUT] {
	return &iterableStream[OUT]{
		supply: func() iterator[OUT] {
			nextFromInputStream := input.iterator()
			var nextFromOutputStream iterator[OUT]
			return func() (OUT, bool) {
				for {
					for nextFromOutputStream == nil {
						// apply the mapper to the current input item and generate an output stream
						// to iterate
						nextInputElem, ok := nextFromInputStream()
						if !ok {
							return finishedIterator[OUT]()
						}
						if outStream := mapper(nextInputElem); outStream != nil {
							nextFromOutputStream = outStream.iterator()
						}
					}
					if outItem, ok := nextFromOutputStream(); ok {
						return outItem, true
					} else {
						// item's resulting outputStream has been iterated. Look for next input item
						nextFromOutputStream = nil
					}
				}
			}
		},
	}
}

func (is *iterableStream[T]) FlatMap(mapper func(T) Stream[T]) Stream[T] {
	return FlatMap[T, T](is, mapper)
}

// Peek peturns a stream consisting of the elements of this stream, additionally performing
// the provided action on each element as elements are consumed from the resulting stream.
// This function is equivalent to invoking input.Peek(consumer) as method.
func Peek[T any](input Stream[T], consumer func(T)) Stream[T] {
	return input.Peek(consumer)
}
func (is *iterableStream[T]) Peek(consumer func(T)) Stream[T] {
	return &iterableStream[T]{
		infinite: is.isInfinite(),
		supply: func() iterator[T] {
			next := is.iterator()
			return func() (T, bool) {
				n, ok := next()
				if !ok {
					var zeroVal T
					return zeroVal, false
				}
				consumer(n)
				return n, true
			}
		},
	}
}

// Skip returns a stream consisting of the remaining elements of this stream after discarding
// the first n elements of the stream.
// This function is equivalent to invoking input.Skip(n) as method.
func Skip[T any](input Stream[T], n int) Stream[T] {
	return input.Skip(n)
}
func (is *iterableStream[T]) Skip(n int) Stream[T] {
	return &iterableStream[T]{
		infinite: is.isInfinite(),
		supply: func() iterator[T] {
			next := is.iterator()
			skipped := 0
			return func() (T, bool) {
				var it T
				var ok bool
				for it, ok = next(); ok && skipped < n; it, ok = next() {
					skipped++
				}
				return it, ok
			}
		},
	}
}
