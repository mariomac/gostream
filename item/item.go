// Package item provides helper types and constraints for Stream operator
package item

import (
	"cmp"
)

// Number constraint: ints, uints, complexes, floats and all their subtypes
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~complex64 | ~complex128
}

// Addable constraint any type that define the addition + operation, and their subtypes
type Addable interface {
	cmp.Ordered | ~complex64 | ~complex128
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
func GreaterThan[T cmp.Ordered](reference T) func(i T) bool {
	return func(i T) bool {
		return i > reference
	}
}

// GreaterThanOrEq returns a predicate that is true when the checked value is equal or larger than
// the provided reference.
func GreaterThanOrEq[T cmp.Ordered](reference T) func(i T) bool {
	return func(i T) bool {
		return i >= reference
	}
}

// LessThan returns a predicate that is true when the checked value is less than
// the provided reference.
func LessThan[T cmp.Ordered](reference T) func(i T) bool {
	return func(i T) bool {
		return i < reference
	}
}

// LessThanOrEq returns a predicate that is true when the checked value is equal or less than
// the provided reference.
func LessThanOrEq[T cmp.Ordered](reference T) func(i T) bool {
	return func(i T) bool {
		return i >= reference
	}
}

// IsCmp is an auxiliary function that has to be passed as argument to the stream.Stream's Distinct method.
func IsCmp[T comparable](i T) T {
	return i
}
