package stream

// Map returns a Stream consisting of the results of individually applying
// the mapper function to each elements of the input Stream.
func Map[IT, OT any](input Stream[IT], mapper func(IT) OT) Stream[OT] {
	return newIterableStream[OT](func() iterator[OT] {
		return &mapperIterator[IT, OT]{
			mapper: mapper,
			input:  input.iterator(),
		}
	})
}

func (bs *abstractStream[T]) Map(mapper func(T) T) Stream[T] {
	return Map[T, T](bs.implementor, mapper)
}

func (as *abstractStream[T]) Filter(predicate func(T) bool) Stream[T] {
	return newIterableStream(func() iterator[T] {
		return &filterIterator[T]{
			predicate: predicate,
			input:     as.implementor.iterator(),
		}
	})
}
