package stream

import (
	"iter"
	"slices"

	"github.com/mariomac/gostream/item"
	"github.com/mariomac/gostream/order"
)

// ForEach invokes the consumer function for each item of the Stream.
// This function is equivalent to invoking input.ForEach(consumer) as method.
func ForEach[T any](input Stream[T], consumer func(T)) {
	input.ForEach(consumer)
}

func (is *iterableStream[T]) ForEach(consumer func(T)) {
	next := is.iterator()
	for in, ok := next(); ok; in, ok = next() {
		consumer(in)
	}
}

// Iter makes iterableStream compatible with Go's "for ... range" syntax.
// It returns a function that can be used in range loops.
// This function is equivalent to referencing the input.Iter method (without invoking it)
func Iter[T any](input Stream[T]) func(func(int, T) bool) {
	return input.Iter
}

func (is *iterableStream[T]) Iter(yield func(int, T) bool) {
	next := is.iterator()
	idx := 0
	for item, ok := next(); ok; item, ok = next() {
		if !yield(idx, item) {
			return
		}
		idx++
	}
}

// Seq returns the input Stream[T] as a Go standard iter.Seq[T].
// This function is equivalent to referencing the input.Seq method (without invoking it)
func Seq[T any](input Stream[T]) iter.Seq[T] {
	return input.Seq
}

func (is *iterableStream[T]) Seq(yield func(T) bool) {
	next := is.iterator()
	for item, ok := next(); ok; item, ok = next() {
		if !yield(item) {
			return
		}
	}
}

// Seq2 returns the input Stream[item.Pair[K, V]] as a Go standard iter.Seq2[K,V].
// Observe that this function will only work on Streams created from maps (or
// explicitly created as streams of item.Pair[K, V])
func Seq2[K comparable, V any](input Stream[item.Pair[K, V]]) iter.Seq2[K, V] {
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
// This function is equivalent to invoking input.ToSlice() as method.
func ToSlice[T any](input Stream[T]) []T {
	return input.ToSlice()
}

func (is *iterableStream[T]) ToSlice() []T {
	assertFinite[T](is)
	return slices.Collect(is.Seq)
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

// Reduce performs a reduction on the elements of this stream, using an associative
// accumulation function, and returns an value describing the reduced value, if any.
// If no reduced value (e.g. because the stream is empty), the second returned value
// is false.
// This function is equivalent to invoking input.Reduce(accumulator) as method.
func Reduce[T any](input Stream[T], accumulator func(a, b T) T) (T, bool) {
	return input.Reduce(accumulator)
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

// AllMatch returns whether all elements of this stream match the provided predicate.
// If this operation finds an item where the predicate is false, it stops processing
// the rest of the stream.
// This function is equivalent to invoking input.AllMatch(predicate) as method
func AllMatch[T any](input Stream[T], predicate func(T) bool) bool {
	return input.AllMatch(predicate)
}

func (is *iterableStream[T]) AllMatch(predicate func(T) bool) bool {
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
// This function is equivalent to invoking input.AnyMatch(predicate) as method
func AnyMatch[T any](input Stream[T], predicate func(T) bool) bool {
	return input.AnyMatch(predicate)
}

func (is *iterableStream[T]) AnyMatch(predicate func(T) bool) bool {
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
func NoneMatch[T any](input Stream[T], predicate func(T) bool) bool {
	return input.NoneMatch(predicate)
}

func (is *iterableStream[T]) NoneMatch(predicate func(T) bool) bool {
	return !is.AnyMatch(predicate)
}

// Count of elements in this stream.
// This function is equivalent to invoking input.Count() as method.
func Count[T any](input Stream[T]) int {
	return input.Count()
}

func (is *iterableStream[T]) Count() int {
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
// This function is equivalent to invoking input.FindFirst() as method.
func FindFirst[T any](input Stream[T]) (T, bool) {
	return input.FindFirst()
}

func (is *iterableStream[T]) FindFirst() (T, bool) {
	return is.iterator()()
}

// Max returns the maximum element of this stream according to the provided Comparator,
// along with true if the stream is not empty. If the stream is empty, returns the zero
// value along with false.
// This function is equivalent to invoking input.Max(cmp) as method.
func Max[T any](input Stream[T], cmp order.Comparator[T]) (T, bool) {
	return input.Max(cmp)
}

func (is *iterableStream[T]) Max(cmp order.Comparator[T]) (T, bool) {
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
// This function is equivalent to invoking input.Min(cmp) as method.
func Min[T any](input Stream[T], cmp order.Comparator[T]) (T, bool) {
	return input.Min(cmp)
}

func (is *iterableStream[T]) Min(cmp order.Comparator[T]) (T, bool) {
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
