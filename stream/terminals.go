package stream

func (bs *abstractStream[T]) ForEach(fn func(T)) {
	next := bs.implementor.iterator()
	for in, ok := next(); ok; in, ok = next() {
		fn(in)
	}
}

func (st *abstractStream[T]) ToSlice() []T {
	// TODO: use "count" for better performance
	var res []T
	next := st.implementor.iterator()
	for r, ok := next(); ok; r, ok = next() {
		res = append(res, r)
	}
	return res
}

func (st *sliceStream[T]) ToSlice() []T {
	return st.items
}
