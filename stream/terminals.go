package stream

import (
	"iter"
	"slices"

	"github.com/mariomac/gostream/item"
	"github.com/mariomac/gostream/order"
)

// ForEach invokes the consumer function for each item of the Stream.
func (is *Stream[T]) ForEach(consumer func(T)) {
	next := is.iterator()
	for in, ok := next(); ok; in, ok = next() {
		consumer(in)
	}
}

// Iter makes Stream compatible with Go's "for ... range" syntax.
// It returns a function that can be used in range loops.
// The difference with the Seq method is that Iter returns an
// iter.Seq2[int, T] whose first value is the position of the sequence
// in the stream.
func (is *Stream[T]) Iter() iter.Seq2[int, T] {
	next := is.iterator()
	return func(yield func(int, T) bool) {
		idx := 0
		for item, ok := next(); ok; item, ok = next() {
			if !yield(idx, item) {
				return
			}
			idx++
		}
	}
}

// Seq returns the input *Stream[T] as a Go standard iter.Seq[T].
// The difference with the Iter method is that the returned iter.Seq
// only iterates through the ordered values, without returning its
// position within the stream
func (is *Stream[T]) Seq() iter.Seq[T] {
	next := is.iterator()
	return func(yield func(T) bool) {
		for item, ok := next(); ok; item, ok = next() {
			if !yield(item) {
				return
			}
		}
	}
}

// Seq2 returns the input *Stream[item.Pair[K, V]] as a Go standard iter.Seq2[K,V].
// Observe that this function will only work on Streams created from maps (or
// explicitly created as streams of item.Pair[K, V])
// TODO esto hay que moverlo
func Seq2[K comparable, V any](input *Stream[item.Pair[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		next := input.iterator()
		for item, ok := next(); ok; item, ok = next() {
			if !yield(item.Key, item.Val) {
				return
			}
		}
	}
}

// ToSlice returns a Slice Containing all the elements of this Stream.
func (is *Stream[T]) ToSlice() []T {
	assertFinite[T](is)
	return slices.Collect(is.Seq())
}

// ToMap returns a map Containing all the item.Pair elements of this Stream, where
// the Key/Val fields of the item.Pair represents the key/value of the map, respectively.
// TODO esto hay que moverlo
func ToMap[K comparable, V any](input *Stream[item.Pair[K, V]]) map[K]V {
	assertFinite(input)
	out := map[K]V{}
	input.ForEach(func(i item.Pair[K, V]) {
		out[i.Key] = i.Val
	})
	return out
}

// Reduce performs a reduction on the elements of this stream, using an associative
// accumulation function, and returns an value describing the reduced value, if any.
// If no reduced value (e.g. because the stream is empty), the second returned value
// is false.
func (is *Stream[T]) Reduce(accumulator func(a, b T) T) (T, bool) {
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

// AllMatch returns whether all elements of this stream match the provided predicate.
// If this operation finds an item where the predicate is false, it stops processing
// the rest of the stream.
func (is *Stream[T]) AllMatch(predicate func(T) bool) bool {
	assertFinite[T](is)
	next := is.iterator()
	for r, ok := next(); ok; r, ok = next() {
		if !predicate(r) {
			return false
		}
	}
	return true
}

// AnyMatch returns whether any elements of this stream match the provided predicate.
// If this operation finds an item where the predicate is true, it stops processing
// the rest of the stream.
func (is *Stream[T]) AnyMatch(predicate func(T) bool) bool {
	assertFinite[T](is)
	next := is.iterator()
	for r, ok := next(); ok; r, ok = next() {
		if predicate(r) {
			return true
		}
	}
	return false
}

// NoneMatch returns whether no elements of this stream match the provided predicate.
// If this operation finds an item where the predicate is true, it stops processing
// the rest of the stream.
// This function is equivalent to invoking input.NoneMatch(predicate) as method.
func (is *Stream[T]) NoneMatch(predicate func(T) bool) bool {
	return !is.AnyMatch(predicate)
}

// Count of elements in this stream.
// This function is equivalent to invoking input.Count() as method.
func (is *Stream[T]) Count() int {
	assertFinite[T](is)
	count := 0
	next := is.iterator()
	for _, ok := next(); ok; _, ok = next() {
		count++
	}
	return count
}

// FindFirst returns the first element of this Stream along with true or, if the
// stream is empty, the zero value of the inner type along with false.
func (is *Stream[T]) FindFirst() (T, bool) {
	return is.iterator()()
}

// Max returns the maximum element of this stream according to the provided Comparator,
// along with true if the stream is not empty. If the stream is empty, returns the zero
// value along with false.
func (is *Stream[T]) Max(cmp order.Comparator[T]) (T, bool) {
	assertFinite[T](is)
	next := is.iterator()
	max, ok := next()
	if !ok {
		return max, false
	}
	for n, ok := next(); ok; n, ok = next() {
		if cmp(n, max) > 0 {
			max = n
		}
	}
	return max, true
}

// Min returns the minimum element of this stream according to the provided Comparator,
// along with true if the stream is not empty. If the stream is empty, returns the zero
// value along with false.
func (is *Stream[T]) Min(cmp order.Comparator[T]) (T, bool) {
	assertFinite[T](is)
	next := is.iterator()
	min, ok := next()
	if !ok {
		return min, false
	}
	for n, ok := next(); ok; n, ok = next() {
		if cmp(n, min) < 0 {
			min = n
		}
	}
	return min, true
}
