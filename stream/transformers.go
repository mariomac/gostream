package stream

import (
	"slices"
)

// Map returns a Stream consisting of the results of individually applying
// the mapper function to each element of the target Stream.
func (is *Stream[T]) Map[OT any](mapper func(T) OT) *Stream[OT] {
	return &Stream[OT]{
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

// IsCmp is an auxiliary function that has to be passed as argument to the stream.Stream's Distinct method,
// to help with some limitations of Go generic type system.
func Cmp[T comparable](i T) T {
	return i
}

// Distinct returns a stream consisting of the distinct elements (according to equality operator)
// of the input stream.
// Due to the limitations of the Go generics system, you need to pass Cmp[T]
// as argument to let compiler checking at runtime that the stream type T is comparable.
func (is *Stream[T]) Distinct[O comparable](cmp func(T) O) *Stream[O] {
	return &Stream[O]{supply: func() iterator[O] {
		next := is.iterator()
		elems := map[O]struct{}{}
		return func() (O, bool) {
			for {
				nt, ok := next()
				if !ok {
					return *new(O), false
				}
				n := cmp(nt)
				if _, ok := elems[n]; !ok {
					elems[n] = struct{}{}
					return n, true
				}
			}
		}
	}}
}

// Sorted Stream a stream consisting of the elements of the target Stream, sorted according
// to the provided Comparator.
func (is *Stream[T]) Sorted(comparator Comparator[T]) *Stream[T] {
	return &Stream[T]{
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
//
// Due to the lazy nature of streams, if any of the mapped streams is infinite it will remain
// unnoticed and some operations (Count, Reduce, Sorted, AllMatch...) will not end.
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

// Peek returns a stream consisting of the elements of this stream, additionally performing
// the provided action on each element as elements are consumed from the resulting stream.
func (is *Stream[T]) Peek(consumer func(T)) *Stream[T] {
	return &Stream[T]{
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

// MapStream2 returns a Stream2 (two output types for key and value) consisting of the results of individually applying
// the mapper function to each key-value pair of the target Stream.
func (is *Stream[T]) MapStream2[K comparable, V any](mapper func(i T) (K, V)) *Stream2[K, V] {
	return &Stream2[K, V]{
		supply: func() iterator2[K, V] {
			next := is.iterator()
			return func() (K, V, bool) {
				n, ok := next()
				if !ok {
					return *new(K), *new(V), false
				}
				k, v := mapper(n)
				return k, v, true
			}
		},
	}
}
