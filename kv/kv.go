// Package fn provides helper tools for functional operation of key-value pairs
package kv

// Pair of Key-Value made to manage maps and other key-value structures
type Pair[K comparable, V any] struct {
	Key K
	Val V
}
