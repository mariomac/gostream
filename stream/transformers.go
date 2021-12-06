package stream

import (
	"github.com/mariomac/gostream/order"
)

// Map returns a Stream consisting of the results of individually applying
// the mapper function to each elements of the input Stream.
func Map[IT, OT any](input Stream[IT], mapper func(IT) OT) Stream[OT] {
	return &iterableStream[OT]{supply: func() iterator[OT] {
		next := input.iterator()
		return func() (OT, bool) {
			n, ok := next()
			if !ok {
				var zeroVal OT
				return zeroVal, false
			}
			return mapper(n), true
		}
	}}
}

func (bs *iterableStream[T]) Map(mapper func(T) T) Stream[T] {
	return Map[T, T](bs, mapper)
}

func (as *iterableStream[T]) Filter(predicate func(T) bool) Stream[T] {
	return &iterableStream[T]{supply: func() iterator[T] {
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
	return &iterableStream[T]{supply: func() iterator[T] {
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

func (cs *comparableStream[T]) Distinct() Comparable[T] {
	return &comparableStream[T]{iterableStream[T]{supply: func() iterator[T] {
		next := cs.iterator()
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
	}}}
}

// TODO: panic if trying to sort an infinite stream
func (is *iterableStream[T]) Sorted(comparator order.Comparator[T]) Stream[T] {
	return &iterableStream[T]{supply: func() iterator[T] {
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
	}}
}
