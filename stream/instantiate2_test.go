package stream

import (
	"cmp"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mariomac/gostream/item"
	"github.com/mariomac/gostream/order"
)

func TestIterate2(t *testing.T) {
	gen := Iterate2(2, 4, func(k, v int) (int, int) {
		return k * k, v + 1
	})
	assert.Equal(t, map[int]int{
		2:     4,
		4:     5,
		16:    6,
		256:   7,
		65536: 8,
	}, gen.Limit(5).ToMap())
	// test that iterating for the second time produces the same results
	assert.Equal(t, map[int]int{
		2:     4,
		4:     5,
		16:    6,
		256:   7,
		65536: 8,
	}, gen.Limit(5).ToMap())
}

func TestGenerate2(t *testing.T) {
	cnt, cnt2 := 0, "a"
	gen := Generate2(func() (int, string) {
		cnt++
		cnt2 += "a"
		return cnt, cnt2
	})
	assert.Equal(t,
		map[int]string{
			1: "a",
			2: "aa",
			3: "aaa",
		},
		gen.Limit(3).ToMap())
}

func TestConcat2(t *testing.T) {
	concat := OfMap(map[string]int{"a": 1, "b": 2, "c": 3}).
		Concat(OfMap(map[string]int{"d": 1, "e": 2, "f": 3}))
	assert.Equal(t, map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 1,
		"e": 1,
		"f": 1,
	}, concat.ToMap())
	// test that iterating for the second time produces the same results
	assert.Equal(t, map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 1,
		"e": 1,
		"f": 1,
	}, concat.ToMap())
}

func TestEmpty2(t *testing.T) {
	assert.Empty(t, Empty2[int, int]().Map(func(int, int) (float64, float64) {
		return 1, 2
	}).ToMap())
}

func TestOfMap(t *testing.T) {
	months := OfMap(map[int]string{
		1: "Jan", 2: "Feb", 3: "Mar", 4: "Apr", 5: "May", 6: "Jun",
		7: "Jul", 8: "Aug", 9: "Sep", 10: "Oct", 11: "Nov", 12: "Dec",
	})



	monthNames := months.Map(func(p item.Pair[int, string]) string {
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

	monthNames := months.Map(func(p item.Pair[int, string]) string {
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
