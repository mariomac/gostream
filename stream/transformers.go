package stream

import (
	"reflect"
	"slices"

	"github.com/mariomac/gostream/order"
)

// Map returns a Stream consisting of the results of individually applying
// the mapper function to each elements of the target Stream.
func (is *Stream[T]) Map[OT any](mapper func(T) OT) *Stream[OT] {
	return &Stream[OT]{
		infinite: is.isInfinite(),
		supply: func() iterator[OT] {
			next := is.iterator()
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

// Filter returns a Stream consisting of the items of the target stream that match the given
// predicate (this is, applying the predicate function over the item returns true).
func (is *Stream[T]) Filter(predicate func(T) bool) *Stream[T] {
	return &Stream[T]{
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
		},
	}
}

// Limit returns a stream consisting of the elements of this stream, truncated to
// be no longer than maxSize in length.
func (is *Stream[T]) Limit(maxSize int) *Stream[T] {
	return &Stream[T]{
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
		},
	}
}

// Distinct returns a stream consisting of the distinct elements according to the equality
// operator ==, preserving encounter order. Values that cannot safely be compared are passed
// through unchanged, including every occurrence of such a value.
func (is *Stream[T]) Distinct() *Stream[T] {
	typ := reflect.TypeFor[T]()
	if !typ.Comparable() {
		return &Stream[T]{infinite: is.isInfinite(), supply: is.supply}
	}
	// Interfaces, including those nested in structs or arrays, may contain
	// values that cannot be used as map keys. Scalar types need no per-value check.
	checkValues := typ.Kind() == reflect.Interface || typ.Kind() == reflect.Struct || typ.Kind() == reflect.Array
	return &Stream[T]{infinite: is.isInfinite(), supply: func() iterator[T] {
		next := is.iterator()
		elems := map[any]struct{}{}
		return func() (T, bool) {
			for {
				n, ok := next()
				if !ok {
					var zeroVal T
					return zeroVal, false
				}
				if checkValues && !reflect.ValueOf(&n).Elem().Comparable() {
					return n, true
				}
				key := any(n)
				if _, ok := elems[key]; !ok {
					elems[key] = struct{}{}
					return n, true
				}
			}
		}
	}}
}

// Sorted Stream a stream consisting of the elements of the target Stream, sorted according
// to the provided order.Comparator.
func (is *Stream[T]) Sorted(comparator order.Comparator[T]) *Stream[T] {
	assertFinite[T](is)
	return &Stream[T]{
		infinite: false,
		supply: func() iterator[T] {
			items := is.ToSlice()
			slices.SortFunc(items, comparator)
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
func (is *Stream[IN]) FlatMap[OUT any](mapper func(IN) *Stream[OUT]) *Stream[OUT] {
	return &Stream[OUT]{
		supply: func() iterator[OUT] {
			nextFromInputStream := is.iterator()
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

// Peek peturns a stream consisting of the elements of this stream, additionally performing
// the provided action on each element as elements are consumed from the resulting stream.
func (is *Stream[T]) Peek(consumer func(T)) *Stream[T] {
	return &Stream[T]{
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
func (is *Stream[T]) Skip(n int) *Stream[T] {
	return &Stream[T]{
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
