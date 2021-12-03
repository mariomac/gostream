package str

func Of[T any](elems ...T) Stream[T] {
	return OfSlice(elems)
}

func OfSlice[T any](elems []T) Stream[T] {
	return (&sliceStream[T]{Items: elems}).attach()
}
