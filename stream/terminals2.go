package stream

import (
	"iter"
	"maps"
)

// ForEach invokes the consumer function for each item of the Stream2.
func (is *Stream2[K, V]) ForEach(consumer func(K, V)) {
	next := is.iterator()
	for k, v, ok := next(); ok; k, v, ok = next() {
		consumer(k, v)
	}
}

// Seq2 returns the input *Stream2[K, V] as a Go standard iter.Seq2[K, V].
func (is *Stream2[K, V]) Seq2() iter.Seq2[K, V] {
	next := is.iterator()
	return func(yield func(K, V) bool) {
		for k, v, ok := next(); ok; k, v, ok = next() {
			if !yield(k, v) {
				return
			}
		}
	}
}

// ToMap returns a map Containing all the key-value pair elements of this Stream2.
func (is *Stream2[K, V]) ToMap() map[K]V {
	return maps.Collect(is.Seq2())
}

// AllMatch returns whether all key-value pairs of this stream match the provided predicate.
// If this operation finds an item where the predicate is false, it stops processing
// the rest of the stream.
func (is *Stream2[K, V]) AllMatch(predicate func(K, V) bool) bool {
	next := is.iterator()
	for k, v, ok := next(); ok; k, v, ok = next() {
		if !predicate(k, v) {
			return false
		}
	}
	return true
}

// AnyMatch returns whether any elements of this stream match the provided predicate.
// If this operation finds an item where the predicate is true, it stops processing
// the rest of the stream.
func (is *Stream2[K, V]) AnyMatch(predicate func(K, V) bool) bool {
	next := is.iterator()
	for k, v, ok := next(); ok; k, v, ok = next() {
		if predicate(k, v) {
			return true
		}
	}
	return false
}

// NoneMatch returns whether no elements of this stream match the provided predicate.
// If this operation finds an item where the predicate is true, it stops processing
// the rest of the stream.
func (is *Stream2[K, V]) NoneMatch(predicate func(K, V) bool) bool {
	return !is.AnyMatch(predicate)
}

// Count of elements in this stream.
func (is *Stream2[K, V]) Count() int {
	count := 0
	next := is.iterator()
	for _, _, ok := next(); ok; _, _, ok = next() {
		count++
	}
	return count
}

// FindFirst returns the first key-value pair of this Stream2 along with true or, if the
// stream is empty, the zero value of the inner type along with false.
func (is *Stream2[K, V]) FindFirst() (K, V, bool) {
	return is.iterator()()
}

// Max returns the maximum key-value pair of this stream according to the provided Comparator2,
// along with true if the stream is not empty. If the stream is empty, returns the zero
// value along with false.
func (is *Stream2[K, V]) Max(cmp Comparator2[K, V]) (K, V, bool) {
	next := is.iterator()
	maxK, maxV, ok := next()
	if !ok {
		return maxK, maxV, false
	}
	for k, v, ok := next(); ok; k, v, ok = next() {
		if cmp(k, v, maxK, maxV) > 0 {
			maxK, maxV = k, v
		}
	}
	return maxK, maxV, true
}

// Min returns the minimum element of this stream according to the provided Comparator,
// along with true if the stream is not empty. If the stream is empty, returns the zero
// value along with false.
func (is *Stream2[K, V]) Min(cmp Comparator2[K, V]) (K, V, bool) {
	next := is.iterator()
	minK, minV, ok := next()
	if !ok {
		return minK, minV, false
	}
	for k, v, ok := next(); ok; k, v, ok = next() {
		if cmp(k, v, minK, minV) < 0 {
			minK, minV = k, v
		}
	}
	return minK, minV, true
}
