package stream

// Map returns a Stream consisting of the results of individually applying
// the mapper function to each elements of the input Stream.
func Map[IT, OT any](input Stream[IT], mapper func(IT) OT) Stream[OT] {
	return newIterableStream[OT](func() iterator[OT] {
		next := input.iterator()
		return func() (OT, bool) {
			n, ok := next()
			if !ok {
				var zeroVal OT
				return zeroVal, ok
			}
			return mapper(n), true
		}
	})
}

func (bs *abstractStream[T]) Map(mapper func(T) T) Stream[T] {
	return Map[T, T](bs.implementor, mapper)
}

func (as *abstractStream[T]) Filter(predicate func(T) bool) Stream[T] {
	return newIterableStream(func() iterator[T] {
		next := as.implementor.iterator()
		return func() (T, bool) {
			for {
				n, ok := next()
				if !ok {
					var zeroVal T
					return zeroVal, ok
				}
				if predicate(n) {
					return n, true
				}
			}
		}
	})
}
