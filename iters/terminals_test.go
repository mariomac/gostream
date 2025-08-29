package iters

import (
	"cmp"
	"slices"
	"testing"

	"github.com/mariomac/gostream/item"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestReduce(t *testing.T) {
	// test empty iter.Seq
	_, ok := Reduce(slices.Values([]int{}), item.Add[int])
	assert.False(t, ok)

	// test one-element iter.Seq
	red, ok := Reduce(slices.Values([]int{8}), item.Add[int])
	assert.True(t, ok)
	assert.Equal(t, 8, red)

	// test multi-element iter.Seq
	red, ok = Reduce(slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}), item.Add[int])
	assert.True(t, ok)
	assert.Equal(t, 55, red)
}

func TestIterableStream_AllMatch(t *testing.T) {
	// for empty iter.Seq, following Java behavior as reference
	assert.True(t, AllMatch(slices.Values([]string{}), item.IsZero[string]))
	assert.True(t, AllMatch(slices.Values([]string{"hello", "world"}), item.Not(item.IsZero[string])))
	assert.False(t, AllMatch(slices.Values([]string{"", "world"}), item.Not(item.IsZero[string])))
}

func TestIterableStream_AnyMatch(t *testing.T) {
	// for empty iter.Seq, following Java behavior as reference
	assert.False(t, AnyMatch(slices.Values([]string{}), item.IsZero[string]))
	assert.True(t, AnyMatch(slices.Values([]string{"hello", "world"}),item.Not(item.IsZero[string])))
	assert.True(t, AnyMatch(slices.Values([]string{"", "world"}),item.Not(item.IsZero[string])))
	assert.False(t, AnyMatch(slices.Values([]string{"", ""}),item.Not(item.IsZero[string])))
}

func TestIterableStream_NoneMatch(t *testing.T) {
	// for empty iter.Seq, following Java behavior as reference
	assert.True(t, NoneMatch(slices.Values([]string{}),item.IsZero[string]))
	assert.False(t, NoneMatch(slices.Values([]string{"hello", "world"}),item.Not(item.IsZero[string])))
	assert.False(t, NoneMatch(slices.Values([]string{"", "world"}),item.Not(item.IsZero[string])))
	assert.True(t, NoneMatch(slices.Values([]string{"", ""}),item.Not(item.IsZero[string])))
}

func TestCount(t *testing.T) {
	assert.Equal(t, 0, Count(slices.Values([]int{}),))
	assert.Equal(t, 0, Count(Skip(slices.Values([]int{1,2,3}), 3)))
	assert.Equal(t, 3, Count(slices.Values([]int{1,2,3})))
	assert.Equal(t, 3, Count(Skip(slices.Values([]int{1,2,3,4,5,6}), 3)))
	assert.Equal(t, 8, Count(Limit(Iterate[int](1, item.Increment[int]), 8)))
}

func TestFindFirst(t *testing.T) {
	_, ok := FindFirst(slices.Values([]int{}))
	require.False(t, ok)

	_, ok = FindFirst(Skip(slices.Values([]int{1,2, 3}), 3))
	require.False(t, ok)

	n, ok := FindFirst(slices.Values([]int{1, 2, 3}))
	require.True(t, ok)
	assert.Equal(t, 1, n)

	n, ok = FindFirst(Skip(slices.Values([]int{1, 2, 3, 4, 5, 6}), 3))
	require.True(t, ok)
	assert.Equal(t, 4, n)

	n, ok = FindFirst(Limit(Iterate[int](1, item.Increment[int]), 8))
	require.True(t, ok)
	assert.Equal(t, 1, n)
}

func TestMax(t *testing.T) {
	_, ok := Max(slices.Values([]int{}))
	require.False(t, ok)

	_, ok = Max(Skip(slices.Values([]int{1, 2, 3}), 3))
	require.False(t, ok)

	n, ok := Max(Skip(slices.Values([]int{1, 2, 3}), 2))
	require.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = Max(slices.Values([]int{1, 3, 2}))
	require.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = Max(Skip(slices.Values([]int{1, 2, 3, 4, 5, 6}), 3),)
	require.True(t, ok)
	assert.Equal(t, 6, n)
}

func TestMin(t *testing.T) {
	_, ok := Min(slices.Values([]int{}),)
	require.False(t, ok)

	n, ok := Min(slices.Values([]int{1, 2, 3}),)
	require.True(t, ok)
	assert.Equal(t, 1, n)
}

func TestMaxFunc(t *testing.T) {
	_, ok := MaxFunc(slices.Values([]int{}), cmp.Compare[int])
	require.False(t, ok)

	_, ok = MaxFunc(Skip(slices.Values([]int{1, 2, 3}), 3), cmp.Compare[int])
	require.False(t, ok)

	n, ok := MaxFunc(Skip(slices.Values([]int{1, 2, 3}), 2), cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = MaxFunc(slices.Values([]int{1, 3, 2}), cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = MaxFunc(Skip(slices.Values([]int{1, 2, 3, 4, 5, 6}), 3), cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 6, n)
}

func TestMinFunc(t *testing.T) {
	_, ok := MinFunc(slices.Values([]int{}), cmp.Compare[int])
	require.False(t, ok)

	n, ok := MinFunc(slices.Values([]int{1, 2, 3}), cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 1, n)
}



func TestForEach2(t *testing.T) {
	t.Fail()
}

