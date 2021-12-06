package stream

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterate(t *testing.T) {
	gen := Iterate(2, func(n int) int {
		return n * n
	})
	assert.Equal(t,
		[]int{2, 4, 16, 256, 65536},
		gen.Limit(5).ToSlice())
	// test that iterating for the second time produces the same results
	assert.Equal(t,
		[]int{2, 4, 16, 256, 65536},
		gen.Limit(5).ToSlice())
}

func TestGenerate(t *testing.T) {
	cnt := 0
	gen := Generate(func() int {
		cnt++
		return cnt
	})
	assert.Equal(t,
		[]int{1, 2, 3, 4, 5},
		gen.Limit(5).ToSlice())
}

func TestConcat(t *testing.T) {
	concat := Concat[int](
		Of(1, 2, 3, 4, 5, 6),
		Of(7, 8, 9, 10),
	)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, concat.ToSlice())
	// test that iterating for the second time produces the same results
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, concat.ToSlice())
}

func TestEmpty(t *testing.T) {
	assert.Empty(t, Empty[int]().Map(rand.Intn).ToSlice())
}
