package stream

import (
	"cmp"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	in := Of[int](1, 2, 3, 4, 5)
	double := in.Map(func(i int) int {
		return i * 2
	})
	assert.Equal(t, []int{2, 4, 6, 8, 10}, double.ToSlice())

	strings := in.Map(strconv.Itoa)
	assert.Equal(t, []string{"1", "2", "3", "4", "5"}, strings.ToSlice())
}

func TestForEach(t *testing.T) {
	var copy []int
	Of[int](1, 2, 3, 4, 5).ForEach(func(i int) {
		copy = append(copy, i)
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, copy)
}

func TestFilter(t *testing.T) {
	in := Of[int](1, 2, 3, 4, 5)
	odds := in.Filter(func(n int) bool {
		return n%2 == 1
	})
	assert.Equal(t, []int{1, 3, 5}, odds.ToSlice())
	empty := in.Filter(func(_ int) bool {
		return false
	})
	assert.Empty(t, empty.ToSlice())
	all := in.Filter(func(_ int) bool {
		return true
	})
	assert.Equal(t, []int{1, 2, 3, 4, 5}, all.ToSlice())
}

func TestLimit(t *testing.T) {
	count := 0
	items := Generate(func() int {
		count++
		return count
	}).Limit(7).ToSlice()

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, items)
}

func TestDistinct(t *testing.T) {
	assert.Equal(t,
		[]int{1, 2, 3, 4, 5},
		Of(1, 1, 2, 3, 3, 3, 4, 5, 1, 2, 3, 4, 5).Distinct().ToSlice(),
	)
}

func TestDistinctNonComparable(t *testing.T) {
	items := [][]int{{1}, {1}, {2}}
	assert.Equal(t, items, OfSlice(items).Distinct().ToSlice())
}

func TestDistinctInterfaceValues(t *testing.T) {
	items := []any{nil, nil, 1, 1, "1", []int{2}, []int{2}, map[string]int{"a": 1}}
	assert.Equal(t,
		[]any{nil, 1, "1", []int{2}, []int{2}, map[string]int{"a": 1}},
		OfSlice(items).Distinct().ToSlice())
}

func TestDistinctNestedInterfaceValues(t *testing.T) {
	type entry struct{ Value any }
	items := []entry{{1}, {1}, {[]int{2}}, {[]int{2}}}
	assert.Equal(t, []entry{{1}, {[]int{2}}, {[]int{2}}}, OfSlice(items).Distinct().ToSlice())
	arrays := [][1]any{{1}, {1}, {[]int{2}}, {[]int{2}}}
	assert.Equal(t, [][1]any{{1}, {[]int{2}}, {[]int{2}}}, OfSlice(arrays).Distinct().ToSlice())
}

func TestDistinctLazyAndReusable(t *testing.T) {
	consumed := 0
	distinct := Of(1, 1, 2, 2, 3).Peek(func(int) { consumed++ }).Distinct()
	assert.Zero(t, consumed)
	assert.Equal(t, []int{1, 2}, distinct.Limit(2).ToSlice())
	assert.Equal(t, 3, consumed)
	assert.Equal(t, []int{1, 2, 3}, distinct.ToSlice())
	assert.Empty(t, Empty[int]().Distinct().ToSlice())
}

func TestDistinctInfinite(t *testing.T) {
	distinct := Iterate(0, func(n int) int { return n + 1 }).Distinct()
	assert.Panics(t, func() { distinct.ToSlice() })
	assert.Equal(t, []int{0, 1, 2}, distinct.Limit(3).ToSlice())
	assert.Panics(t, func() { Generate(func() []int { return []int{1} }).Distinct().ToSlice() })
}

func TestSort(t *testing.T) {
	assert.Equal(t,
		[]int{1, 1, 2, 3, 5, 6, 7, 8, 8},
		Of(1, 7, 8, 3, 2, 1, 5, 8, 6).
			Sorted(cmp.Compare[int]).ToSlice())
}

func TestFlapMap(t *testing.T) {
	generateCharSequence := func(in string) *Stream[byte] {
		return OfSlice([]byte(in))
	}
	generateNillableCharSequence := func(in string) *Stream[byte] {
		if len(in) == 0 {
			return nil
		}
		return OfSlice([]byte(in))
	}

	assert.Empty(t, Of[string]().FlatMap(generateCharSequence).ToSlice())

	chars := Of[string]("", "hello", "my", "", "friends!", "").
		FlatMap(generateCharSequence).ToSlice()
	assert.Equal(t, []byte("hellomyfriends!"), chars)

	chars = Of[string]("", "hello", "my", "", "friends!", "").
		FlatMap(generateNillableCharSequence).ToSlice()
	assert.Equal(t, []byte("hellomyfriends!"), chars)
}

func TestPeek(t *testing.T) {
	var actions []string

	sl := Of(1, 2, 3, 4, 5, 6).
		Filter(func(n int) bool {
			return n%2 == 1
		}).
		Peek(func(n int) {
			actions = append(actions, fmt.Sprint("processed ", n))
		}).
		ToSlice()

	require.Equal(t, []int{1, 3, 5}, sl)
	assert.Equal(t, []string{
		"processed 1", "processed 3", "processed 5",
	}, actions)
}

func TestSkip(t *testing.T) {
	assert.Empty(t, Empty[int]().Skip(3).ToSlice())
	assert.Empty(t, Of(1, 2).Skip(3).ToSlice())
	assert.Empty(t, Of(1, 2, 3).Skip(3).ToSlice())
	assert.Equal(t, []int{4, 5, 6}, Of(1, 2, 3, 4, 5, 6).Skip(3).ToSlice())
}
