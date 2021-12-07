// Package item provides helper types and constraints for Stream operator
package item

// Pair of Key-Value made to manage maps and other key-value structures
// TODO: use "constraints.Map" when it is defined in the constraints package
type Pair[K comparable, V any] struct {
	Key K
	Val V
}
