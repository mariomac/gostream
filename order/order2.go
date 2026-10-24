// Package order provides helper functions and contraints to allow ordering
// streams
package order

import "cmp"

// Comparator2 function compares its two key-value pairs for order. Returns a negative
// integer, zero, or a positive integer as the first argument is less than,
// equal to, or greater than the second.
type Comparator2[K comparable, V any] func(k1 K, v1 V, k2 K, v2 V) int

// ByKey uses the source comparator to compare the key of two key-value entries
func ByKey[K comparable, V any](cmp Comparator[K]) Comparator2[K, V] {
	return func(k1 K, v1 V, k2 K, v2 V) int {
		return cmp(k1, k2)
	}
}

// ByVal uses the source comparator to compare the value of two key-value entries
func ByVal[K, V cmp.Ordered](cmp Comparator[V]) Comparator2[K, V] {
	return func(k1 K, v1 V, k2 K, v2 V) int {
		return cmp(v1, v2)
	}
}
