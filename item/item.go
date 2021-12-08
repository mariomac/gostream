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

// Not negates the boolean result of the input condition function
func Not[T any](condition func(i T) bool) func(i T) bool {
	return func(i T) bool {
		return !condition(i)
	}
}

// IsZero returns true if the input value corresponds to the zero value of its type:
// 0 for numeric values, empty string, false, nil pointer, etc...
func IsZero[T comparable](input T) bool {
	var zero T
	return input == zero
}
