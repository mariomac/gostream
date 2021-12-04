package str

func (bs *abstractStream[T]) ForEach(fn func(T)) {
	it := bs.implementor.iterator()
	for in, ok := it.next(); ok; in, ok = it.next() {
		fn(in)
	}
}

func (st *abstractStream[T]) ToSlice() []T {
	// TODO: use "count" for better performance
	var res []T
	it := st.implementor.iterator()
	for r, ok := it.next(); ok; r, ok = it.next() {
		res = append(res, r)
	}
	return res
}

func (st *sliceStream[T]) ToSlice() []T {
	return st.items
}
