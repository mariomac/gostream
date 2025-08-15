package stream

import (
	"cmp"
	"maps"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mariomac/gostream/item"
	"github.com/mariomac/gostream/order"
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

func TestOfMap_SortedByKey(t *testing.T) {
	months := OfMap(map[int]string{
		1: "Jan", 2: "Feb", 3: "Mar", 4: "Apr", 5: "May", 6: "Jun",
		7: "Jul", 8: "Aug", 9: "Sep", 10: "Oct", 11: "Nov", 12: "Dec",
	}).Sorted(order.ByKey[int, string](cmp.Compare[int]))

	monthNames := Map(months, func(p item.Pair[int, string]) string {
		return p.Val
	}).ToSlice()

	assert.Equal(t,
		[]string{
			"Jan", "Feb", "Mar", "Apr", "May", "Jun",
			"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
		},
		monthNames,
	)
}

func TestOfMap_SortedByVal(t *testing.T) {
	months := OfMap(map[int]string{
		1: "Jan", 2: "Feb", 3: "Mar", 4: "Apr", 5: "May", 6: "Jun",
		7: "Jul", 8: "Aug", 9: "Sep", 10: "Oct", 11: "Nov", 12: "Dec",
	}).Sorted(order.ByVal[int, string](cmp.Compare[string]))

	monthNames := Map(months, func(p item.Pair[int, string]) string {
		return p.Val
	}).ToSlice()

	assert.Equal(t,
		[]string{
			"Apr", "Aug", "Dec", "Feb", "Jan", "Jul",
			"Jun", "Mar", "May", "Nov", "Oct", "Sep",
		},
		monthNames,
	)
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

func TestOfSeq2(t *testing.T) {
	// Create an iter.Seq2 from key-value pairs
	seq2 := maps.All(map[string]int{"a": 1, "b": 2, "c": 3})

	stream := OfSeq2(seq2)
	assert.Equal(t, 3, stream.Count())
	assert.True(t, stream.AnyMatch(item.Equals(
		item.Pair[string, int]{Key: "a", Val: 1})))
	assert.True(t, stream.AnyMatch(item.Equals(
		item.Pair[string, int]{Key: "b", Val: 2})))
	assert.True(t, stream.AnyMatch(item.Equals(
		item.Pair[string, int]{Key: "c", Val: 3})))

	// test that iterating for the second time produces the same results
	stream = OfSeq2(seq2)
	assert.Equal(t, 3, stream.Count())
	assert.True(t, stream.AnyMatch(item.Equals(
		item.Pair[string, int]{Key: "a", Val: 1})))
	assert.True(t, stream.AnyMatch(item.Equals(
		item.Pair[string, int]{Key: "b", Val: 2})))
	assert.True(t, stream.AnyMatch(item.Equals(
		item.Pair[string, int]{Key: "c", Val: 3})))
}
