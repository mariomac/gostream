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
	concat := Of(1, 2, 3, 4, 5, 6).
		Concat(Of(7, 8, 9, 10))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, concat.ToSlice())
	// test that iterating for the second time produces the same results
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, concat.ToSlice())
}

func TestEmpty(t *testing.T) {
	assert.Empty(t, Empty[int]().Map(rand.Intn).ToSlice())
}

func TestOfChannel(t *testing.T) {
	elems := make(chan string)
	go func() {
		elems <- "por"
		elems <- "el"
		elems <- "puente"
		elems <- "de"
		elems <- "aranda"
		close(elems)
	}()
	assert.Equal(t,
		[]string{"por", "el", "puente", "de", "aranda"},
		OfChannel[string](elems).ToSlice())
}

func TestOfSeq(t *testing.T) {
	// Create an iter.Seq from a slice using slices.Values
	values := []int{1, 2, 3, 4, 5}
	seq := func(yield func(int) bool) {
		for _, v := range values {
			if !yield(v) {
				return
			}
		}
	}

	stream := OfSeq(seq)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, stream.ToSlice())

	// test that iterating for the second time produces the same results
	stream = OfSeq(seq)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, stream.ToSlice())
}
