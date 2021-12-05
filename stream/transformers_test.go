package stream

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	in := Of[int](1, 2, 3, 4, 5)
	double := in.Map(func(i int) int {
		return i * 2
	})
	assert.Equal(t, []int{2, 4, 6, 8, 10}, double.ToSlice())

	strings := Map(in, strconv.Itoa)
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