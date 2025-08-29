// Package item provides helper types and constraints for Stream operator
package item

import "golang.org/x/exp/constraints"

// Number constraint: ints, uints, complexes, floats and all their subtypes
type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

// Addable constraint any type that define the addition + and subtraction - operation, and their subtypes
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

// Neg inverts the sign of the given numeric argument
func Neg[T Number](a T) T {
	return -a
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

// Equals returns a predicate that is true when the checked value is
// equal to the provided reference.
func Equals[T comparable](reference T) func(i T) bool {
	return func(i T) bool {
		return i == reference
	}
}

// GreaterThan returns a predicate that is true when the checked value is larger than
// the provided reference.
func GreaterThan[T constraints.Ordered](reference T) func(i T) bool {
	return func(i T) bool {
		return i > reference
	}
}

// GreaterThanOrEq returns a predicate that is true when the checked value is equal or larger than
// the provided reference.
func GreaterThanOrEq[T constraints.Ordered](reference T) func(i T) bool {
	return func(i T) bool {
		return i >= reference
	}
}

// LessThan returns a predicate that is true when the checked value is less than
// the provided reference.
func LessThan[T constraints.Ordered](reference T) func(i T) bool {
	return func(i T) bool {
		return i < reference
	}
}

// LessThanOrEq returns a predicate that is true when the checked value is equal or less than
// the provided reference.
func LessThanOrEq[T constraints.Ordered](reference T) func(i T) bool {
	return func(i T) bool {
		return i >= reference
	}
}
