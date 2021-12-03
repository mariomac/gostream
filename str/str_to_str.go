package str

// some functions are implemented here and referenced from the
// method and others are implemented in the method and referenced
// here
// the criteria is: type-changing functions are more likely to be
// invoked as functions (e.g. Map). Type-keeping functions are
// more likely to be invoked as methods

func Map[IT any, OT any](is Stream[IT], fn func(IT) OT) Stream[OT] {
	return (&connectedStream[IT, OT]{
		appliedFn: fn,
		Input:     is,
	}).attach()
}

func (bs *abstractStream[T]) Map(fn func(T) T) Stream[T] {
	return Map[T, T](bs.implementor, fn)
}

func Filter[T any](is Stream[T], predicate func(T) bool) Stream[T] {
	return is.Filter(predicate)
}

func (as *abstractStream[T]) Filter(predicate func(T) bool) Stream[T] {
	return (&filterStream[T]{predicate: predicate, Input: as.implementor}).attach()
}
