package stream

// Map returns a Stream consisting of the results of individually applying
// the mapper function to each elements of the input Stream.
func Map[IT any, OT any](input Stream[IT], mapper func(IT) OT) Stream[OT] {
	return (&mapperStream[IT, OT]{
		mapper: mapper,
		input:  input,
	}).attach()
}

func (bs *abstractStream[T]) Map(mapper func(T) T) Stream[T] {
	return Map[T, T](bs.implementor, mapper)
}

func (as *abstractStream[T]) Filter(predicate func(T) bool) Stream[T] {
	return (&filterStream[T]{predicate: predicate, input: as.implementor}).attach()
}
