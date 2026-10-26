package stream

import (
	cmp "cmp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToSlice(t *testing.T) {
	// Testing ToSlice of concrete streams
	slice := OfSlice([]int{1, 2, 3, 4, 5, 6})
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, slice.ToSlice())

	// Testing ToSlice of connectedStreams
	assert.Equal(t, []int{1, 2, 3},
		slice.Filter(func(n int) bool {
			return n <= 3
		}).ToSlice())
}

func add(a, b int) int { return a + b }
func TestReduce(t *testing.T) {
	// test empty stream
	_, ok := Empty[int]().Reduce(add)
	assert.False(t, ok)

	// test one-element stream
	red, ok := Of(8).Reduce(add)
	assert.True(t, ok)
	assert.Equal(t, 8, red)

	// test multi-element stream
	red, ok = Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Reduce(add)
	assert.True(t, ok)
	assert.Equal(t, 55, red)
}

func notZero(s string) bool { return s != "" }
func isZero(s string) bool  { return s == "" }
func TestStream_AllMatch(t *testing.T) {
	// for empty streams, following Java behavior as reference
	assert.True(t, Empty[string]().AllMatch(isZero))
	assert.True(t, Of("hello", "world").AllMatch(notZero))
	assert.False(t, Of("", "world").AllMatch(notZero))
}

func TestStream_AnyMatch(t *testing.T) {
	// for empty streams, following Java behavior as reference
	assert.False(t, Empty[string]().AnyMatch(isZero))
	assert.True(t, Of("hello", "world").AnyMatch(notZero))
	assert.True(t, Of("", "world").AnyMatch(notZero))
	assert.False(t, Of("", "").AnyMatch(notZero))
}

func TestStream_NoneMatch(t *testing.T) {
	// for empty streams, following Java behavior as reference
	assert.True(t, Empty[string]().NoneMatch(isZero))
	assert.False(t, Of("hello", "world").NoneMatch(notZero))
	assert.False(t, Of("", "world").NoneMatch(notZero))
	assert.True(t, Of("", "").NoneMatch(notZero))
}

func testInc(a int) int { return a + 1 }

func TestCount(t *testing.T) {
	assert.Equal(t, 0, Empty[int]().Count())
	assert.Equal(t, 0, Of(1, 2, 3).Skip(3).Count())
	assert.Equal(t, 3, Of(1, 2, 3).Count())
	assert.Equal(t, 3, Of(1, 2, 3, 4, 5, 6).Skip(3).Count())
	assert.Equal(t, 8, Iterate[int](1, testInc).Limit(8).Count())
}

func TestFindFirst(t *testing.T) {
	_, ok := Empty[int]().FindFirst()
	require.False(t, ok)

	_, ok = Of(1, 2, 3).Skip(3).FindFirst()
	require.False(t, ok)

	n, ok := Of(1, 2, 3).FindFirst()
	require.True(t, ok)
	assert.Equal(t, 1, n)

	n, ok = Of(1, 2, 3, 4, 5, 6).Skip(3).FindFirst()
	require.True(t, ok)
	assert.Equal(t, 4, n)

	n, ok = Iterate[int](1, testInc).Limit(8).FindFirst()
	require.True(t, ok)
	assert.Equal(t, 1, n)
}

func TestMax(t *testing.T) {
	_, ok := Empty[int]().Max(cmp.Compare[int])
	require.False(t, ok)

	_, ok = Of(1, 2, 3).Skip(3).Max(cmp.Compare[int])
	require.False(t, ok)

	n, ok := Of(1, 2, 3).Skip(2).Max(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = Of(1, 2, 3).Max(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = Of(1, 2, 3, 4, 5, 6).Skip(3).Max(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 6, n)

	n, ok = Iterate[int](1, testInc).Limit(8).Max(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 8, n)
}

func TestMin(t *testing.T) {
	_, ok := Empty[int]().Min(cmp.Compare[int])
	require.False(t, ok)

	_, ok = Of(1, 2, 3).Skip(3).Min(cmp.Compare[int])
	require.False(t, ok)

	n, ok := Of(1, 2, 3).Skip(2).Min(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = Of(1, 2, 3).Min(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 1, n)

	n, ok = Of(1, 2, 3, 4, 5, 6).Skip(3).Min(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 4, n)

	n, ok = Iterate[int](1, testInc).Limit(8).Min(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 1, n)
}
