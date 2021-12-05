package stream

// Of creates an Stream from a variable number of elements that are passed as
// arguments.
func Of[T any](elems ...T) Stream[T] {
	return OfSlice(elems)
}

// OfSlice creates a Stream from a slice.
func OfSlice[T any](elems []T) Stream[T] {
	return &iterableStream[T]{supply: func() iterator[T] {
		items := elems
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

// Generate an infinite sequential stream where each element is generated by the provided supplier function.
// Due to the stateful nature of the supplier, multiple operations towards the same stream might provide
// different results.
func Generate[T any](supplier func() T) Stream[T] {
	return &iterableStream[T]{supply: func() iterator[T] {
		return func() (T, bool) {
			return supplier(), true
		}
	}}
}

// Iterate returns an infinite sequential ordered Stream produced by iterative application of a function
// f to an initial element seed, producing a Stream consisting of seed, f(seed), f(f(seed)), etc.
// The first element (position 0) in the Stream will be the provided seed. For n > 0, the element at
// position n, will be the result of applying the function f to the element at position n - 1.
// Due to the stateful nature of the supplier, multiple operations towards the same stream might provide
// different results.
func Iterate[T any](seed T, f func(T) T) Stream[T] {
	return &iterableStream[T]{supply: func() iterator[T] {
		lastElement := seed
		return func() (T, bool) {
			i := lastElement
			lastElement = f(lastElement)
			return i, true
		}
	}}
}

// Concat creates a lazily concatenated stream whose elements are all the elements of the first stream
// followed by all the elements of the second stream.
func Concat[T any](a, b Stream[T]) Stream[T] {
	return &iterableStream[T]{supply: func() iterator[T] {
		first := true
		next := a.iterator()
		return func() (T, bool) {
			n, ok := next()
			if ok {
				return n, true
			}
			if first {
				first = false
				next = b.iterator()
			} else {
				next = finishedIterator[T]
			}
			return next()
		}
	}}
}
