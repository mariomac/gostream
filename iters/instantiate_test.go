package iters

import (
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterate(t *testing.T) {
	gen := Iterate(2, func(n int) int {
		return n * n
	})
	assert.Equal(t,
		[]int{2, 4, 16, 256, 65536},
		slices.Collect(Limit(5, gen)))
	// test that iterating for the second time produces the same results
	assert.Equal(t,
		[]int{2, 4, 16, 256, 65536},
		slices.Collect(Limit(5, gen)))
}

func TestGenerate(t *testing.T) {
	cnt := 0
	gen := Generate(func() int {
		cnt++
		return cnt
	})
	assert.Equal(t,
		[]int{1, 2, 3, 4, 5},
		slices.Collect(Limit(5, gen)))
}

func TestConcat(t *testing.T) {
	concat := Concat[int](
		slices.Values([]int{1, 2, 3, 4, 5, 6}),
		slices.Values([]int{7, 8, 9, 10}),
	)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, slices.Collect(concat))
	// test that iterating for the second time produces the same results
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, slices.Collect(concat))
}

func TestEmpty(t *testing.T) {
	assert.Empty(t, slices.Collect(Empty[int]()))
}

func TestSeq2KeyValues(t *testing.T) {
	// Create an iter.Seq2 from key-value pairs
	seq2 := maps.All(map[string]int{"a": 1, "b": 2, "c": 3})

	assert.Equal(t, []string{"a", "b", "c"},
		slices.Sorted(Keys(seq2)))

	assert.Equal(t, []int{1, 2, 3},
		slices.Sorted(Values(seq2)))
}
