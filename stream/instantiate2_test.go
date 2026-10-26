package stream

import (
	"cmp"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
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
	cnt, cnt2 := 0, ""
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

func TestEmpty2(t *testing.T) {
	assert.Empty(t, Empty2[int, int]().Map(func(int, int) (float64, float64) {
		return 1, 2
	}).ToMap())
}

func TestOfMap(t *testing.T) {
	months := OfMap(map[int]string{
		1: "Jan", 2: "Feb", 3: "Mar", 4: "Apr", 5: "May", 6: "Jun",
		7: "Jul", 8: "Aug", 9: "Sep", 10: "Oct", 11: "Nov", 12: "Dec",
	}).SortedByKey(cmp.Compare).Values().ToSlice()
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
	}).SortedByVal(cmp.Compare).Values().ToSlice()
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

func TestConcat2(t *testing.T) {
	dupe := func(i int) (int, int) { return i, i }
	dedupe := func(i, _ int) int { return i }
	type testCase struct {
		title string
		input []*Stream2[int, int]
	}
	for _, tc := range []testCase{{
		title: "single",
		input: []*Stream2[int, int]{Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).MapStream2(dupe)},
	}, {
		title: "basic double",
		input: []*Stream2[int, int]{Of(1, 2, 3, 4, 5, 6).MapStream2(dupe), Of(7, 8, 9, 10).MapStream2(dupe)},
	}, {
		title: "multiple",
		input: []*Stream2[int, int]{Of(1, 2, 3).MapStream2(dupe), Of(4, 5, 6).MapStream2(dupe), Of(7, 8).MapStream2(dupe), Of(9, 10).MapStream2(dupe)},
	}, {
		title: "multiple with empties",
		input: []*Stream2[int, int]{Of(1, 2, 3, 4, 5, 6).MapStream2(dupe), Empty2[int, int](), Empty2[int, int](), Of(7, 8).MapStream2(dupe), Of(9, 10).MapStream2(dupe)},
	}, {
		title: "multiple with empties at start and end",
		input: []*Stream2[int, int]{Empty2[int, int](), Of(1, 2, 3).MapStream2(dupe), Of(4, 5, 6).MapStream2(dupe), Of(7, 8).MapStream2(dupe), Of(9, 10).MapStream2(dupe), Empty2[int, int]()},
	}} {
		t.Run(tc.title, func(t *testing.T) {
			concat := Concat2(tc.input...)
			assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, concat.MapStream(dedupe).ToSlice())
			// test that iterating for the second time produces the same results
			assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, concat.MapStream(dedupe).ToSlice())
		})
	}
}

func TestConcat2_Empty(t *testing.T) {
	type testCase struct {
		title string
		input []*Stream2[int, int]
	}
	for _, tc := range []testCase{
		{title: "no arguments"},
		{title: "empty", input: []*Stream2[int, int]{Empty2[int, int]()}},
		{title: "multiple empty", input: []*Stream2[int, int]{Empty2[int, int](), Empty2[int, int](), Empty2[int, int]()}},
	} {
		t.Run(tc.title, func(t *testing.T) {
			concat := Concat2(tc.input...)
			assert.Empty(t, concat.Count())
		})
	}
}

func TestZip(t *testing.T) {
	assert.Equal(t,
		map[int]string{1: "one", 2: "two", 3: "three"},
		Zip(
			Of(1, 2, 3, 4),
			Of("one", "two", "three"),
		).ToMap())
	assert.Equal(t,
		map[string]int{"one": 1, "two": 2, "three": 3},
		Zip(
			Of("one", "two", "three"),
			Of(1, 2, 3, 4),
		).ToMap())

	assert.Empty(t,
		Zip(
			Empty[int](),
			Of(1, 2, 3, 4),
		).ToMap())
	assert.Empty(t,
		Zip(
			Of(1, 2, 3, 4),
			Empty[int](),
		).ToMap())
}
