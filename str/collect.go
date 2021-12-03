package str

func (st *abstractStream[T]) AsSlice() []T {
	// TODO: use "count" for better performance
	var res []T
	it := st.implementor.iterator()
	for r, ok := it.next(); ok; r, ok = it.next() {
		res = append(res, r)
	}
	return res
}

func (st *sliceStream[T]) AsSlice() []T {
	return st.Items
}
