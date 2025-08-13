package stream

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mariomac/gostream/item"
)

func TestIsZero(t *testing.T) {
	strs := Of("hello", "", "colleague", "", "mine").
		Filter(item.Not(item.IsZero[string])).
		ToSlice()
	assert.Equal(t,
		[]string{"hello", "colleague", "mine"},
		strs)
}

func TestIsZero_Pointers(t *testing.T) {
	one := 1
	two := 2
	three := 3
	elems := Of(&one, nil, &two, nil, &three).
		Filter(item.Not(item.IsZero[*int])).
		ToSlice()
	assert.Equal(t,
		[]*int{&one, &two, &three},
		elems)
}

func TestIsZero_Structs(t *testing.T) {
	type foo struct {
		A int
		B string
	}
	elems := Of(foo{A: 1}, foo{}, foo{B: "hello"}, foo{}).
		Filter(item.Not(item.IsZero[foo])).
		ToSlice()
	assert.Equal(t,
		[]foo{{A: 1}, {B: "hello"}},
		elems)
}

func TestNeg(t *testing.T) {
	assert.Equal(t,
		[]int{-3, -2, -1, 0, 1, 2, 3},
		Of(3, 2, 1, 0, -1, -2, -3).Map(item.Neg[int]).ToSlice(),
	)
}
