// Package order provides helper functions and contraints to allow ordering
// streams
package order

import (
	"cmp"
	"strings"

	"golang.org/x/exp/constraints"

	"github.com/mariomac/gostream/item"
)

// Comparator function compares its two arguments for order. Returns a negative
// integer, zero, or a positive integer as the first argument is less than,
// equal to, or greater than the second.
type Comparator[T any] func(a, b T) int

// Natural implements the Comparator for those elements whose type
// has a natural order (numbers and strings).
// Deprecated in favor of cmp.Compare
func Natural[T constraints.Ordered](a, b T) int {
	return cmp.Compare(a, b)
}

// IgnoreCase implements order.Comparator for strings, ignoring the case.
func IgnoreCase(a, b string) int {
	return cmp.Compare(strings.ToLower(a), strings.ToLower(b))
}

// Inverse result of the Comparator function for inverted sorts
func Inverse[T any](cmp Comparator[T]) Comparator[T] {
	return func(a, b T) int {
		return -cmp(a, b)
	}
}

// ByKey uses the source comparator to compare the key of two item.Pair entries
func ByKey[K comparable, V any](cmp Comparator[K]) Comparator[item.Pair[K, V]] {
	return func(a, b item.Pair[K, V]) int {
		return cmp(a.Key, b.Key)
	}
}

// ByVal uses the source comparator to compare the value of two item.Pair entries
func ByVal[K, V comparable](cmp Comparator[V]) Comparator[item.Pair[K, V]] {
	return func(a, b item.Pair[K, V]) int {
		return cmp(a.Val, b.Val)
	}
}
