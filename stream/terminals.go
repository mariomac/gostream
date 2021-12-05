package stream

func (bs *iterableStream[T]) ForEach(fn func(T)) {
	next := bs.iterator()
	for in, ok := next(); ok; in, ok = next() {
		fn(in)
	}
}

func (st *iterableStream[T]) ToSlice() []T {
	// TODO: use "count" for better performance
	var res []T
	next := st.iterator()
	for r, ok := next(); ok; r, ok = next() {
		res = append(res, r)
	}
	return res
}
