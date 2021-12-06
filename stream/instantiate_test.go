package stream

import (
	"github.com/mariomac/gostream/kv"
	"github.com/mariomac/gostream/order"
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

func TestOfMap_SortedByKey(t *testing.T) {
	months := OfMap(map[int]string{
		1: "Jan", 2: "Feb", 3: "Mar", 4: "Apr", 5: "May", 6: "Jun",
		7: "Jul", 8: "Aug", 9: "Sep", 10: "Oct", 11: "Nov", 12: "Dec",
	}).Sorted(order.ByKey[int, string](order.Natural[int]))

	monthNames := Map(months, func(p kv.Pair[int, string]) string {
		return p.Val
	}).ToSlice()

	assert.Equal(t,
		[]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
			"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		monthNames,
	)
}

func TestOfMap_SortedByVal(t *testing.T) {
	months := OfMap(map[int]string{
		1: "Jan", 2: "Feb", 3: "Mar", 4: "Apr", 5: "May", 6: "Jun",
		7: "Jul", 8: "Aug", 9: "Sep", 10: "Oct", 11: "Nov", 12: "Dec",
	}).Sorted(order.ByVal[int, string](order.Natural[string]))

	monthNames := Map(months, func(p kv.Pair[int, string]) string {
		return p.Val
	}).ToSlice()

	assert.Equal(t,
		[]string{"Apr", "Aug", "Dec", "Feb", "Jan", "Jul",
			"Jun", "Mar", "May", "Nov", "Oct", "Sep"},
		monthNames,
	)
}
