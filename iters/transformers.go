package iters

import (
	"iter"
)

// Map returns a iter.Seq consisting of the results of individually applying
// the mapper function to each elements of the input iter.Seq.
func Map[IT, OT any](input iter.Seq[IT], mapper func(IT) OT) iter.Seq[OT] {
	return func(yield func(OT) bool) {
		for i := range input {
			if !yield(mapper(i)) {
				return
			}
		}
	}
}

// Filter returns a iter.Seq consisting of the items of this stream that match the given
// predicate (this is, applying the predicate function over the item returns true).
func Filter[T any](input iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range input {
			if predicate(i) {
				if !yield(i) {
					return
				}
			}
		}
	}
}

// Limit returns a stream consisting of the elements of this stream, truncated to
// be no longer than maxSize in length.
// This function is equivalent to invoking input.Limit(maxSize) as method.
func Limit[T any](input iter.Seq[T], maxSize int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		for i := range input {
			if count == maxSize {
				return
			}
			if !yield(i) {
				return
			}
			count++
		}
	}
}

// Distinct returns a stream consisting of the distinct elements (according to equality operator)
// of the input stream.
// This function needs to internally store the previous distinct elements in memory, so it might
// not be suitable for large or infinite streams with high variability in their items.
func Distinct[T comparable](input iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		elems := map[T]struct{}{}
		for i := range input {
			if _, ok := elems[i]; !ok {
				elems[i] = struct{}{}
				if !yield(i) {
					return
				}
			}
		}
	}
}

// FlatMap returns a stream consisting of the results of replacing each element of this stream
// with the contents of a mapped stream produced by applying the provided mapping function to
// each element. Each mapped stream is closed after its contents have been placed into this
// stream. (If a mapped stream is null an empty stream is used, instead.)
//
// Due to the lazy nature of streams, if any of the mapped streams is infinite it will remain
// unnoticed and some operations (Count, Reduce, Sorted, AllMatch...) will not end.
//
// When both the input and output type are the same, the operation can be
// invoked as the method input.FlatMap(mapper).
func FlatMap[IN, OUT any](input iter.Seq[IN], mapper func(IN) iter.Seq[OUT]) iter.Seq[OUT] {
	return func(yield func(OUT) bool) {
		for i := range input {
			for j := range mapper(i) {
				if !yield(j) {
					return
				}
			}
		}
	}
}

// Peek peturns a stream consisting of the elements of this stream, additionally performing
// the provided action on each element as elements are consumed from the resulting stream.
// This function is equivalent to invoking input.Peek(consumer) as method.
func Peek[T any](input iter.Seq[T], consumer func(T)) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range input {
			consumer(i)
			if !yield(i) {
				return
			}
		}
	}
}

// Skip returns a stream consisting of the remaining elements of this stream after discarding
// the first n elements of the stream.
// This function is equivalent to invoking input.Skip(n) as method.
func Skip[T any](input iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, _ := iter.Pull(input)
		var it T
		var ok bool
		skipped := 0
		for it, ok = next(); ok && skipped < n-1; it, ok = next() {
			skipped++
		}
		for it, ok = next(); ok; it, ok = next() {
			if !yield(it) {
				return
			}
		}
	}
}
