// Package item provides helper types and constraints for Stream operator
package item

import "constraints"

// Number constraint: ints, uints, complexes, floats and all their subtypes
type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

// Number constraint any type that define the addition + operation, and their subtypes
type Addable interface {
	constraints.Ordered | constraints.Complex
}

// Pair of Key-Value made to manage maps and other key-value structures
// TODO: use "constraints.Map" when it is defined in the constraints package
type Pair[K comparable, V any] struct {
	Key K
	Val V
}

// Add the two arguments using the plus + operator
func Add[T Addable](a, b T) T {
	return a + b
}

// Multiply the two arguments using the multiplication * operator
func Multiply[T Number](a, b T) T {
	return a * b
}

// Increment the argument
func Increment[T Number](a T) T {
	return a + 1
}
