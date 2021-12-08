package stream

import "github.com/mariomac/gostream/item"

func (bs *iterableStream[T]) ForEach(fn func(T)) {
	next := bs.iterator()
	for in, ok := next(); ok; in, ok = next() {
		fn(in)
	}
}

func (st *iterableStream[T]) ToSlice() []T {
	assertFinite[T](st)
	// TODO: use "count" for better performance
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
	assertFinite(input)
	out := map[K]V{}
	input.ForEach(func(i item.Pair[K, V]) {
		out[i.Key] = i.Val
	})
	return out
}

func (is *iterableStream[T]) Reduce(accumulator func(a, b T) T) (T, bool) {
	assertFinite[T](is)
	next := is.iterator()
	accum, ok := next()
	if !ok {
		return accum, false
	}
	for r, ok := next(); ok; r, ok = next() {
		accum = accumulator(accum, r)
	}
	return accum, true
}
