// Package order provides helper functions and contraints to allow ordering
// streams
package order

import (
	"github.com/mariomac/gostream/kv"
	"sort"
	"strings"
)

// Comparator function return true if the left argument precedes
// the right argument.
// Some bundled implementors are order.Integer, order.Natural, order.IgnoreCase
type Comparator[T any] func(a, b T) bool

// those bundled types that have a natural order
type natural interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Natural implements the Comparator for those elements whose type
// has a natural order (numbers and strings)

func Natural[T natural](a, b T) bool {
	return a < b
}

// IgnoreCase implements order.Comparator for strings, ignoring the case.
func IgnoreCase(a, b string) bool {
	return strings.ToLower(a) < strings.ToLower(b)
}

// Inverse result of the Comparator function for inverted sorts
func Inverse[T any](less Comparator[T]) Comparator[T] {
	return func(a, b T) bool {
		return !less(a, b)
	}
}

// ByKey uses the source comparator to compare the key of two kv.Pair entries
func ByKey[K comparable, V any](less Comparator[K]) Comparator[kv.Pair[K, V]] {
	return func(a, b kv.Pair[K, V]) bool {
		return less(a.Key, b.Key)
	}
}

// ByVal uses the source comparator to compare the value of two kv.Pair entries
func ByVal[K, V comparable](less Comparator[V]) Comparator[kv.Pair[K, V]] {
	return func(a, b kv.Pair[K, V]) bool {
		return less(a.Val, b.Val)
	}
}

// SortSlice sorts the given slice according to the criteria in the provided comparator
func SortSlice[T any](slice []T, comparator Comparator[T]) {
	sort.Sort(&sorter[T]{items: slice, comparator: comparator})
}

type sorter[T any] struct {
	items      []T
	comparator Comparator[T]
}

func (s *sorter[T]) Len() int {
	return len(s.items)
}

func (s *sorter[T]) Less(i, j int) bool {
	return s.comparator(s.items[i], s.items[j])
}

func (s *sorter[T]) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}
