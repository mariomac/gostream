package stream

// Of creates an Stream from a variable number of elements that are passed as
// arguments.
func Of[T any](elems ...T) Stream[T] {
	return OfSlice(elems)
}

// OfSlice creates a Stream from a slice.
func OfSlice[T any](elems []T) Stream[T] {
	return (&sliceStream[T]{items: elems}).attach()
}
