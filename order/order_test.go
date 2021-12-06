package order

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNaturalOrder(t *testing.T) {
	assert.True(t, Natural(1, 2))
	assert.False(t, Natural(2, 1))
	assert.False(t, Natural(-3, -3))

	assert.True(t, Natural(1.0, 1.1))
	assert.False(t, Natural(1.1, 1.0))
	assert.False(t, Natural(-3.0, -3.0))

	assert.True(t, Natural("Hello", "amici"))
	assert.False(t, Natural("aaa", "Bbb"))
	assert.False(t, Natural("lololo", "lololo"))
}

func TestIgnoreCase(t *testing.T) {
	assert.True(t, IgnoreCase("aaa", "Bbb"))
	assert.False(t, IgnoreCase("Hello", "amici"))
	assert.False(t, IgnoreCase("LoloLo", "lOlolo"))
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
