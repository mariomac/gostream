package stream

import "github.com/mariomac/gostream/item"

func (bs *iterableStream[T]) ForEach(fn func(T)) {
	next := bs.iterator()
	for in, ok := next(); ok; in, ok = next() {
		fn(in)
	}
}

func (st *iterableStream[T]) ToSlice() []T {
	// TODO: use "count" for better performance
	// TODO: assert it's finite
	var res []T
	next := st.iterator()
	for r, ok := next(); ok; r, ok = next() {
		res = append(res, r)
	}
	return res
}

// ToMap returns a map Containing all the item.Pair elements of this Stream, where
// the Key/Val fields of the item.Pair represents the key/value of the map, respectively.
func ToMap[K comparable, V any](input Stream[item.Pair[K, V]]) map[K]V {
	out := map[K]V{}
	input.ForEach(func(i item.Pair[K, V]) {
		out[i.Key] = i.Val
	})
	return out
}
