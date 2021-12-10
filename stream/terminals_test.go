package stream

import (
	"strings"
	"testing"

	"github.com/mariomac/gostream/item"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
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

func TestToMap(t *testing.T) {
	processed := OfMap(map[string]int{
		"Barcelona": 1,
		"Madrid":    2,
		"Paris":     3,
	}).Map(func(p item.Pair[string, int]) item.Pair[string, int] {
		p.Key = strings.ToLower(p.Key)
		if strings.Contains(p.Key, "i") {
			p.Val++
		}
		return p
	})

	assert.Equal(t, map[string]int{
		"barcelona": 1,
		"madrid":    3,
		"paris":     4,
	}, ToMap(processed))
}

func TestReduce(t *testing.T) {
	// test empty stream
	_, ok := Empty[int]().Reduce(item.Add[int])
	assert.False(t, ok)

	// test one-element stream
	red, ok := Of(8).Reduce(item.Add[int])
	assert.True(t, ok)
	assert.Equal(t, 8, red)

	// test multi-element stream
	red, ok = Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).Reduce(item.Add[int])
	assert.True(t, ok)
	assert.Equal(t, 55, red)
}

func TestIterableStream_AllMatch(t *testing.T) {
	// for empty streams, following Java behavior as reference
	assert.True(t, Empty[string]().AllMatch(item.IsZero[string]))
	assert.True(t, Of("hello", "world").AllMatch(item.Not(item.IsZero[string])))
	assert.False(t, Of("", "world").AllMatch(item.Not(item.IsZero[string])))
}

func TestIterableStream_AnyMatch(t *testing.T) {
	// for empty streams, following Java behavior as reference
	assert.False(t, Empty[string]().AnyMatch(item.IsZero[string]))
	assert.True(t, Of("hello", "world").AnyMatch(item.Not(item.IsZero[string])))
	assert.True(t, Of("", "world").AnyMatch(item.Not(item.IsZero[string])))
	assert.False(t, Of("", "").AnyMatch(item.Not(item.IsZero[string])))
}

func TestIterableStream_NoneMatch(t *testing.T) {
	// for empty streams, following Java behavior as reference
	assert.True(t, Empty[string]().NoneMatch(item.IsZero[string]))
	assert.False(t, Of("hello", "world").NoneMatch(item.Not(item.IsZero[string])))
	assert.False(t, Of("", "world").NoneMatch(item.Not(item.IsZero[string])))
	assert.True(t, Of("", "").NoneMatch(item.Not(item.IsZero[string])))
}

func TestCount(t *testing.T) {
	assert.Equal(t, 0, Empty[int]().Count())
	assert.Equal(t, 0, Of(1, 2, 3).Skip(3).Count())
	assert.Equal(t, 3, Of(1, 2, 3).Count())
	assert.Equal(t, 3, Of(1, 2, 3, 4, 5, 6).Skip(3).Count())
	assert.Equal(t, 8, Iterate[int](1, item.Increment[int]).Limit(8).Count())
}
