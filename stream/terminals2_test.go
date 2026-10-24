package stream

import (
	cmp "cmp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mariomac/gostream/item"
)

func TestToMap(t *testing.T) {
	processed := OfMap(map[string]int{
		"Barcelona": 1,
		"Madrid":    2,
		"Paris":     3,
	}).Map(func(k string, v int) (string, int) {
		k = strings.ToLower(k)
		if strings.Contains(k, "i") {
			v++
		}
		return k, v
	})

	assert.Equal(t, map[string]int{
		"barcelona": 1,
		"madrid":    3,
		"paris":     4,
	}, processed.ToMap())
}

func isEmpty(k, _ string) bool { return k == "" }
func TestStream2_AllMatch(t *testing.T) {
	// for empty streams, following Java behavior as reference
	assert.True(t, Empty2[string, string]().AllMatch(isEmpty))
	assert.True(t, OfMap(map[string]string{
		"hello": "fo",
		"world": "bar",
	}).AllMatch(item.Not2[string, string](isEmpty)))
	assert.False(t, OfMap(map[string]string{
		"hello": "fo",
		"":      "bar",
	}).AllMatch(item.Not2[string, string](isEmpty)))
}

func TestStream2_AnyMatch(t *testing.T) {
	// for empty streams, following Java behavior as reference
	assert.False(t, Empty2[string, string]().AnyMatch(isEmpty))
	assert.True(t, OfMap(map[string]string{
		"hello": "fo",
		"world": "bar",
	}).AnyMatch(item.Not2[string, string](isEmpty)))
	assert.True(t, OfMap(map[string]string{
		"hello": "fo",
		"":      "bar",
	}).AnyMatch(item.Not2[string, string](isEmpty)))
	assert.False(t, OfMap(map[string]string{
		"": "bar",
	}).AnyMatch(item.Not2[string, string](isEmpty)))
}

func TestStream2_NoneMatch(t *testing.T) {
	// for empty streams, following Java behavior as reference
	assert.True(t, Empty2[string, string]().NoneMatch(isEmpty))
	assert.True(t, OfMap(map[string]string{
		"hello": "fo",
		"world": "bar",
	}).NoneMatch(isEmpty))
	assert.False(t, OfMap(map[string]string{
		"hello": "fo",
		"":      "bar",
	}).NoneMatch(isEmpty))
	assert.False(t, OfMap(map[string]string{
		"": "bar",
	}).NoneMatch(isEmpty))
}

func TestCount2(t *testing.T) {
	assert.Equal(t, 0, Empty2[int, int]().Count())
	assert.Equal(t, 0, OfMap(map[int]int{1: 1, 2: 1, 3: 1}).Skip(3).Count())
	assert.Equal(t, 3, OfMap(map[int]int{1: 1, 2: 1, 3: 1}).Count())
	assert.Equal(t, 3, OfMap(map[int]int{1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1}).Skip(3).Count())
	assert.Equal(t, 8, Iterate2(1, 1, func(k, v int) (int, int) {
		return k + 1, v + 1
	}).Limit(8).Count())
}

func TestFindFirst2(t *testing.T) {
	_, _, ok := Empty2[int, int]().FindFirst()
	require.False(t, ok)

	_, _, ok = OfMap(map[int]int{1: 1, 2: 1, 3: 1}).Skip(3).FindFirst()
	require.False(t, ok)

	_, n, ok := OfMap(map[int]int{1: 1, 2: 1, 3: 1}).FindFirst()
	require.True(t, ok)
	assert.Equal(t, 1, n)

	n, ok = Of(1, 2, 3, 4, 5, 6).Skip(3).FindFirst()
	require.True(t, ok)
	assert.Equal(t, 4, n)

	n, ok = Iterate[int](1, item.Increment[int]).Limit(8).FindFirst()
	require.True(t, ok)
	assert.Equal(t, 1, n)
}

func TestMax(t *testing.T) {
	_, _, ok := Empty2[int, int]().Max(cmp.Compare[int])
	require.False(t, ok)

	_, _, ok = OfMap(map[int]int{1: 1, 2: 1, 3: 1}).Skip(3).Max(cmp.Compare[int])
	require.False(t, ok)

	_, n, ok := OfMap(map[int]int{1: 1, 2: 1, 3: 1}).Skip(2).Max(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = OfMap(map[int]int{1: 1, 2: 1, 3: 1}).Max(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = Of(1, 2, 3, 4, 5, 6).Skip(3).Max(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 6, n)

	n, ok = Iterate[int](1, item.Increment[int]).Limit(8).Max(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 8, n)
}

func TestMin(t *testing.T) {
	_, _, ok := Empty[int]().Min(cmp.Compare[int])
	require.False(t, ok)

	_, _, ok = OfMap(map[int]int{1: 1, 2: 1, 3: 1}).Skip(3).Min(cmp.Compare[int])
	require.False(t, ok)

	_, n, ok := OfMap(map[int]int{1: 1, 2: 1, 3: 1}).Skip(2).Min(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 3, n)

	n, ok = OfMap(map[int]int{1: 1, 2: 1, 3: 1}).Min(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 1, n)

	n, ok = Of(1, 2, 3, 4, 5, 6).Skip(3).Min(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 4, n)

	n, ok = Iterate[int](1, item.Increment[int]).Limit(8).Min(cmp.Compare[int])
	require.True(t, ok)
	assert.Equal(t, 1, n)
}
