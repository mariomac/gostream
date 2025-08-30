package iters

import (
	"cmp"
	"iter"

	"github.com/mariomac/gostream/order"
)

// ForEach invokes the consumer function for each item of the iter.Seq.
func ForEach[T any](input iter.Seq[T], consumer func(T)) {
	for v := range input {
		consumer(v)
	}
}

// ForEach2 invokes the consumer function for each pair of items of the iter.Seq2
func ForEach2[T1, T2 any](input iter.Seq2[T1, T2], consumer func(T1, T2)) {
	for v1, v2 := range input {
		consumer(v1, v2)
	}
}

// Reduce performs a reduction on the elements of the input Seq, using an associative
// accumulation function, and returns a value describing the reduced value, if any.
// If no reduced value (e.g. because the iter.Seq is empty), the second returned value
// is false.
func Reduce[T any](input iter.Seq[T], accumulator func(a, b T) T) (T, bool) {
	pull, _ := iter.Pull(input)
	accum, ok := pull()
	if !ok {
		return accum, false
	}
	for r, ok := pull(); ok; r, ok = pull() {
		accum = accumulator(accum, r)
	}
	return accum, true
}

// AllMatch returns whether all elements of this iter.Seq match the provided predicate.
// If this operation finds an item where the predicate is false, it stops processing
// the rest of the iter.Seq.
func AllMatch[T any](input iter.Seq[T], predicate func(T) bool) bool {
	for r := range input {
		if !predicate(r) {
			return false
		}
	}
	return true
}

// AnyMatch returns whether any elements of the iter.Seq match the provided predicate.
// If this operation finds an item where the predicate is true, it stops processing
// the rest of the iter.Seq.
func AnyMatch[T any](input iter.Seq[T], predicate func(T) bool) bool {
	for r := range input {
		if predicate(r) {
			return true
		}
	}
	return false
}

// NoneMatch returns whether no elements of the iter.Seq match the provided predicate.
// If this operation finds an item where the predicate is true, it stops processing
// the rest of the iter.Seq.
func NoneMatch[T any](input iter.Seq[T], predicate func(T) bool) bool {
	return !AnyMatch(input, predicate)
}

// Count of elements in the iter.Seq.
func Count[T any](input iter.Seq[T]) int {
	c := 0
	for range input {
		c++
	}
	return c
}

// FindFirst returns the first element of this iter.Seq along with true or, if the
// iter.Seq is empty, the zero value of the inner type along with false.
func FindFirst[T any](input iter.Seq[T]) (T, bool) {
	for i := range input {
		return i, true
	}
	var t T
	return t, false
}

// Max returns the maximum element of the iter.Seq.
// along with true if the iter.Seq is not empty. If the iter.Seq is empty, returns the zero
// value along with false.
func Max[T cmp.Ordered](input iter.Seq[T]) (T, bool) {
	it, _ := iter.Pull(input)
	max, ok := it()
	if !ok {
		return max, false
	}
	for n, ok := it(); ok; n, ok = it() {
		if n > max {
			max = n
		}
	}
	return max, true
}

// MaxFunc returns the maximum element of the iter.Seq according to the provided Comparator,
// along with true if the iter.Seq is not empty. If the iter.Seq is empty, returns the zero
// value along with false.
func MaxFunc[T any](input iter.Seq[T], cmp order.Comparator[T]) (T, bool) {
	it, _ := iter.Pull(input)
	max, ok := it()
	if !ok {
		return max, false
	}
	for n, ok := it(); ok; n, ok = it() {
		if cmp(n, max) > 0 {
			max = n
		}
	}
	return max, true
}

// Min returns the minimum element of the iter.Seq,
// along with true if the iter.Seq is not empty. If the iter.Seq is empty, returns the zero
// value along with false.
func Min[T cmp.Ordered](input iter.Seq[T]) (T, bool) {
	next, _ := iter.Pull(input)
	min, ok := next()
	if !ok {
		return min, false
	}
	for n, ok := next(); ok; n, ok = next() {
		if n < min {
			min = n
		}
	}
	return min, true
}

// MinFunc returns the minimum element of the iter.Seq according to the provided Comparator,
// along with true if the iter.Seq is not empty. If the iter.Seq is empty, returns the zero
// value along with false.
func MinFunc[T any](input iter.Seq[T], cmp order.Comparator[T]) (T, bool) {
	next, _ := iter.Pull(input)
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
