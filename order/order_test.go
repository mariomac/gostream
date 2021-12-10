package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNaturalOrder(t *testing.T) {
	assert.Negative(t, Natural(1, 2))
	assert.Positive(t, Natural(2, 1))
	assert.Zero(t, Natural(-3, -3))

	assert.Negative(t, Natural(1.0, 1.1))
	assert.Positive(t, Natural(1.1, 1.0))
	assert.Zero(t, Natural(-3.0, -3.0))

	assert.Negative(t, Natural("Hello", "amici"))
	assert.Positive(t, Natural("aaa", "Bbb"))
	assert.Zero(t, Natural("lololo", "lololo"))
}

func TestIgnoreCase(t *testing.T) {
	assert.Negative(t, IgnoreCase("aaa", "Bbb"))
	assert.Positive(t, IgnoreCase("Hello", "amici"))
	assert.Zero(t, IgnoreCase("LoloLo", "lOlolo"))
}

func TestSortSlice(t *testing.T) {
	slice := []int{1, 7, 8, 3, 2, 1, 5, 8, 6}
	SortSlice(slice, Natural[int])
	assert.Equal(t,
		[]int{1, 1, 2, 3, 5, 6, 7, 8, 8},
		slice)

	SortSlice(slice, Inverse(Natural[int]))
	assert.Equal(t,
		[]int{8, 8, 7, 6, 5, 3, 2, 1, 1},
		slice)
}
