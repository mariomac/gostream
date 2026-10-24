package stream

import (
	"slices"

	"github.com/mariomac/gostream/order"
)

// Map returns a Stream2 consisting of the results of individually applying
// the mapper function to each key-value pair of the target Stream2.
func (is *Stream2[K, V]) Map[OK comparable, OV comparable](mapper func(K, V) (OK, OV)) *Stream2[OK, OV] {
	return &Stream2[OK, OV]{
		supply: func() iterator2[OK, OV] {
			next := is.iterator()
			return func() (OK, OV, bool) {
				k, v, ok := next()
				if !ok {
					return *new(OK), *new(OV), false
				}
				okk, ov := mapper(k, v)
				return okk, ov, true
			}
		},
	}
}

// MapStream returns a Stream (single output type) consisting of the results of individually applying
// the mapper function to each key-value pair of the target Stream2.
func (is *Stream2[K, V]) MapStream[O any](mapper func(K, V) O) *Stream[O] {
	return &Stream[O]{
		supply: func() iterator[O] {
			next := is.iterator()
			return func() (O, bool) {
				k, v, ok := next()
				if !ok {
					return *new(O), false
				}
				o := mapper(k, v)
				return o, true
			}
		},
	}
}

// Filter returns a Stream consisting of the items of the target stream that match the given
// predicate (this is, applying the predicate function over the item returns true).
func (is *Stream2[K, V]) Filter(predicate func(K, V) bool) *Stream2[K, V] {
	return &Stream2[K, V]{
		supply: func() iterator2[K, V] {
			next := is.iterator()
			return func() (K, V, bool) {
				for {
					k, v, ok := next()
					if !ok {
						return *new(K), *new(V), false
					}
					if predicate(k, v) {
						return k, v, true
					}
				}
			}
		},
	}
}

// Limit returns a stream consisting of the elements of this stream, truncated to
// be no longer than maxSize in length.
func (is *Stream2[K, V]) Limit(maxSize int) *Stream2[K, V] {
	return &Stream2[K, V]{
		supply: func() iterator2[K, V] {
			next := is.iterator()
			count := 0
			return func() (K, V, bool) {
				if count == maxSize {
					return *new(K), *new(V), false
				}
				k, v, ok := next()
				if !ok {
					count = maxSize
					return k, v, false
				}
				count++
				return k, v, true
			}
		},
	}
}

// FlatMap returns a stream consisting of the results of replacing each element of this stream
// with the contents of a mapped stream produced by applying the provided mapping function to
// each element. Each mapped stream is closed after its contents have been placed into this
// stream. (If a mapped stream is null an empty stream is used, instead.)
//
// Due to the lazy nature of streams, if any of the mapped streams is infinite it will remain
// unnoticed and some operations (Count, Reduce, Sorted, AllMatch...) will not end.
func (is *Stream2[KIN, VIN]) FlatMap[KOUT comparable, VOUT any](mapper func(KIN, VIN) *Stream2[KOUT, VOUT]) *Stream2[KOUT, VOUT] {
	return &Stream2[KOUT, VOUT]{
		supply: func() iterator2[KOUT, VOUT] {
			nextFromInputStream := is.iterator()
			var nextFromOutputStream iterator2[KOUT, VOUT]
			return func() (KOUT, VOUT, bool) {
				for {
					for nextFromOutputStream == nil {
						// apply the mapper to the current input item and generate an output stream
						// to iterate
						nextInputKey, nextInputVal, ok := nextFromInputStream()
						if !ok {
							return finishedIterator2[KOUT, VOUT]()
						}
						if outStream := mapper(nextInputKey, nextInputVal); outStream != nil {
							nextFromOutputStream = outStream.iterator()
						}
					}
					if outKey, outVal, ok := nextFromOutputStream(); ok {
						return outKey, outVal, true
					} else {
						// item's resulting outputStream has been iterated. Look for next input item
						nextFromOutputStream = nil
					}
				}
			}
		},
	}
}

// Peek returns a stream consisting of the elements of this stream, additionally performing
// the provided action on each element as elements are consumed from the resulting stream.
func (is *Stream2[K, V]) Peek(consumer func(K, V)) *Stream2[K, V] {
	return &Stream2[K, V]{
		supply: func() iterator2[K, V] {
			next := is.iterator()
			return func() (K, V, bool) {
				k, v, ok := next()
				if !ok {
					return k, v, false
				}
				consumer(k, v)
				return k, v, true
			}
		},
	}
}

// Concat returns a lazily concatenated stream whose elements are all the elements of the
// target Stream followed by all the elements of the argument Stream
// todo: consider multiple arguments
func (is *Stream2[K, V]) Concat(append *Stream2[K, V]) *Stream2[K, V] {
	return &Stream2[K, V]{
		supply: func() iterator2[K, V] {
			first := true
			next := is.iterator()
			return func() (K, V, bool) {
				k, v, ok := next()
				if ok {
					return k, v, true
				}
				if first {
					first = false
					next = append.iterator()
				} else {
					next = finishedIterator2[K, V]
				}
				return next()
			}
		},
	}
}

// TODO: document and test
func (is *Stream2[K, V]) Keys() *Stream[K] {
	return &Stream[K]{
		supply: func() iterator[K] {
			next := is.iterator()
			return func() (K, bool) {
				k, _, ok := next()
				return k, ok
			}
		},
	}
}

// TODO: document and test
func (is *Stream2[K, V]) Values() *Stream[V] {
	return &Stream[V]{
		supply: func() iterator[V] {
			next := is.iterator()
			return func() (V, bool) {
				_, v, ok := next()
				return v, ok
			}
		},
	}
}

type tuple[K comparable, V any] struct {
	k K
	v V
}

// TODO: poner ejemplo de orrdenar natural by value
func (is *Stream2[K, V]) Sorted(cmp order.Comparator2[K, V]) *Stream2[K, V] {
	return &Stream2[K, V]{
		supply: func() iterator2[K, V] {
			items := is.MapStream(func(k K, v V) tuple[K, V] {
				return tuple[K, V]{k: k, v: v}
			}).ToSlice()
			slices.SortFunc(items, func(a, b tuple[K, V]) int {
				return cmp(a.k, a.v, b.k, b.v)
			})
			return func() (K, V, bool) {
				if len(items) == 0 {
					return finishedIterator2[K, V]()
				}
				n := items[0]
				items = items[1:]
				return n.k, n.v, true
			}
		},
	}
}

// Skip returns a stream consisting of the remaining elements of this stream after discarding
// the first n elements of the stream.
func (is *Stream2[K, V]) Skip(n int) *Stream2[K, V] {
	return &Stream2[K, V]{
		supply: func() iterator2[K, V] {
			next := is.iterator()
			skipped := 0
			return func() (K, V, bool) {
				var k K
				var v V
				var ok bool
				for k, v, ok = next(); ok && skipped < n; k, v, ok = next() {
					skipped++
				}
				return k, v, ok
			}
		},
	}
}
