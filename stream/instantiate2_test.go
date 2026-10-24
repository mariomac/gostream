package stream

import (
	"cmp"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"

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
	}).SortedByKey().Values().ToSlice()
	assert.Equal(t,
		[]string{
			"Jan", "Feb", "Mar", "Apr", "May", "Jun",
			"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
		},
		months,
	)
}

func TestOfMap_SortedByVal(t *testing.T) {
	months := OfMap(map[int]string{
		1: "Jan", 2: "Feb", 3: "Mar", 4: "Apr", 5: "May", 6: "Jun",
		7: "Jul", 8: "Aug", 9: "Sep", 10: "Oct", 11: "Nov", 12: "Dec",
	}).Sorted(order.ByVal[int, string](cmp.Compare[string])).
		Values()
	assert.Equal(t,
		[]string{
			"Apr", "Aug", "Dec", "Feb", "Jan", "Jul",
			"Jun", "Mar", "May", "Nov", "Oct", "Sep",
		},
		months,
	)
}

func TestOfSeq2(t *testing.T) {
	// Create an iter.Seq2 from key-value pairs
	seq2 := maps.All(map[string]int{"a": 1, "b": 2, "c": 3})
	stream := OfSeq2(seq2)
	assert.Equal(t, 3, stream.Count())
	assert.True(t, stream.AnyMatch(func(s string, i int) bool {
		return s == "a" && i == 1
	}))
	assert.True(t, stream.AnyMatch(func(s string, i int) bool {
		return s == "b" && i == 2
	}))
	assert.True(t, stream.AnyMatch(func(s string, i int) bool {
		return s == "c" && i == 3
	}))

	// test that iterating for the second time produces the same results
	assert.Equal(t, 3, stream.Count())
	assert.True(t, stream.AnyMatch(func(s string, i int) bool {
		return s == "a" && i == 1
	}))
	assert.True(t, stream.AnyMatch(func(s string, i int) bool {
		return s == "b" && i == 2
	}))
	assert.True(t, stream.AnyMatch(func(s string, i int) bool {
		return s == "c" && i == 3
	}))
}
