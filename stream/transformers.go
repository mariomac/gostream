package stream

import (
	"github.com/mariomac/gostream/order"
)

// Map returns a Stream consisting of the results of individually applying
// the mapper function to each elements of the input Stream.
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

func (bs *iterableStream[T]) Map(mapper func(T) T) Stream[T] {
	return Map[T, T](bs, mapper)
}

func (as *iterableStream[T]) Filter(predicate func(T) bool) Stream[T] {
	return &iterableStream[T]{
		infinite: as.infinite,
		supply: func() iterator[T] {
			next := as.iterator()
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

func (as *iterableStream[T]) Limit(maxSize int) Stream[T] {
	return &iterableStream[T]{
		infinite: false,
		supply: func() iterator[T] {
			next := as.iterator()
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
